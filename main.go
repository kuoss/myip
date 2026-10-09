package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strings"
	"time"
)

var ErrConfigIsNil = errors.New("config is nil")

type Config struct {
	Addr    string
	Proxies []string
}

func loadConfig() *Config {
	cfg := &Config{Addr: ":8080"}

	// Load Addr
	addr := os.Getenv("APP_ADDR")
	if addr != "" {
		cfg.Addr = addr
	}

	// Load Proxies
	proxiesStr := os.Getenv("APP_PROXIES")
	if proxiesStr != "" {
		cfg.Proxies = strings.Split(proxiesStr, ",")
	}

	log.Println("IP App starting...")
	log.Println("Addr:", cfg.Addr)
	log.Println("Proxies:", cfg.Proxies)

	return cfg
}

// parseProxies accepts CIDRs ("10.0.0.0/8") and single IPs ("192.168.1.33").
func parseProxies(proxies []string) ([]netip.Prefix, error) {
	prefixes := make([]netip.Prefix, 0, len(proxies))
	for _, p := range proxies {
		p = strings.TrimSpace(p)
		if strings.Contains(p, "/") {
			prefix, err := netip.ParsePrefix(p)
			if err != nil {
				return nil, fmt.Errorf("invalid CIDR: %s", p)
			}
			prefixes = append(prefixes, prefix.Masked())
			continue
		}
		ip, err := netip.ParseAddr(p)
		if err != nil {
			return nil, fmt.Errorf("invalid IP address: %s", p)
		}
		ip = ip.Unmap()
		prefixes = append(prefixes, netip.PrefixFrom(ip, ip.BitLen()))
	}
	return prefixes, nil
}

type clientIPResolver struct {
	trusted []netip.Prefix
}

func (r clientIPResolver) isTrusted(ip netip.Addr) bool {
	ip = ip.Unmap()
	for _, prefix := range r.trusted {
		if prefix.Contains(ip) {
			return true
		}
	}
	return false
}

// clientIP returns the remote address, or, when the request comes from a trusted proxy,
// the address it forwarded: X-Forwarded-For read right to left up to the first untrusted
// hop, then X-Real-IP.
func (r clientIPResolver) clientIP(req *http.Request) string {
	remote := strings.TrimSpace(req.RemoteAddr)
	if host, _, err := net.SplitHostPort(remote); err == nil {
		remote = host
	}
	remoteIP, err := netip.ParseAddr(remote)
	if err != nil || !r.isTrusted(remoteIP) {
		return remote
	}
	if ip, ok := r.fromForwardedFor(req.Header.Get("X-Forwarded-For")); ok {
		return ip
	}
	if ip, err := netip.ParseAddr(strings.TrimSpace(req.Header.Get("X-Real-IP"))); err == nil {
		return ip.String()
	}
	return remote
}

func (r clientIPResolver) fromForwardedFor(header string) (string, bool) {
	if header == "" {
		return "", false
	}
	items := strings.Split(header, ",") // never empty
	for i := len(items) - 1; ; i-- {
		ip, err := netip.ParseAddr(strings.TrimSpace(items[i]))
		if err != nil {
			return "", false
		}
		if i == 0 || !r.isTrusted(ip) {
			return ip.String(), true
		}
	}
}

func setupRouter(cfg *Config) (http.Handler, error) {
	if cfg == nil {
		return nil, ErrConfigIsNil
	}

	trusted, err := parseProxies(cfg.Proxies)
	if err != nil {
		return nil, fmt.Errorf("parse proxies: %w", err)
	}
	resolver := clientIPResolver{trusted: trusted}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintln(w, resolver.clientIP(r))
	})

	return mux, nil
}

func run() error {
	log.Println("IP App started...")

	cfg := loadConfig()

	handler, err := setupRouter(cfg)
	if err != nil {
		return fmt.Errorf("setupRouter err: %w", err)
	}

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return srv.ListenAndServe() //nolint:wrapcheck // returned as is to main
}

func main() {
	err := run()
	if err != nil {
		log.Printf("run err: %s", err.Error())
	}
}
