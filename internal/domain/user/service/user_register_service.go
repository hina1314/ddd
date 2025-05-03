package service

import (
	"context"
	"study/internal/domain/user/entity"
	"study/internal/domain/user/repository"
	"study/util"
)

// UserRegisterService 用户领域服务
type UserRegisterService struct {
	userRepo repository.UserRepository
}

// NewUserRegisterService 创建用户领域服务
func NewUserRegisterService(userRepo repository.UserRepository) *UserRegisterService {
	return &UserRegisterService{
		userRepo: userRepo,
	}
}

// RegisterUser 注册新用户（包含账户创建）
func (s *UserRegisterService) RegisterUser(ctx context.Context, phone, email, password string) (*entity.User, error) {
	var (
		user *entity.User
		err  error
	)
	passwordHash, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	r := util.NewRandUtil()
	username := r.String(6)
	// 创建用户
	user, err = entity.NewUser(phone, email, username, passwordHash)
	if err != nil {
		return nil, err
	}

	// 保存用户
	if err = s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
