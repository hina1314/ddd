package hotel

import (
	"context"
	"github.com/shopspring/decimal"
	"study/db/model"
	"study/internal/domain/hotel/entity"
	"study/internal/domain/hotel/repository"
	"study/internal/infra"
	"study/util/errors"
	"time"
)

type HotelRepositoryImpl struct {
	db model.TxManager
}

func NewHotelRepository(store model.TxManager) repository.HotelRepository {
	return &HotelRepositoryImpl{
		db: store,
	}
}

func (h *HotelRepositoryImpl) FindSkuByID(ctx context.Context, skuID int64) (*entity.HotelSku, error) {
	sku, err := h.db.Querier(ctx).FindSkuByID(ctx, skuID)
	if err != nil {
		if infra.IsNotFoundError(err) {
			return nil, errors.Wrap(err, errors.ErrHotelSkuNotFound, "find_sku_error")
		}
		return nil, err
	}
	return &entity.HotelSku{
		ID:           sku.ID,
		HotelID:      sku.HotelID,
		RoomTypeID:   sku.RoomTypeID,
		SalesPrice:   sku.SalesPrice,
		RefundStatus: sku.RefundStatus,
		MerchantID:   sku.MerchantID,
	}, nil
}

func (h *HotelRepositoryImpl) FindRoomItemIDsByRoomTypeID(ctx context.Context, hotelID, roomTypeID int64) ([]int64, error) {
	q := h.db.Querier(ctx)
	items, err := q.FindRoomItemIDsByRoomTypeID(ctx, model.FindRoomItemIDsByRoomTypeIDParams{
		HotelID:    hotelID,
		RoomTypeID: roomTypeID,
	})
	if err != nil {
		return nil, err
	}

	var ids []int64
	for _, item := range items {
		ids = append(ids, item)
	}
	return ids, nil
}

func (h *HotelRepositoryImpl) FindAvailableRoomItems(ctx context.Context, roomItemIDs []int64, checkIn, checkOut time.Time, limit int) ([]int64, error) {
	q := h.db.Querier(ctx)

	available, err := q.FindAvailableRoomItems(ctx, model.FindAvailableRoomItemsParams{
		RoomItemIds: roomItemIDs,
		StartDate:   checkIn,
		EndDate:     checkOut,
		Limit:       int32(limit),
	})
	if err != nil {
		return nil, err
	}

	var ids []int64
	for _, item := range available {
		ids = append(ids, item)
	}
	return ids, nil
}

func (h *HotelRepositoryImpl) GetPrice(ctx context.Context, skuID int64, start, end time.Time) ([]entity.HotelSkuDayPrice, error) {
	prices, err := h.db.Querier(ctx).GetPrice(ctx, model.GetPriceParams{
		SkuID:     skuID,
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return nil, err
	}

	if len(prices) == 0 {
		return nil, errors.New(errors.ErrNoSkuPrice, "no sku price")
	}
	result := make([]entity.HotelSkuDayPrice, len(prices))
	for i, p := range prices {
		result[i] = entity.HotelSkuDayPrice{
			ID:           p.ID,
			HotelID:      p.HotelID,
			SkuID:        p.SkuID,
			RoomTypeID:   p.RoomTypeID,
			Date:         p.Date,
			MarketPrice:  decimal.Decimal{},
			SalePrice:    decimal.Decimal{},
			TicketPrice:  decimal.Decimal{},
			TicketStatus: p.TicketStatus,
			CreatedAt:    time.Time{},
		}
	}
	return result, nil
}

func (h *HotelRepositoryImpl) AddRoomDate(ctx context.Context, roomDates []entity.HotelRoomDate) error {
	// Prepare arrays for bulk insert
	orderIDs := make([]int64, len(roomDates))
	hotelIDs := make([]int64, len(roomDates))
	roomItemIDs := make([]int64, len(roomDates))
	dates := make([]time.Time, len(roomDates))
	statuses := make([]int16, len(roomDates))
	createdAts := make([]time.Time, len(roomDates))
	updatedAts := make([]time.Time, len(roomDates))

	for i, rd := range roomDates {
		hotelIDs[i] = rd.HotelID
		roomItemIDs[i] = rd.RoomItemID
		dates[i] = rd.Date
		statuses[i] = int16(rd.Status)
		createdAts[i] = rd.CreatedAt
		updatedAts[i] = rd.UpdatedAt
	}

	err := h.db.Querier(ctx).AddRoomDate(ctx, model.AddRoomDateParams{
		OrderID:    orderIDs,
		HotelID:    hotelIDs,
		RoomItemID: roomItemIDs,
		Date:       dates,
		Status:     statuses,
		CreatedAt:  createdAts,
		UpdatedAt:  updatedAts,
	})
	if err != nil {
		return errors.New("add_room_date_error", err.Error())
	}
	return nil
}
