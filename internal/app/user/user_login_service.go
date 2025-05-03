package user

import (
	"context"
	"study/internal/api/handler/dto"
)

func (s *UserService) LoginUser(ctx context.Context, phone, email, password string) (*dto.UserResponse, error) {
	record, err := s.userLoginService.AuthenticateUser(ctx, phone, email, password)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.token.CreateToken(record.ID, record.Phone, record.Email.String(), s.cfg.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		Phone:       record.Phone,
		Username:    record.Username,
		Email:       record.Email.String(),
		AccessToken: accessToken,
		CreatedAt:   record.CreatedAt,
	}, nil
}
