package service

import (
	"context"
	"errors"
	"study/internal/domain/hotel/repository"
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
	availableItems, err := s.hotelRepo.FindAvailableRoomItems(ctx, hotelID, roomTypeID, num, checkIn, checkOut)
	if err != nil {
		return nil, err
	}
	if len(availableItems) == 0 {
		return nil, errors.New("no available rooms")
	}

	// 更新房态为 booked
	err = s.hotelRepo.UpdateRoomStatus(ctx, availableItems, checkIn, checkOut, "booked")
	if err != nil {
		return nil, err
	}

	return availableItems, nil
}

// ReleaseRoom 释放房间（订单取消或退款时）
func (s *StockService) ReleaseRoom(ctx context.Context, roomItemIDs []int64, checkIn, checkOut time.Time) error {
	return s.hotelRepo.UpdateRoomStatus(ctx, roomItemIDs, checkIn, checkOut, "available")
}
