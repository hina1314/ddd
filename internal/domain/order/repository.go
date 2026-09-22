package order

import "context"

// Repository 订单仓储接口
type Repository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id int64) (*Order, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*Order, error)
	Update(ctx context.Context, order *Order) error
	GetByUserID(ctx context.Context, userID int64, offset, limit int) ([]*Order, error)
	AddCart(ctx context.Context, cart *Cart) error
}
