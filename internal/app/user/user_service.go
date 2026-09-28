package user

import (
	"context"
	"database/sql"
	stdErrors "errors"
	"github.com/hina1314/ddd/config"
	"github.com/hina1314/ddd/internal/api/handler/dto"
	"github.com/hina1314/ddd/internal/domain/user/repository"
	"github.com/hina1314/ddd/internal/domain/user/service"
	"github.com/hina1314/ddd/token"
	"github.com/hina1314/ddd/util/errors"
)

type UserService struct {
	userRegisterService *service.UserRegisterService
	userLoginService    *service.UserLoginService
	userUpdateService   *service.UserUpdateService
	userRepo            repository.UserRepository
	cfg                 config.Config
	token               token.Maker
}

func NewUserService(
	userRegisterService *service.UserRegisterService,
	userLoginService *service.UserLoginService,
	userUpdateService *service.UserUpdateService,
	userRepo repository.UserRepository,
	cfg config.Config,
	tokenMaker token.Maker,
) *UserService {
	return &UserService{
		userRegisterService: userRegisterService,
		userLoginService:    userLoginService,
		userUpdateService:   userUpdateService,
		userRepo:            userRepo,
		cfg:                 cfg,
		token:               tokenMaker,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, userId int64) (*dto.UserResponse, error) {
	record, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		if stdErrors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(errors.ErrUserNotFound, "user not found")
		}
		return nil, err
	}

	return &dto.UserResponse{
		Phone:     record.Phone,
		Username:  record.Username,
		Email:     record.Email.String(),
		CreatedAt: record.CreatedAt,
	}, nil
}
