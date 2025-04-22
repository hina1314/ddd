//go:build wireinject
// +build wireinject

package di

import (
	"database/sql"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
	_ "github.com/lib/pq"
	"study/config"
	"study/db/model"
	"study/internal/api/handler"
	"study/internal/api/response"
	"study/internal/app"
	"study/internal/domain/user"
	"study/internal/infra/repository"
	"study/token"
	"study/util/errors"
	"study/util/i18n"
)

// Dependencies 包含应用程序的所有依赖。
type Dependencies struct {
	ResponseHandler *response.ResponseHandler
	UserHandler     *handler.UserHandler
	TokenMaker      token.Maker
	Config          config.Config // 使用值类型
	server          *fiber.App    // 非导出字段
}

// NewServer 返回 Fiber 服务器实例。
func (d *Dependencies) NewServer() *fiber.App {
	if d.server == nil {
		d.server = fiber.New()
	}
	return d.server
}

// NewDependencies 初始化所有依赖。
func NewDependencies(cfg config.Config) (*Dependencies, error) {
	deps, err := initializeDependencies(cfg)
	if err != nil {
		return nil, err
	}
	deps.Config = cfg
	return deps, nil
}

func initializeDependencies(cfg config.Config) (*Dependencies, error) {
	wire.Build(
		// 基础设施层
		newFiberApp, // 新增提供者
		newDB,
		newTokenMaker,
		newErrorHandler,
		newValidator,
		wire.Bind(new(i18n.Translator), new(*i18n.FileTranslator)),
		newFileTranslator,
		newTranslationService,
		repository.NewUserRepository,
		repository.NewUserAccountRepository,

		// 领域层
		user.NewDomainService,

		// 应用层
		app.NewUserService,

		// 表现层
		response.NewResponseHandler,
		handler.NewUserHandler,

		// 返回值
		wire.Struct(new(Dependencies), "*"),
	)
	return nil, nil
}

func newFiberApp() *fiber.App {
	return fiber.New()
}

func newDB(cfg config.Config) (model.TxManager, error) {
	db, err := sql.Open("postgres", cfg.DBSource)
	if err != nil {
		return nil, err
	}
	return model.NewStore(db), nil
}

func newTokenMaker(cfg config.Config) (token.Maker, error) {
	return token.NewPasetoMaker(cfg.TokenSymmetricKey)
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
func newValidator() *validator.Validate {
	v := validator.New()
	// 注册自定义 tag "phone"
	if err := v.RegisterValidation("phone", errors.PhoneValidator); err != nil {
		panic(err)
	}
	return v
}
