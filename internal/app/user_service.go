package app

import (
	"context"
	"errors"
	"study/db/model"
	"study/internal/domain/user"
)

type UserService struct {
	domainService *user.DomainService
	txManager     model.TxManager
}

func NewUserService(domainService *user.DomainService, txManager model.TxManager) *UserService {
	return &UserService{domainService: domainService, txManager: txManager}
}

func (s *UserService) RegisterUser(ctx context.Context, name, phone, email, password string) (*user.User, error) {
	// 应用层只负责协调，具体逻辑交给领域服务
	// 开始事务
	var newUser *user.User
	err := s.txManager.ExecTx(ctx, func(tx model.Tx) error {
		var err error
		newUser, err = s.domainService.RegisterUser(ctx, name, phone, email, password)
		return err
	})
	if err != nil {
		if err.Error() == "duplicate_user" {
			return nil, errors.New("unique_violation")
		}
		return nil, err
	}

	return newUser, nil
}

//func (s *UserService) RegisterUser(ctx context.Context, name, phone, email, password string) (*user.User, error) {
//	// 应用层只负责协调，具体逻辑交给领域服务
//	// 开始事务
//	tx, err := s.txManager.Begin(ctx)
//	if err != nil {
//		return nil, fmt.Errorf("start transaction: %w", err)
//	}
//
//	// 确保事务结束时回滚（如果未提交）
//	defer func() {
//		if err != nil {
//			_ = tx.Rollback() // 忽略回滚错误，避免覆盖原始错误
//		}
//	}()
//
//	ctx = context.WithValue(ctx, model.TxKey{}, tx)
//	newUser, err := s.domainService.RegisterUser(ctx, name, phone, email, password)
//	if err != nil {
//		if err.Error() == "duplicate_user" {
//			return nil, errors.New("unique_violation")
//		}
//		return nil, err
//	}
//
//	// 提交事务
//	if err = tx.Commit(); err != nil {
//		return nil, fmt.Errorf("commit transaction: %w", err)
//	}
//	return newUser, nil
//}
