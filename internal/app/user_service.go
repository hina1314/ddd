package app

import (
	"context"
	"errors"
	"study/db/model"
	"study/internal/domain/user"
)

type UserService struct {
	domainService *user.DomainService
	model.Store
}

func NewUserService(domainService *user.DomainService, store model.Store) *UserService {
	return &UserService{domainService: domainService, Store: store}
}

func (s *UserService) RegisterUser(ctx context.Context, name, phone, email, password string) (*user.User, error) {
	// 应用层只负责协调，具体逻辑交给领域服务
	err := s.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	newUser, err := s.domainService.RegisterUser(ctx, name, phone, email, password)
	if err != nil {
		txErr := s.Rollback()
		if txErr != nil {
			return nil, err
		}
		if err.Error() == "duplicate_user" {
			return nil, errors.New("unique_violation")
		}
		return nil, err
	}
	err = s.Commit()
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
