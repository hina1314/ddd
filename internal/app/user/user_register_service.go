package user

import (
	"context"
	"database/sql"
	stdErr "errors"
	"study/internal/api/handler/dto"
	"study/internal/domain/user/entity"
	"study/util/errors"
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
		return nil, errors.New(errors.ErrUserAlreadyExists, "user already exists")
	}

	// 开始事务
	tx, txCtx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternalError, "Error txManager.Begin")
	}

	// 确保事务结束时回滚（如果未提交）
	defer func() {
		if err != nil {
			_ = tx.Rollback() // 忽略回滚错误，避免覆盖原始错误
		}
	}()

	user, err = s.userRegisterService.RegisterUser(txCtx, phone, email, password)
	if err != nil {
		return nil, err
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternalError, "Error txManager.Commit")
	}
	return &dto.UserResponse{
		Phone:       user.Phone,
		Username:    user.Username,
		Email:       user.Email.String(),
		AccessToken: "",
		CreatedAt:   user.CreatedAt,
	}, nil
}
