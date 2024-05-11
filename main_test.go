package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) { //nolint:tparallel // using t.Setenv()
	t.Run("error", func(t *testing.T) {
		t.Setenv("APP_DEBUG", "hello")

		cfg, err := loadConfig()
		require.EqualError(t, err, "process err: envconfig.Process: assigning APP_DEBUG to Debug: converting 'hello' to type bool. details: strconv.ParseBool: parsing \"hello\": invalid syntax")
		require.Nil(t, cfg)
	})
	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		cfg, err := loadConfig()
		require.NoError(t, err)
		require.Equal(t, &Config{Addr: ":80"}, cfg)
	})
}

func TestSetupRouter(t *testing.T) {
	t.Parallel()
	t.Run("error with nil", func(t *testing.T) {
		t.Parallel()

		router, err := setupRouter(nil)
		require.Equal(t, ErrConfigIsNil, err)
		require.Nil(t, router)
	})

	t.Run("error with hello world", func(t *testing.T) {
		t.Parallel()

		router, err := setupRouter(&Config{Proxies: []string{"hello", "world"}})
		require.EqualError(t, err, "setTrustedProxies err: invalid IP address: hello")
		require.Nil(t, router)
	})

	t.Run("ok with httptest request", func(t *testing.T) {
		t.Parallel()

		router, err := setupRouter(&Config{})
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		router.ServeHTTP(w, req)
		require.Equal(t, 200, w.Code)
		require.Equal(t, "192.0.2.1\n", w.Body.String())
	})

	t.Run("ok with http request", func(t *testing.T) {
		t.Parallel()

		router, err := setupRouter(&Config{})
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, "/", http.NoBody)
		require.NoError(t, err)
		router.ServeHTTP(w, req)
		require.Equal(t, 200, w.Code)
		require.Equal(t, "\n", w.Body.String())
	})
}

func TestRun(t *testing.T) {
	t.Run("error APP_DEBUG", func(t *testing.T) {
		t.Setenv("APP_DEBUG", "hello")

		err := run()
		require.EqualError(t, err, `loadConfig err: process err: envconfig.Process: assigning APP_DEBUG to Debug: converting 'hello' to type bool. details: strconv.ParseBool: parsing "hello": invalid syntax`)
	})

	t.Run("error APP_ADDR", func(t *testing.T) {
		t.Setenv("APP_ADDR", "hello")

		err := run()
		require.EqualError(t, err, "setupRouter err: listen tcp: address hello: missing port in address")
	})

	t.Run("error APP_PROXIES", func(t *testing.T) {
		t.Setenv("APP_PROXIES", "hello")

		err := run()
		require.EqualError(t, err, "setupRouter err: setTrustedProxies err: invalid IP address: hello")
	})
}

func Test_main(t *testing.T) {
	t.Setenv("APP_ADDR", "hello")
	main()
}
