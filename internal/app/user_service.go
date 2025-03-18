package app

import (
	"context"

	"study/internal/domain/user"
)

type UserService struct {
	repo    user.UserRepository
	service *user.DomainService
}

func NewUserService(userRepo user.UserRepository, userAccountRepo user.UserAccountRepository) *UserService {
	return &UserService{
		repo:    userRepo,
		service: user.NewDomainService(userRepo, userAccountRepo),
	}
}

func (s *UserService) GetUserById(ctx context.Context, id int64) (*user.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) RegisterUser(ctx context.Context, name, phone, email, password string) (*user.User, error) {
	newUser, err := s.service.RegisterUser(ctx, name, phone, email, password)
	if err != nil {
		return nil, err
	}
	err = s.repo.Save(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return newUser, err
}
