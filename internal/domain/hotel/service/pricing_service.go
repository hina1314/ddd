package service

import (
	"fmt"
	"github.com/shopspring/decimal"
	"study/internal/domain/hotel/entity"
	"study/util/errors"
)

type PricingService struct {
}

func NewPricingService() *PricingService {
	return &PricingService{}
}

// CalculateTotalPrice 根据价格类型和每日价格列表，计算总价
func (s *PricingService) CalculateTotalPrice(datePrices []entity.HotelSkuDayPrice, priceType int8, roomCount int) (decimal.Decimal, error) {
	unitPrice := decimal.Zero

	for _, datePrice := range datePrices {
		switch priceType {
		case 1:
			unitPrice = unitPrice.Add(datePrice.SalePrice)
		case 2:
			if !datePrice.TicketStatus {
				return decimal.Zero, errors.New(errors.ErrTicketNotSupport, fmt.Sprintf("can't use coupon on %v", datePrice.Date))
			}
			unitPrice = unitPrice.Add(datePrice.TicketPrice)
		default:
			return decimal.Zero, errors.New("invalid_price_type", "unsupported price type")
		}
	}

	totalRooms := decimal.NewFromInt(int64(roomCount))
	totalPrice := unitPrice.Mul(totalRooms)
	return totalPrice, nil
}
