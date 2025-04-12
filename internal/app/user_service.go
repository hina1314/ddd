package app

import (
	"context"
	"study/db/model"
	"study/internal/api/handler/dto"
	"study/internal/domain/user"
	"study/token"
	"study/util/errors"
	"time"
)

type UserService struct {
	domainService *user.DomainService
	txManager     model.TxManager
	token         token.Maker
}

func NewUserService(domainService *user.DomainService, txManager model.TxManager, tokenMaker token.Maker) *UserService {
	return &UserService{domainService: domainService, txManager: txManager, token: tokenMaker}
}

func (s *UserService) RegisterUser(ctx context.Context, name, phone, email, password string) (*dto.UserResponse, error) {
	// 应用层只负责协调，具体逻辑交给领域服务
	// 开始事务
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		//return nil, fmt.Errorf("start transaction: %w", err)
		return nil, errors.Wrap(err, errors.ErrInternalError, "Error txManager.Begin")
	}

	// 确保事务结束时回滚（如果未提交）
	defer func() {
		if err != nil {
			_ = tx.Rollback() // 忽略回滚错误，避免覆盖原始错误
		}
	}()

	ctx = context.WithValue(ctx, model.TxKey{}, tx)
	record, err := s.domainService.RegisterUser(ctx, name, phone, email, password)
	if err != nil {
		return nil, err
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternalError, "Error txManager.Commit")
	}
	return &dto.UserResponse{
		Phone:       record.Phone,
		Username:    record.Username,
		Email:       record.Email.String(),
		AccessToken: "",
		CreatedAt:   time.Time{},
	}, nil
}

func (s *UserService) LoginUser(ctx context.Context, phoneOrEmail, password string) (*dto.UserResponse, error) {
	record, err := s.domainService.AuthenticateUser(ctx, phoneOrEmail, password)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.token.CreateToken(uint64(record.ID), 3600)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		Phone:       record.Phone,
		Username:    record.Username,
		Email:       record.Email.String(),
		AccessToken: accessToken,
		CreatedAt:   time.Time{},
	}, nil
}
