package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	os.Clearenv()
	os.Setenv("APP_DEBUG", "hello")
	cfg, err := loadConfig()
	require.EqualError(t, err, "process err: envconfig.Process: assigning APP_DEBUG to Debug: converting 'hello' to type bool. details: strconv.ParseBool: parsing \"hello\": invalid syntax")
	require.Equal(t, Config{}, cfg)

	os.Clearenv()
	cfg, err = loadConfig()
	require.NoError(t, err)
	require.Equal(t, Config{Addr: ":80"}, cfg)
}

func TestSetupRouter(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		router, err := setupRouter(Config{Proxies: []string{"hello", "world"}})
		require.EqualError(t, err, "setTrustedProxies err: invalid IP address: hello")
		require.Nil(t, router)
	})
	t.Run("ok", func(t *testing.T) {
		router, err := setupRouter(Config{})
		require.NoError(t, err)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(w, req)
		require.Equal(t, 200, w.Code)
		require.Equal(t, "\n", w.Body.String())
	})
}

func TestRun(t *testing.T) {
	t.Run("error APP_DEBUG", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("APP_DEBUG", "hello")
		err := run()
		require.EqualError(t, err, `loadConfig err: process err: envconfig.Process: assigning APP_DEBUG to Debug: converting 'hello' to type bool. details: strconv.ParseBool: parsing "hello": invalid syntax`)
	})

	t.Run("error APP_ADDR", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("APP_ADDR", "hello")
		err := run()
		require.EqualError(t, err, "listen tcp: address hello: missing port in address")
	})

	t.Run("error APP_PROXIES", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("APP_PROXIES", "hello")
		err := run()
		require.EqualError(t, err, "setupRouter err: setTrustedProxies err: invalid IP address: hello")
	})
}

func TestMain(t *testing.T) {
	os.Clearenv()
	os.Setenv("APP_ADDR", "hello")
	main()
}
