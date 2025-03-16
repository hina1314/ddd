package svc

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"study/db/model"
	"study/token"
	"study/util"
)

type ServiceContext struct {
	Config     util.Config
	Db         model.Store
	TokenMaker token.Maker
}

func NewServiceContext(c util.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		Db:         initDB(c),
		TokenMaker: initTokenMaker(c),
	}
}

func initDB(conf util.Config) model.Store {

	sqlConn, err := sql.Open("postgres", conf.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	return model.NewStore(sqlConn)
}

func initTokenMaker(conf util.Config) token.Maker {
	tokenMaker, err := token.NewPasetoMaker(conf.TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker:%w", err)
	}
	return tokenMaker
}
