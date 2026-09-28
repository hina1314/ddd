package user

import (
	"context"
	"database/sql"
	stdErrors "errors"
	"github.com/hina1314/ddd/internal/api/handler/dto"
	"github.com/hina1314/ddd/internal/app/assemble"
	"github.com/hina1314/ddd/internal/domain/user/entity"
	"github.com/hina1314/ddd/util"
	"github.com/hina1314/ddd/util/errors"
)

func (s *UserService) UpdateUser(ctx context.Context, cmd *assemble.UpdateUserCommand) (*dto.UserResponse, error) {
	if cmd.Phone == nil && cmd.Email == nil && cmd.Username == nil && cmd.Password == nil {
		return nil, errors.New(errors.ErrInvalidInput, "at least one profile field is required")
	}
	if cmd.Username != nil && *cmd.Username == "" {
		return nil, errors.New(errors.ErrInvalidInput, "username must not be empty")
	}
	if cmd.Password != nil && (len(*cmd.Password) < 8 || len(*cmd.Password) > 72) {
		return nil, errors.New(errors.ErrInvalidInput, "password must contain 8 to 72 bytes")
	}
	var (
		user *entity.User
		err  error
	)

	user, err = s.userRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		if stdErrors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(errors.ErrUserNotFound, "user not found")
		}
		return nil, err
	}

	if cmd.Phone != nil {
		user.Phone = *cmd.Phone
	}
	if cmd.Email != nil {
		user.Email = entity.NewEmail(*cmd.Email)
	}
	if cmd.Username != nil {
		user.Username = *cmd.Username
	}
	if user.Phone == "" && user.Email.IsNil() {
		return nil, errors.New(errors.ErrInvalidInput, "phone and email cannot both be empty")
	}

	if cmd.Password != nil {
		passwordHash, err := util.HashPassword(*cmd.Password)
		if err != nil {
			return nil, errors.Wrap(err, errors.ErrInvalidInput, "invalid password")
		}
		user.Password = passwordHash
	}

	record, err := s.userUpdateService.UpdateUser(ctx, user)
	if err != nil {
		if stdErrors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(errors.ErrUserNotFound, "user not found")
		}
		return nil, err
	}

	return &dto.UserResponse{
		Phone:       record.Phone,
		Username:    record.Username,
		Email:       record.Email.String(),
		AccessToken: "",
		CreatedAt:   record.CreatedAt,
	}, nil
}
