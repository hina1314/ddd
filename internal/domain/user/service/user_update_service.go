package service

import (
	"context"
	"github.com/hina1314/ddd/internal/domain/user/entity"
	"github.com/hina1314/ddd/internal/domain/user/repository"
)

// UserUpdateService 用户修改服務
type UserUpdateService struct {
	userRepo repository.UserRepository
}

func NewUserUpdateService(userRepo repository.UserRepository) *UserUpdateService {
	return &UserUpdateService{
		userRepo: userRepo,
	}
}

// UpdateUser 修改用戶信息
func (s *UserUpdateService) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	var err error
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
