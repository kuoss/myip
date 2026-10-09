package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("APP_ADDR", "")
	t.Setenv("APP_PROXIES", "")

	cfg := loadConfig()
	if cfg.Addr != ":8080" || cfg.Proxies != nil {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadConfig_FromEnv(t *testing.T) {
	t.Setenv("APP_ADDR", ":8000")
	t.Setenv("APP_PROXIES", "10.0.0.0/8,192.168.1.33")

	cfg := loadConfig()
	if cfg.Addr != ":8000" || len(cfg.Proxies) != 2 || cfg.Proxies[1] != "192.168.1.33" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestSetupRouter_Errors(t *testing.T) {
	t.Parallel()

	if _, err := setupRouter(nil); err != ErrConfigIsNil {
		t.Fatalf("nil config: err = %v", err)
	}

	_, err := setupRouter(&Config{Proxies: []string{"hello", "world"}})
	if err == nil || err.Error() != "parse proxies: invalid IP address: hello" {
		t.Fatalf("invalid proxy: err = %v", err)
	}

	_, err = setupRouter(&Config{Proxies: []string{"10.0.0.0/33"}})
	if err == nil || err.Error() != "parse proxies: invalid CIDR: 10.0.0.0/33" {
		t.Fatalf("invalid CIDR: err = %v", err)
	}
}

func serve(t *testing.T, cfg *Config, method, target, remoteAddr string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	handler, err := setupRouter(cfg)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(method, target, http.NoBody)
	req.RemoteAddr = remoteAddr
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func TestClientIP(t *testing.T) {
	t.Parallel()

	proxies := &Config{Proxies: []string{"10.0.0.0/8", "192.168.1.33"}}
	tests := []struct {
		name    string
		cfg     *Config
		remote  string
		headers map[string]string
		want    string
	}{
		{"no proxies: remote address", &Config{}, "192.0.2.1:1234", nil, "192.0.2.1\n"},
		{"no proxies: forwarded header ignored", &Config{}, "192.0.2.1:1234", map[string]string{"X-Forwarded-For": "203.0.113.9"}, "192.0.2.1\n"},
		{"untrusted remote: header ignored", proxies, "192.0.2.1:1234", map[string]string{"X-Forwarded-For": "203.0.113.9"}, "192.0.2.1\n"},
		{"trusted CIDR: forwarded client", proxies, "10.1.2.3:1234", map[string]string{"X-Forwarded-For": "203.0.113.9"}, "203.0.113.9\n"},
		{"trusted single IP: forwarded client", proxies, "192.168.1.33:1234", map[string]string{"X-Forwarded-For": "203.0.113.9"}, "203.0.113.9\n"},
		{"skips trusted hops right to left", proxies, "10.1.2.3:1234", map[string]string{"X-Forwarded-For": "198.51.100.7, 203.0.113.9, 10.9.9.9"}, "203.0.113.9\n"},
		{"all hops trusted: leftmost", proxies, "10.1.2.3:1234", map[string]string{"X-Forwarded-For": "10.4.4.4, 10.9.9.9"}, "10.4.4.4\n"},
		{"X-Real-IP fallback", proxies, "10.1.2.3:1234", map[string]string{"X-Real-IP": "203.0.113.9"}, "203.0.113.9\n"},
		{"invalid forwarded: X-Real-IP", proxies, "10.1.2.3:1234", map[string]string{"X-Forwarded-For": "bogus", "X-Real-IP": "203.0.113.9"}, "203.0.113.9\n"},
		{"trusted, no headers: remote", proxies, "10.1.2.3:1234", nil, "10.1.2.3\n"},
		{"IPv6 remote", &Config{}, "[2001:db8::1]:1234", nil, "2001:db8::1\n"},
		{"empty remote address", &Config{}, "", nil, "\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := serve(t, tt.cfg, http.MethodGet, "/", tt.remote, tt.headers)
			if w.Code != http.StatusOK || w.Body.String() != tt.want {
				t.Fatalf("code=%d body=%q, want %q", w.Code, w.Body.String(), tt.want)
			}
		})
	}
}

func TestRoutes(t *testing.T) {
	t.Parallel()

	if w := serve(t, &Config{}, http.MethodGet, "/other", "192.0.2.1:1", nil); w.Code != http.StatusNotFound {
		t.Fatalf("GET /other: code=%d, want 404", w.Code)
	}
	if w := serve(t, &Config{}, http.MethodPost, "/", "192.0.2.1:1", nil); w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST /: code=%d, want 405", w.Code)
	}
}

func TestRun(t *testing.T) {
	getRandomAddr := func() string {
		server := httptest.NewServer(nil)
		server.Close()
		u, _ := url.Parse(server.URL)

		return ":" + u.Port()
	}

	t.Run("ok", func(t *testing.T) {
		addr := getRandomAddr()
		t.Setenv("APP_ADDR", addr)
		t.Setenv("APP_PROXIES", "")

		go run() //nolint:errcheck // smoke test
		for range 50 {
			resp, err := http.Get("http://127.0.0.1" + addr + "/")
			if err == nil {
				_ = resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					t.Fatalf("code=%d", resp.StatusCode)
				}
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
		t.Fatal("server did not start")
	})

	t.Run("error APP_ADDR", func(t *testing.T) {
		t.Setenv("APP_ADDR", "hello")
		t.Setenv("APP_PROXIES", "")

		err := run()
		if err == nil || err.Error() != "listen tcp: address hello: missing port in address" {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("error APP_PROXIES", func(t *testing.T) {
		t.Setenv("APP_ADDR", "")
		t.Setenv("APP_PROXIES", "hello")

		err := run()
		if err == nil || err.Error() != "setupRouter err: parse proxies: invalid IP address: hello" {
			t.Fatalf("err = %v", err)
		}
	})
}

func Test_main(t *testing.T) {
	t.Setenv("APP_ADDR", "hello")
	main()
}
