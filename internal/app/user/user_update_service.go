package user

import (
	"context"
	"study/internal/api/handler/dto"
	"study/internal/app/assemble"
	"study/internal/domain/user/entity"
	"study/util"
)

func (s *UserService) UpdateUser(ctx context.Context, cmd *assemble.UpdateUserCommand) (*dto.UserResponse, error) {
	var (
		user *entity.User
		err  error
	)

	user, err = s.userRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	user.Phone = cmd.Phone
	user.Email = entity.NewEmail(cmd.Email)
	user.Username = cmd.Username

	if cmd.Password != "" {
		passwordHash, err := util.HashPassword(cmd.Password)
		if err != nil {
			return nil, err
		}
		user.Password = passwordHash
	}

	record, err := s.userUpdateService.UpdateUser(ctx, user)
	if err != nil {
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
