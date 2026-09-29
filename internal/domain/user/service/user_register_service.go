package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hina1314/ddd/internal/domain/user/entity"
	"github.com/hina1314/ddd/internal/domain/user/passwordpolicy"
	"github.com/hina1314/ddd/internal/domain/user/repository"
	"github.com/hina1314/ddd/internal/domain/user/usererrors"
	"github.com/hina1314/kit/errors"
	passwordutil "github.com/hina1314/kit/password"
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

// RegisterUser 注册新用户。
func (s *UserRegisterService) RegisterUser(ctx context.Context, phone, email, password string) (*entity.User, error) {
	if passwordpolicy.TooShort(password) {
		return nil, errors.New(usererrors.ErrMinLength, "password must contain at least 8 characters").
			WithParams(map[string]interface{}{"field": "password", "min": passwordpolicy.MinCharacters})
	}
	if passwordpolicy.TooLong(password) {
		return nil, errors.New(usererrors.ErrPasswordTooLong, "password must contain at most 72 bytes").
			WithParams(map[string]interface{}{"field": "password", "max": passwordpolicy.MaxBytes})
	}
	var (
		user *entity.User
		err  error
	)
	passwordHash, err := passwordutil.HashPassword(password)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInvalidInput, "invalid password")
	}

	username := "user_" + uuid.NewString()[:12]
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
