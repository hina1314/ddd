package order

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"
)

// Service 订单领域服务
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GenerateOrderNo 生成订单号
func (s *Service) GenerateOrderNo(userID int64) string {
	timestamp := time.Now().Unix()
	hash := md5.Sum([]byte(fmt.Sprintf("%d_%d", userID, timestamp)))
	return fmt.Sprintf("ORDER_%d_%x", timestamp, hash[:4])
}

// CreateOrder 创建订单
func (s *Service) CreateOrder(ctx context.Context, userID int64, items []OrderItem) (*Order, error) {
	order := &Order{
		OrderNo:   s.GenerateOrderNo(userID),
		UserID:    userID,
		Status:    StatusPending,
		Items:     items,
		CreatedAt: time.Now(),
	}

	// 设置过期时间（30分钟后）
	expireAt := time.Now().Add(30 * time.Minute)
	order.ExpireAt = &expireAt

	// 计算总金额
	order.CalculateTotal()

	err := s.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}
