package repository

import (
	"context"
	"study/internal/domain/hotel/entity"
)

type HotelSkuRepository interface {
	GetHotelSku(ctx context.Context, skuId int64) (*entity.HotelSku, error)
	GetPrice(ctx context.Context, startDate, endDate string) ([]entity.DatePrice, error)
}
