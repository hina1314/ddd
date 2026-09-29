//go:build wireinject
// +build wireinject

package di

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
	"github.com/hina1314/ddd/config"
	"github.com/hina1314/ddd/db/model"
	"github.com/hina1314/ddd/internal/api/handler"
	"github.com/hina1314/ddd/internal/api/presentation"
	"github.com/hina1314/ddd/internal/api/validation"
	"github.com/hina1314/ddd/internal/app/user"
	userService "github.com/hina1314/ddd/internal/domain/user/service"
	userRepo "github.com/hina1314/ddd/internal/infra/user"
	"github.com/hina1314/ddd/token"
	"github.com/hina1314/kit/database"
	"github.com/hina1314/kit/errors"
	"github.com/hina1314/kit/i18n"
	"github.com/hina1314/kit/response"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"time"
)

// Dependencies 包含应用程序的所有依赖。
type Dependencies struct {
	DB              *model.SQLStore
	ResponseHandler *response.ResponseHandler
	UserHandler     *handler.UserHandler
	TokenMaker      token.Maker
	Config          config.Config // 使用值类型
	server          *fiber.App    // 非导出字段
}

// NewServer 返回 Fiber 服务器实例。
func (d *Dependencies) NewServer() *fiber.App {
	if d.server == nil {
		d.server = newFiberApp(d.ResponseHandler, d.Config)
	}
	return d.server
}

// NewDependencies 初始化所有依赖。
func NewDependencies(cfg config.Config) (*Dependencies, error) {
	deps, _, err := initializeDependencies(cfg)
	if err != nil {
		return nil, err
	}
	deps.Config = cfg
	return deps, nil
}

func initializeDependencies(cfg config.Config) (*Dependencies, func(), error) {
	wire.Build(
		// 基础设施层
		newFiberApp, // 新增提供者
		newDB,
		wire.Bind(new(model.TxManager), new(*model.SQLStore)),
		newTokenMaker,
		newErrorHandler,
		newValidator,
		wire.Bind(new(i18n.Translator), new(*i18n.FileTranslator)),
		newFileTranslator,
		newTranslationService,
		userRepo.NewUserRepository,
		// 领域层
		// user
		userService.NewUserLoginService,
		userService.NewUserRegisterService,
		userService.NewUserUpdateService,

		// 应用层
		user.NewUserService,
		// 表现层
		presentation.NewResponseHandler,
		handler.NewUserHandler,

		// 返回值
		wire.Struct(new(Dependencies), "*"),
	)
	return nil, nil, nil
}

func newFiberApp(
	responseHandler *response.ResponseHandler,
	cfg config.Config,
) *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler:       responseHandler.HandleError,
		ReadTimeout:        cfg.ServerReadTimeout,
		WriteTimeout:       cfg.ServerWriteTimeout,
		IdleTimeout:        cfg.ServerIdleTimeout,
		BodyLimit:          cfg.ServerBodyLimit,
		ServerHeader:       cfg.AppName,
		ProxyHeader:        "X-Real-IP",
		TrustProxy:         true,
		TrustProxyConfig:   fiber.TrustProxyConfig{Proxies: cfg.ServerTrustedProxies},
		EnableIPValidation: true,
	})
}

func newDB(cfg config.Config) (*model.SQLStore, func(), error) {
	db, err := database.Open(context.Background(), database.Config{
		Driver: cfg.DBDriver, Source: cfg.DBSource,
		MaxOpenConns: cfg.DBMaxOpenConns, MaxIdleConns: cfg.DBMaxIdleConns,
		ConnMaxLifetime: cfg.DBConnMaxLifetime, ConnMaxIdleTime: cfg.DBConnMaxIdleTime,
		PingTimeout: 3 * time.Second,
	})
	if err != nil {
		return nil, nil, err
	}
	// Wire 在后续依赖装配失败时调用清理函数；成功后由入口管理关闭。
	cleanup := func() {
		if err := db.Close(); err != nil {
			slog.Error("close database after dependency initialization failure", "error", err)
		}
	}
	return model.NewStore(db), cleanup, nil
}

func newTokenMaker(cfg config.Config) (token.Maker, error) {
	return token.NewPasetoMaker(cfg.TokenSymmetricKey, cfg.TokenIssuer, cfg.TokenAudience)
}

func newErrorHandler(cfg config.Config) *errors.ErrorHandler {
	return errors.NewErrorHandler(cfg.Debug, false)
}

func newFileTranslator() *i18n.FileTranslator {
	return i18n.NewFileTranslator("en")
}

func newTranslationService(translator i18n.Translator, cfg config.Config) (*i18n.TranslationService, error) {
	if err := translator.LoadTranslations("./config/i18n"); err != nil {
		return nil, err
	}
	return i18n.NewTranslationService(translator, cfg.DefaultLocale), nil
}
func newValidator() (*validator.Validate, error) {
	v := validator.New()
	if err := v.RegisterValidation("phone", validation.PhoneValidator); err != nil {
		return nil, fmt.Errorf("register phone validation: %w", err)
	}
	if err := v.RegisterValidation("password", validation.PasswordValidator); err != nil {
		return nil, fmt.Errorf("register password validation: %w", err)
	}
	return v, nil
}
