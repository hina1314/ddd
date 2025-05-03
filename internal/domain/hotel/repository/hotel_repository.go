package repository

import (
	"context"
	"study/internal/domain/hotel/entity"
	"time"
)

type HotelRepository interface {
	FindSkuByID(ctx context.Context, skuID int64) (*entity.HotelSku, error)
	FindRoomItemIDsByRoomTypeID(ctx context.Context, hotelID, roomTypeID int64) ([]int64, error)
	FindAvailableRoomItems(ctx context.Context, roomItemIDs []int64, checkIn, checkOut time.Time, num int) ([]int64, error)
	GetPrice(ctx context.Context, skuID int64, start, end time.Time) ([]entity.HotelSkuDayPrice, error)
	AddRoomDate(ctx context.Context, roomDates []entity.HotelRoomDate) error
}
