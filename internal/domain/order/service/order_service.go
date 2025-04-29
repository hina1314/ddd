package service

import (
	"fmt"
	"strings"
	"study/internal/domain/hotel/repository"
	"study/internal/domain/hotel/service"
	repository2 "study/internal/domain/order/repository"
	"time"
)

type OrderService struct {
	OrderRepo    repository2.OrderRepository
	userPlanRepo repository.UserPlanRepository
	hotelRepo    repository.HotelRepository
	stockService service.StockService
}

func NewOrderService(userPlanRepo repository.UserPlanRepository, hotelSkuRepo repository.HotelRepository) *OrderService {
	return &OrderService{userPlanRepo: userPlanRepo, hotelRepo: hotelSkuRepo}
}

type Contact struct {
	Name  string
	Phone string
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
