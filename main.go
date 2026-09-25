package main

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"study/config"
	"study/internal/api/router"
	"study/internal/di"
	"sync/atomic"
	"syscall"
	"time"
)

type databasePinger interface {
	PingContext(context.Context) error
}

// Release 构建通过 -ldflags 注入 tag；本地构建保留 dev。
var version = "dev"

// main 是应用程序的入口点。
func main() {
	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(os.Stdout, nil),
		).With("service", "study-api"),
	)
	slog.Info("application entrypoint reached", "version", version)
	// 加载配置
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	slog.Info("configuration loaded")

	// 初始化依赖
	deps, err := di.NewDependencies(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	slog.Info("dependencies initialized")

	// 创建 Fiber 服务器
	server := deps.NewServer()
	// 设置中间件
	router.SetupMiddleware(server, deps)

	var accepting atomic.Bool
	var dbHealthy atomic.Bool

	dbHealthy.Store(true) // 启动时的 PingContext 已经成功

	server.Get("/livez", func(c fiber.Ctx) error {
		return c.SendString("alive\n")
	})

	server.Get("/readyz", func(c fiber.Ctx) error {
		c.Set("X-App-Version", version)
		if !accepting.Load() || !dbHealthy.Load() {
			return c.Status(fiber.StatusServiceUnavailable).
				SendString("not ready\n")
		}

		return c.SendString("ready\n")
	})

	// 设置业务路由
	router.Setup(server, deps)

	// 仅在本地实验时启用，不访问数据库。
	if os.Getenv("SHUTDOWN_LAB") == "1" {
		server.Get("/_lab/slow", func(c fiber.Ctx) error {
			log.Println("[lab] request started")
			time.Sleep(10 * time.Second)
			log.Println("[lab] work finished, sending response")
			return c.SendString("completed\n")
		})

		server.Get("/_lab/error", func(c fiber.Ctx) error {
			return fiber.NewError(
				fiber.StatusInternalServerError,
				"simulated lab failure",
			)
		})
	}

	healthCtx, stopHealth := context.WithCancel(context.Background())
	defer stopHealth()

	go monitorDatabase(healthCtx, deps.DB, &dbHealthy)

	accepting.Store(true)

	// 启动服务器
	// 接收 Ctrl+C，让程序自己安排退出。
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(quit)

	// 监听放到另一个 goroutine，main 负责协调停机。
	listenErr := make(chan error, 1)
	go func() {
		slog.Info("starting HTTP listener", "address", cfg.ServerAddress)
		listenErr <- server.Listen(cfg.ServerAddress)
	}()

	select {
	case <-quit:
		accepting.Store(false)
		stopHealth()

		log.Println("[shutdown] readiness disabled, waiting for traffic removal")
		time.Sleep(cfg.TrafficDrainDelay)
		log.Println("[shutdown] draining active requests")
	case err := <-listenErr:
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
		return
	}

	exitCode := 0

	// 先完成或中断 HTTP 请求，此时数据库仍保持可用。
	if err := server.ShutdownWithTimeout(cfg.ShutdownTimeout); err != nil {
		log.Printf("[shutdown] HTTP server failed or timed out: %v", err)
		exitCode = 1
	}

	// HTTP 请求处理结束后，再关闭数据库连接池。
	if err := deps.DB.Close(); err != nil {
		log.Printf("[shutdown] database close failed: %v", err)
		exitCode = 1
	} else {
		log.Println("[shutdown] database pool closed")
	}

	if exitCode != 0 {
		os.Exit(exitCode)
	}

	log.Println("[shutdown] completed")
}

func monitorDatabase(
	ctx context.Context,
	db databasePinger,
	healthy *atomic.Bool,
) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	failures := 0
	successes := 0

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			probeCtx, cancel := context.WithTimeout(ctx, time.Second)
			err := db.PingContext(probeCtx)
			cancel()

			if ctx.Err() != nil {
				return
			}

			if err != nil {
				successes = 0
				failures++

				if failures >= 3 {
					failures = 3

					if healthy.CompareAndSwap(true, false) {
						log.Printf(
							"[health] database unhealthy after 3 failures: %v",
							err,
						)
					}
				}

				continue
			}

			failures = 0
			successes++

			if successes >= 2 {
				successes = 2

				if healthy.CompareAndSwap(false, true) {
					log.Println(
						"[health] database healthy after 2 successes",
					)
				}
			}
		}
	}
}
