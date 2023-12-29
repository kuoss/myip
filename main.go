package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug   bool
	Addr    string `default:":80"`
	Proxies []string
}

func loadConfig() (Config, error) {
	var cfg Config
	err := envconfig.Process("app", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("process err: %w", err)
	}
	log.Println("IP App starting...")
	log.Println("Addr:", cfg.Addr)
	log.Println("Proxies:", cfg.Proxies)
	return cfg, nil
}

func setupRouter(cfg Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.SetTrustedProxies(cfg.Proxies)
	router.ForwardedByClientIP = true
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, c.ClientIP()+"\n")
	})
	return router
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("loadConfig err: %w", err)
	}
	log.Println("IP App started...")
	r := setupRouter(cfg)
	return r.Run(cfg.Addr)
}

func main() {
	err := run()
	if err != nil {
		log.Printf("run err: %s", err.Error())
	}
}
