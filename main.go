package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug   bool     `exhaustruct:"optional"`
	Addr    string   `exhaustruct:"optional"`
	Proxies []string `exhaustruct:"optional"`
}

func loadConfig() (*Config, error) {
	cfg := &Config{Addr: ":80"}
	if err := envconfig.Process("app", cfg); err != nil {
		return nil, fmt.Errorf("process err: %w", err)
	}

	log.Println("IP App starting...")
	log.Println("Addr:", cfg.Addr)
	log.Println("Proxies:", cfg.Proxies)

	return cfg, nil
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
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("loadConfig err: %w", err)
	}

	log.Println("IP App started...")

	router, err := setupRouter(cfg)
	if err != nil {
		return fmt.Errorf("setupRouter err: %w", err)
	}

	return router.Run(cfg.Addr) //nolint:wrapcheck // never ends
}

func main() {
	err := run()
	if err != nil {
		log.Printf("run err: %s", err.Error())
	}
}
