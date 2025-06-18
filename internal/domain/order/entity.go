package order

import (
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

// OrderStatus 订单状态
type OrderStatus int

const (
	StatusPending   OrderStatus = iota // 待支付
	StatusPaid                         // 已支付
	StatusCancelled                    // 已取消
	StatusExpired                      // 已过期
)

// Order 订单聚合根
type Order struct {
	ID          int64       `json:"id"`
	OrderNo     string      `json:"order_no"`
	UserID      int64       `json:"user_id"`
	Status      OrderStatus `json:"status"`
	TotalAmount float64     `json:"total_amount"`
	PaidAt      *time.Time  `json:"paid_at,omitempty"`
	ExpireAt    *time.Time  `json:"expire_at,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	Items       []OrderItem `json:"items,omitempty"`
}

// OrderItem 订单项
type OrderItem struct {
	ID        int64           `json:"id"`
	OrderID   int64           `json:"order_id"`
	ProductID int64           `json:"product_id"`
	Quantity  int             `json:"quantity"`
	UnitPrice decimal.Decimal `json:"unit_price"`
}

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderStatus = errors.New("invalid order status")
	ErrOrderExpired       = errors.New("order expired")
)

// CalculateTotal 计算订单总金额
//func (o *Order) CalculateTotal() {
//	total := 0.0
//	for _, item := range o.Items {
//		total += item.UnitPrice * float64(item.Quantity)
//	}
//	o.TotalAmount = total
//}

// CanPay 检查是否可以支付
func (o *Order) CanPay() bool {
	if o.Status != StatusPending {
		return false
	}
	if o.ExpireAt != nil && time.Now().After(*o.ExpireAt) {
		return false
	}
	return true
}

// Pay 支付订单
func (o *Order) Pay() error {
	if !o.CanPay() {
		if o.ExpireAt != nil && time.Now().After(*o.ExpireAt) {
			return ErrOrderExpired
		}
		return ErrInvalidOrderStatus
	}

	now := time.Now()
	o.Status = StatusPaid
	o.PaidAt = &now
	return nil
}

// Cancel 取消订单
func (o *Order) Cancel() error {
	if o.Status != StatusPending {
		return ErrInvalidOrderStatus
	}
	o.Status = StatusCancelled
	return nil
}
