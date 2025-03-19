package app

import (
	"context"
	"errors"
	"study/internal/domain/user"
)

type UserService struct {
	domainService *user.DomainService
}

func NewUserService(domainService *user.DomainService) *UserService {
	return &UserService{domainService: domainService}
}

//func (s *UserService) GetUserById(ctx context.Context, id int64) (*user.User, error) {
//	return s.domainService.GetByID(ctx, id)
//}

func (s *UserService) RegisterUser(ctx context.Context, name, phone, email, password string) (*user.User, error) {
	// 应用层只负责协调，具体逻辑交给领域服务
	newUser, err := s.domainService.RegisterUser(ctx, name, phone, email, password)
	if err != nil {
		if err.Error() == "duplicate_user" {
			return nil, errors.New("unique_violation")
		}
		return nil, err
	}
	return newUser, nil
}
