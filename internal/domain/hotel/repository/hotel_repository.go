package repository

import (
	"context"
	"study/internal/domain/hotel/entity"
	"time"
)

type HotelRepository interface {
	FindSkuByID(ctx context.Context, skuID int64) (*entity.HotelSku, error)
	FindAvailableRoomItems(ctx context.Context, hotelID, roomTypeID int64, num int, checkIn, checkOut time.Time) ([]int64, error)
	UpdateRoomStatus(ctx context.Context, roomItemID []int64, checkIn, checkOut time.Time, status string) error
	GetPrice(ctx context.Context, startDate, endDate string) ([]entity.DatePrice, error)
}
