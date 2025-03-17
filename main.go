package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var ErrConfigIsNil = errors.New("config is nil")

type Config struct {
	Addr    string   `exhaustruct:"optional"`
	Proxies []string `exhaustruct:"optional"`
}

func loadConfig() *Config {
	cfg := &Config{Addr: ":80"}

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

func setupRouter(cfg *Config) (*gin.Engine, error) {
	router := gin.New()

	if cfg == nil {
		return nil, ErrConfigIsNil
	}

	if err := router.SetTrustedProxies(cfg.Proxies); err != nil {
		return nil, fmt.Errorf("setTrustedProxies err: %w", err)
	}

	router.ForwardedByClientIP = true
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, c.ClientIP()+"\n")
	})

	return router, nil
}

func run() error {
	log.Println("IP App started...")

	cfg := loadConfig()

	router, err := setupRouter(cfg)
	if err != nil {
		return fmt.Errorf("setupRouter err: %w", err)
	}

	return router.Run(cfg.Addr) //nolint:wrapcheck //nolint:gocritic
}

func main() {
	err := run()
	if err != nil {
		log.Printf("run err: %s", err.Error())
	}
}
