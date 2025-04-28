package repository

import (
	"context"
	"study/internal/domain/order/entity"
)

// OrderRepository 只针对 Order 聚合根
type OrderRepository interface {
	Save(ctx context.Context, order *entity.Order) error
	FindByID(ctx context.Context, id int64) (*entity.Order, error)
}
