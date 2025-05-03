package service

import (
	"fmt"
	"study/internal/domain/hotel/entity"
	"study/internal/domain/hotel/repository"
	"study/internal/domain/hotel/service"
	repository2 "study/internal/domain/order/repository"
	repository3 "study/internal/domain/user/repository"
	"study/util/errors"
	"time"
)

type OrderService struct {
	OrderRepo    repository2.OrderRepository
	userPlanRepo repository3.UserPlanRepository
	hotelRepo    repository.HotelRepository
	stockService *service.StockService
}

func NewOrderService(
	orderRepo repository2.OrderRepository,
	userPlanRepo repository3.UserPlanRepository,
	hotelSkuRepo repository.HotelRepository,
	stockService *service.StockService,
) *OrderService {
	return &OrderService{
		OrderRepo:    orderRepo,
		userPlanRepo: userPlanRepo,
		hotelRepo:    hotelSkuRepo,
		stockService: stockService,
	}
}

func (s *OrderService) ValidateBookingDates(start, end time.Time) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	if start.Before(today) {
		return errors.New(errors.ErrStartDatePast, "start date cannot be in the past")
	}
	if !end.After(start) {
		return errors.New(errors.ErrStartDateDisorder, "end date must be after start date")
	}
	return nil
}

func (s *OrderService) CalculateQuantities(datePrices []entity.HotelSkuDayPrice, roomCount int, payType string) (totalNum, ticketNum int, err error) {
	totalDays := len(datePrices)
	totalNum = totalDays * roomCount
	if payType == "ticket" {
		ticketNum = totalNum
	}
	return totalNum, ticketNum, nil
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
