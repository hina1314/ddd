package service

import (
	"context"
	entity2 "study/internal/domain/hotel/entity"
	"study/internal/domain/hotel/repository"
	"study/util/errors"
	"time"
)

// StockService 管理库存和房态
type StockService struct {
	hotelRepo repository.HotelRepository
}

func NewStockService(hotelRepo repository.HotelRepository) *StockService {
	return &StockService{hotelRepo: hotelRepo}
}

// AllocateRoom 分配一个或多个可用房间
func (s *StockService) AllocateRoom(ctx context.Context, hotelID, roomTypeID int64, num int, checkIn, checkOut time.Time) ([]int64, error) {
	// 查询房态，获取可用 RoomItem
	roomItemIDs, err := s.hotelRepo.FindRoomItemIDsByRoomTypeID(ctx, hotelID, roomTypeID)
	if err != nil {
		return nil, err
	}

	availableItems, err := s.hotelRepo.FindAvailableRoomItems(ctx, roomItemIDs, checkIn, checkOut, num)
	if err != nil {
		return nil, err
	}

	if len(availableItems) == 0 || len(availableItems) != num {
		return nil, errors.New(errors.ErrNoStock, "no available rooms")
	}

	return availableItems, nil
}

func (s *StockService) LockRoomDates(ctx context.Context, orderID, hotelID int64, roomItemIDs []int64, datePrices []entity2.HotelSkuDayPrice) error {
	var roomDates []entity2.HotelRoomDate
	for _, roomItemID := range roomItemIDs {
		for _, price := range datePrices {
			roomDates = append(roomDates, entity2.HotelRoomDate{
				OrderID:    orderID,
				HotelID:    hotelID,
				RoomItemID: roomItemID,
				Date:       price.Date, // 假设 datePrices 包含日期
				Status:     entity2.RoomDateStatusOrderLock,
				CreatedAt:  time.Now(),
			})
		}
	}
	return s.hotelRepo.AddRoomDate(ctx, roomDates)
}

// ReleaseRoom 释放房间（订单取消或退款时）
func (s *StockService) ReleaseRoom(ctx context.Context, roomItemIDs []int64, checkIn, checkOut time.Time) error {
	return nil
}
