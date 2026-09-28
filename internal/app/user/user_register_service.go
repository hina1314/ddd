package user

import (
	"context"
	"database/sql"
	stdErr "errors"
	"github.com/hina1314/ddd/internal/api/handler/dto"
	"github.com/hina1314/ddd/internal/domain/user/entity"
	"github.com/hina1314/ddd/internal/domain/user/usererrors"
	"github.com/hina1314/kit/errors"
)

func (s *UserService) RegisterUser(ctx context.Context, phone, email, password string) (*dto.UserResponse, error) {
	var (
		user *entity.User
		err  error
	)

	if phone != "" {
		user, err = s.userRepo.GetByPhone(ctx, phone)
	} else {
		user, err = s.userRepo.GetByEmail(ctx, email)
	}

	if err != nil && !stdErr.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if user != nil {
		return nil, errors.New(usererrors.ErrUserAlreadyExists, "user already exists")
	}

	user, err = s.userRegisterService.RegisterUser(ctx, phone, email, password)
	if err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		Phone:       user.Phone,
		Username:    user.Username,
		Email:       user.Email.String(),
		AccessToken: "",
		CreatedAt:   user.CreatedAt,
	}, nil
}
