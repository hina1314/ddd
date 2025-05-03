package order

import (
	"context"
	"github.com/shopspring/decimal"
	"study/db/model"
	"study/internal/domain/order/entity"
	"study/internal/domain/order/repository"
	"study/util/errors"
	"time"
)

type OrderRepositoryImpl struct {
	db model.TxManager
}

func NewOrderRepository(db model.TxManager) repository.OrderRepository {
	return &OrderRepositoryImpl{
		db: db,
	}
}

func (o *OrderRepositoryImpl) Save(ctx context.Context, order *entity.Order) error {
	q := o.db.Querier(ctx)

	// Save the order
	record, err := q.SaveOrder(ctx, model.SaveOrderParams{
		OrderSn:        order.OrderSn,
		UserID:         order.UserID,
		HotelID:        order.HotelID,
		MerchantID:     order.MerchantID,
		TotalPrice:     order.TotalPrice,
		TotalNumber:    int32(order.TotalNumber),
		TotalPayTicket: int32(order.TotalPayTicket),
		Status:         int16(order.Status),
		CreatedAt:      order.CreatedAt,
		ExpireTime:     order.ExpireTime,
	})
	if err != nil {
		return errors.New("save_order_error", "failed to save order: "+err.Error())
	}

	order.ID = record.ID
	// Save order rooms (if any)
	if len(order.Rooms) > 0 {
		// Prepare arrays for bulk insert
		orderIDs := make([]int64, len(order.Rooms))
		roomTypeIDs := make([]int64, len(order.Rooms))
		roomItemIDs := make([]int64, len(order.Rooms))
		prices := make([]decimal.Decimal, len(order.Rooms))
		statuses := make([]int16, len(order.Rooms))
		createdAts := make([]time.Time, len(order.Rooms))

		for i, room := range order.Rooms {
			orderIDs[i] = order.ID // Link to the saved order
			roomTypeIDs[i] = room.RoomTypeID
			roomItemIDs[i] = room.RoomItemID
			prices[i] = room.Price // Convert decimal.Decimal to float64
			statuses[i] = int16(room.Status)
			createdAts[i] = room.CreatedAt
		}

		err = q.SaveOrderRooms(ctx, model.SaveOrderRoomsParams{
			OrderID:    orderIDs,
			RoomTypeID: roomTypeIDs,
			RoomItemID: roomItemIDs,
			Price:      prices,
			Status:     statuses,
			CreatedAt:  createdAts,
		})
		if err != nil {
			return errors.New("save_order_rooms_error", "failed to save order rooms: "+err.Error())
		}
	}

	return nil
}
