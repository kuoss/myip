package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) { //nolint:tparallel // using t.Setenv()
	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		cfg := loadConfig()
		require.Equal(t, &Config{Addr: ":8080"}, cfg)
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
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/", http.NoBody)
		require.NoError(t, err)
		router.ServeHTTP(w, req)
		require.Equal(t, 200, w.Code)
		require.Equal(t, "\n", w.Body.String())
	})
}

func TestRun(t *testing.T) {
	getRandomAddr := func() string {
		server := httptest.NewServer(nil)
		server.Close()
		u, _ := url.Parse(server.URL)

		return ":" + u.Port()
	}

	t.Run("ok", func(t *testing.T) {
		t.Setenv("APP_ADDR", getRandomAddr())

		go run() //nolint:errcheck // smoke test
		time.Sleep(500 * time.Microsecond)
	})

	t.Run("error APP_ADDR", func(t *testing.T) {
		t.Setenv("APP_ADDR", "hello")

		err := run()
		require.EqualError(t, err, "listen tcp: address hello: missing port in address")
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
