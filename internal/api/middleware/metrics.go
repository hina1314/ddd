package middleware

import (
	"crypto/subtle"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "app",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests.",
		},
		[]string{"method", "route", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "app",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request duration in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "route", "status"},
	)

	httpRequestsInFlight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "app",
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

func MetricsAuth(expectedToken string) fiber.Handler {
	return func(c fiber.Ctx) error {
		if expectedToken == "" {
			return c.Next()
		}
		fields := strings.Fields(c.Get(fiber.HeaderAuthorization))
		if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		provided := fields[1]
		if subtle.ConstantTimeCompare([]byte(provided), []byte(expectedToken)) != 1 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Next()
	}
}

type databaseStatsProvider interface {
	Stats() sql.DBStats
}

func RegisterDatabaseMetrics(db databaseStatsProvider) error {
	collectors := []prometheus.Collector{
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "app_db_open_connections", Help: "Open database connections."}, func() float64 { return float64(db.Stats().OpenConnections) }),
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "app_db_in_use_connections", Help: "Database connections currently in use."}, func() float64 { return float64(db.Stats().InUse) }),
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "app_db_idle_connections", Help: "Idle database connections."}, func() float64 { return float64(db.Stats().Idle) }),
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "app_db_wait_count_total", Help: "Total waits for a database connection."}, func() float64 { return float64(db.Stats().WaitCount) }),
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "app_db_wait_duration_seconds_total", Help: "Total time waiting for database connections."}, func() float64 { return db.Stats().WaitDuration.Seconds() }),
	}
	for _, collector := range collectors {
		if err := prometheus.Register(collector); err != nil {
			return err
		}
	}
	return nil
}
