package service

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
	"study/internal/domain/hotel/repository"
	"study/util/errors"
	"time"
)

type OrderService struct {
	userPlanRepo repository.UserPlanRepository
	hotelSkuRepo repository.HotelSkuRepository
}

func NewOrderService(userPlanRepo repository.UserPlanRepository, hotelSkuRepo repository.HotelSkuRepository) *OrderService {
	return &OrderService{userPlanRepo: userPlanRepo, hotelSkuRepo: hotelSkuRepo}
}

type Contact struct {
	Name  string
	Phone string
}

func (o *OrderService) CreateOrder(ctx context.Context, skuId int64, startDate, endDate string, number int, priceType, payType string, contact [][]Contact) error {
	start, _ := time.Parse("2006-01-02", startDate)
	if start.Before(time.Now().Truncate(24 * time.Hour)) {
		return errors.New("xxx", "start date cannot be in the past")
	}

	end, _ := time.Parse("2006-01-02", endDate)
	if !end.After(start) {
		return errors.New("xxx", "end date must be after start date")
	}

	if len(contact) != number || number == 0 {
		return errors.New("xxx", "contact doesn't match room number")
	}

	//dates, err := GetDatesNotLastDay(startDate, endDate)
	//if err != nil {
	//	return err
	//}
	// Clean phone numbers and validate contacts
	for _, roomContacts := range contact {
		for _, each := range roomContacts {
			if each.Name == "" || each.Phone == "" {
				return errors.New("xxx", "contact name and phone are required")
			}
			each.Phone = removePhoneSpaces(each.Phone)
			count, err := o.userPlanRepo.CountConflictingPlans(each.Phone, startDate, endDate)
			if err != nil {
				return err
			}
			if count > 0 {
				return errors.New("xxxx", "booking conflict: "+each.Name+" already has a booking")
			}
		}
	}

	_, err := o.hotelSkuRepo.GetHotelSku(ctx, skuId)
	if err != nil {
		return err
	}
	datePrices, err := o.hotelSkuRepo.GetPrice(ctx, startDate, endDate)
	if err != nil {
		return err
	}

	totalPrice := decimal.Zero
	for _, datePrice := range datePrices {
		totalPrice = totalPrice.Add(datePrice.SalePrice)
	}
	return nil
}

func removePhoneSpaces(phone string) string {
	return strings.ReplaceAll(phone, " ", "")
}

// GetDatesNotLastDay 返回从startDate到endDate（不包括最后一天）的日期列表
func GetDatesNotLastDay(startDate, endDate string) ([]string, error) {
	// 解析输入日期
	layout := "2006-01-02" // Go的时间格式
	startTime, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %v", err)
	}
	endTime, err := time.Parse(layout, endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %v", err)
	}

	// 初始化结果切片
	var dates []string

	// 循环添加日期
	for startTime.Before(endTime) {
		dates = append(dates, startTime.Format(layout))
		startTime = startTime.Add(24 * time.Hour)
	}

	return dates, nil
}
