package entity

import (
	"fmt"
	"github.com/shopspring/decimal"
	"study/internal/domain/hotel/entity"
	"study/util"
	"study/util/errors"
	"time"
)

type OrderStatus string

const (
	OrderStatusInit     OrderStatus = "init"
	OrderStatusPaid     OrderStatus = "paid"
	OrderStatusCheckin  OrderStatus = "checkin"
	OrderStatusCheckout OrderStatus = "checkout"
	OrderStatusRefunded OrderStatus = "refunded"
)

type Order struct {
	ID                   int64
	OrderSn              string
	UserID               int64
	HotelID              int64
	MerchantID           int64
	ProductCategory      string
	TotalPrice           decimal.Decimal
	TotalNumber          int
	TotalPayTicket       int
	TotalRefundedAmount  decimal.Decimal
	TotalRefundedTickets int
	TotalRefundedNumber  int
	AllowRefund          bool
	Status               OrderStatus
	Rooms                []OrderRoom
	ExpireTime           time.Time
	CreatedAt            time.Time
}

// NewOrder 创建订单
func NewOrder(userID int64, sku *entity.HotelSku, totalPrice decimal.Decimal, totalNum, totalTicket, roomNum int) (*Order, error) {
	order := &Order{
		ID:             0,
		OrderSn:        orderSn("LJ"),
		UserID:         userID,
		HotelID:        sku.HotelID,
		MerchantID:     sku.MerchantID,
		TotalPrice:     totalPrice,
		TotalNumber:    totalNum,
		TotalPayTicket: totalTicket,
		AllowRefund:    sku.RefundStatus,
		Status:         OrderStatusInit,
		CreatedAt:      time.Now(),
		ExpireTime:     time.Now().Add(24 * time.Hour),
		Rooms:          make([]OrderRoom, roomNum),
	}
	return order, nil
}

// AddRoom 添加房间子单
func (o *Order) AddRoom(sku *entity.HotelSku, roomItemIDs []int64) error {
	if o.Status != OrderStatusInit {
		return errors.New("xxx", "cannot add room to non-init order")
	}

	//for roomItemID := range roomItemIDs {
	//	o.Rooms = append(o.Rooms, OrderRoom{
	//		ID:         0,
	//		RoomTypeID: sku.RoomTypeID,
	//		RoomItemID: int64(roomItemID),
	//		Price:      sku.SalesPrice,
	//		Status:     OrderRoomStatusInit,
	//		CreatedAt:  time.Now(),
	//	})
	//}
	for i, roomItemID := range roomItemIDs {
		o.Rooms[i] = OrderRoom{
			ID:         0,
			RoomTypeID: sku.RoomTypeID,
			RoomItemID: roomItemID,
			Price:      sku.SalesPrice,
			Status:     OrderRoomStatusInit,
			CreatedAt:  time.Now(),
		}
	}

	return nil
}

func (o *Order) CanRefund() bool {
	return o.AllowRefund
}

func (o *Order) UpdateRefundStats(amount decimal.Decimal, tickets, number int) {
	o.TotalRefundedAmount = o.TotalRefundedAmount.Add(amount)
	o.TotalRefundedTickets += tickets
	o.TotalRefundedNumber += number
}

func (o *Order) MarkAsRefunded() {
	o.Status = OrderStatusRefunded
}

func (o *Order) MarkAsCheckin() {
	o.Status = OrderStatusCheckin
}

func (o *Order) MarkAsCheckout() {
	o.Status = OrderStatusCheckout
}

func orderSn(pre string) string {
	now := time.Now().Format("20060102150405") // Go的时间模板
	r := util.NewRandUtil()
	num := r.Int(1000, 9999)
	return fmt.Sprintf("%v%v%v", pre, now, num)
}
