package repository

import (
	"context"
	"study/internal/domain/hotel/entity"
)

type OrderRepository interface {
	Save(ctx context.Context, order *entity.Order) error
}
