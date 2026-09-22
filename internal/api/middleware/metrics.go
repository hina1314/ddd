package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "study",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests.",
		},
		[]string{"method", "route", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "study",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request duration in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "route", "status"},
	)

	httpRequestsInFlight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "study",
			Subsystem: "http",
			Name:      "requests_in_flight",
			Help:      "Current number of HTTP requests being processed.",
		},
		[]string{"method"},
	)
)

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		httpRequestsInFlight,
	)
}

func Metrics() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Prometheus 抓取本身不计入业务指标。
		if c.Path() == "/metrics" {
			return c.Next()
		}

		method := c.Method()
		start := time.Now()

		httpRequestsInFlight.WithLabelValues(method).Inc()
		defer httpRequestsInFlight.WithLabelValues(method).Dec()

		err := c.Next()
		status := responseStatus(c, err)

		// 使用注册的路由模板，不能使用带用户 ID 的原始路径。
		route := c.Route().Path
		if route == "" {
			route = "/__unmatched__"
		}

		labels := []string{
			method,
			route,
			strconv.Itoa(status),
		}

		httpRequestsTotal.WithLabelValues(labels...).Inc()

		httpRequestDuration.
			WithLabelValues(labels...).
			Observe(time.Since(start).Seconds())

		return err
	}
}

func MetricsHandler() fiber.Handler {
	return adaptor.HTTPHandler(promhttp.Handler())
}
