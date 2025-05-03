package user

import (
	"context"
	"study/config"
	"study/db/model"
	"study/internal/api/handler/dto"
	"study/internal/domain/user/repository"
	"study/internal/domain/user/service"
	"study/token"
)

type UserService struct {
	userRegisterService *service.UserRegisterService
	userLoginService    *service.UserLoginService
	userRepo            repository.UserRepository
	cfg                 config.Config
	txManager           model.TxManager
	token               token.Maker
}

func NewUserService(
	userRegisterService *service.UserRegisterService,
	userLoginService *service.UserLoginService,
	userRepo repository.UserRepository,
	cfg config.Config,
	txManager model.TxManager,
	tokenMaker token.Maker,
) *UserService {
	return &UserService{
		userRegisterService: userRegisterService,
		userLoginService:    userLoginService,
		userRepo:            userRepo,
		cfg:                 cfg,
		txManager:           txManager,
		token:               tokenMaker,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, userId int64) (*dto.UserResponse, error) {
	record, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		Phone:     record.Phone,
		Username:  record.Username,
		Email:     record.Email.String(),
		CreatedAt: record.CreatedAt,
	}, nil
}
