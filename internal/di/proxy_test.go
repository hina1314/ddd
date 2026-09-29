package di

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/hina1314/ddd/config"
	"github.com/hina1314/kit/response"
	"github.com/stretchr/testify/require"
)

func TestTrustedProxyClientLimits(t *testing.T) {
	for _, tc := range []struct {
		name     string
		proxies  []string
		headers  []string
		statuses []int
	}{
		{"trusted proxy separates clients", []string{"0.0.0.0/32"}, []string{"198.51.100.1", "198.51.100.2", "198.51.100.1"}, []int{200, 200, 429}},
		{"empty allowlist ignores forged headers", nil, []string{"198.51.100.1", "198.51.100.2"}, []int{200, 429}},
		{"untrusted peer ignores forged headers", []string{"192.0.2.10"}, []string{"198.51.100.1", "198.51.100.2"}, []int{200, 429}},
		{"invalid header falls back to peer", []string{"0.0.0.0/32"}, []string{"invalid", ""}, []int{200, 429}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Fiber's in-memory Test transport has remote IP 0.0.0.0.
			app := newFiberApp(&response.ResponseHandler{}, config.Config{ServerTrustedProxies: tc.proxies})
			app.Use(limiter.New(limiter.Config{Max: 1, Expiration: time.Minute}))
			app.Post("/v1/login", func(c fiber.Ctx) error { return c.SendStatus(http.StatusOK) })
			for i, ip := range tc.headers {
				req := httptest.NewRequest(http.MethodPost, "/v1/login", nil)
				req.Header.Set("X-Real-IP", ip)
				res, err := app.Test(req)
				require.NoError(t, err)
				require.NoError(t, res.Body.Close())
				require.Equal(t, tc.statuses[i], res.StatusCode)
			}
		})
	}
}
