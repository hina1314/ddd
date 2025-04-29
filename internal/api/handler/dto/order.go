package dto

type Contact struct {
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required,phone"` // e164 是国际手机号格式，也可以自定义正则
}

type RoomContact struct {
	Guests []Contact `json:"guests" validate:"required,dive"` // 每个房间的住客列表
}

type CreateOrderRequest struct {
	SkuID     string        `json:"sku_id"`
	StartDate string        `json:"start_date"`
	EndDate   string        `json:"end_date"`
	Number    int           `json:"number"`                           // 房间数量
	Contact   []RoomContact `json:"contact" validate:"required,dive"` // 一维数组，长度 = 房间数量
	PriceType int8          `json:"price_type" validate:"required,oneof=1 2"`
	PayType   string        `json:"pay_type" validate:"required,oneof=wechat ticket"`
}
