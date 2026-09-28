package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/hina1314/ddd/config"
	"github.com/hina1314/ddd/internal/api/router"
	"github.com/hina1314/ddd/internal/di"
	"github.com/hina1314/kit/health"
	"github.com/hina1314/kit/lifecycle"
	"github.com/hina1314/kit/metrics"
)

// Release 构建通过 -ldflags 注入 tag；本地构建保留 dev。
var version = "dev"

func main() {
	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
	slog.Info("shutdown completed")
}

// run 负责业务装配和资源所有权，公共包负责健康探测与停机协调。
func run() (result error) {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		return err
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)).
		With("service", cfg.AppName, "environment", cfg.Environment))
	slog.Info("application entrypoint reached", "version", version)
	slog.Info("configuration loaded")

	deps, err := di.NewDependencies(cfg)
	if err != nil {
		return err
	}
	// lifecycle.Run 返回时，HTTP 关闭流程与健康监控停止均已完成。
	// 入口后续装配或监听启动失败也会走此清理路径。
	defer func() {
		if err := deps.DB.Close(); err != nil {
			result = errors.Join(result, err)
		} else {
			slog.Info("database pool closed")
		}
	}()
	slog.Info("dependencies initialized")

	instrumentation, err := metrics.New(metrics.Config{})
	if err != nil {
		return err
	}
	if err := instrumentation.RegisterDatabase(deps.DB); err != nil {
		return err
	}
	checks, err := health.New(health.Config{
		Version: version, Interval: cfg.HealthCheckInterval, Timeout: cfg.HealthCheckTimeout,
	}, health.Probe{Name: "database", Check: deps.DB.PingContext, InitiallyHealthy: true})
	if err != nil {
		return err
	}

	server := deps.NewServer()
	router.SetupMiddleware(server, deps, instrumentation)
	server.Get("/livez", checks.Liveness)
	server.Get("/readyz", checks.Readiness)
	router.Setup(server, deps)
	if os.Getenv("SHUTDOWN_LAB") == "1" {
		server.Get("/_lab/slow", func(c fiber.Ctx) error {
			slog.Info("lab request started")
			time.Sleep(10 * time.Second)
			slog.Info("lab work finished")
			return c.SendString("completed\n")
		})
		server.Get("/_lab/error", func(c fiber.Ctx) error {
			return fiber.NewError(fiber.StatusInternalServerError, "simulated lab failure")
		})
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return lifecycle.Run(ctx, lifecycle.Service{
		Listen: func() error {
			slog.Info("starting HTTP listener", "address", cfg.ServerAddress)
			return server.Listen(cfg.ServerAddress)
		},
		Shutdown:     server.ShutdownWithTimeout,
		SetAccepting: checks.SetAccepting,
		Monitor:      checks.Run,
	}, lifecycle.Config{TrafficDrainDelay: cfg.TrafficDrainDelay, ShutdownTimeout: cfg.ShutdownTimeout})
}
