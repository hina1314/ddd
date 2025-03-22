//go:build wireinject
// +build wireinject

package di

import (
	"database/sql"
	"github.com/google/wire"
	_ "github.com/lib/pq"
	"study/internal/api/handler"
	"study/internal/app"
	"study/internal/domain/user"
	"study/internal/infra/repository"
	"study/util"
)

type Dependencies struct {
	UserHandler *handler.UserHandler
}

func InitializeDependencies(cfg util.Config) (*Dependencies, error) {
	wire.Build(
		// 基础设施层
		NewDB,
		repository.NewUserRepository,
		repository.NewUserAccountRepository,

		// 领域层
		user.NewDomainService,

		// 应用层
		app.NewUserService,

		// 表现层
		handler.NewUserHandler,

		// 返回值
		wire.Struct(new(Dependencies), "UserHandler"),
	)
	return &Dependencies{}, nil
}

func NewDB(cfg util.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DBSource)
	if err != nil {
		return nil, err
	}
	return db, nil
}
