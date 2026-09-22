package dto

import "github.com/shopspring/decimal"

type Contact struct {
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required,phone"` // e164 是国际手机号格式，也可以自定义正则
}

type RoomContact struct {
	Guests []Contact `json:"guests" validate:"required,dive"` // 每个房间的住客列表
}

type CreateOrderRequest struct {
	SkuID     int64         `json:"sku_id" validate:"required,numeric"`
	StartDate string        `json:"start_date" validate:"required"`
	EndDate   string        `json:"end_date" validate:"required"`
	Number    int           `json:"number" validate:"required,numeric"` // 房间数量
	Contact   []RoomContact `json:"contact" validate:"required,dive"`   // 一维数组，长度 = 房间数量
	PriceType int8          `json:"price_type" validate:"required,oneof=1 2"`
	PayType   string        `json:"pay_type" validate:"required,oneof=wechat ticket"`
}

type CreateOrderResponse struct {
	OrderID int64           `json:"order_id" validate:"required,numeric"`
	Total   decimal.Decimal `json:"total" validate:"required,gt=0"`
}

// AddCartRequest 添加到购物车
type AddCartRequest struct {
	SkuID    int64 `json:"sku_id" validate:"required,numeric"`
	Quantity int32 `json:"quantity" validate:"required,numeric"`
}
