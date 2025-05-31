package service

import (
	"context"
	"database/sql"
	stdErr "errors"
	"study/internal/domain/user/entity"
	"study/internal/domain/user/repository"
	"study/util"
	"study/util/errors"
)

// UserLoginService 用户登录务
type UserLoginService struct {
	userRepo repository.UserRepository
}

func NewUserLoginService(userRepo repository.UserRepository) *UserLoginService {
	return &UserLoginService{
		userRepo: userRepo,
	}
}

// AuthenticateUser 用户登录认证
func (s *UserLoginService) AuthenticateUser(ctx context.Context, phone, email, password string) (*entity.User, error) {
	// 查找用户
	var user *entity.User
	var err error

	if phone != "" {
		user, err = s.userRepo.GetByPhone(ctx, phone)
	} else {
		user, err = s.userRepo.GetByEmail(ctx, email)
	}

	if err != nil {
		if stdErr.Is(err, sql.ErrNoRows) {
			return nil, errors.New(errors.ErrUserNotFound, "User not found")
		}
		return nil, err
	}

	err = util.CheckPassword(password, user.Password)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrUserInfoIncorrect, "incorrect password")
	}

	return user, nil
}
