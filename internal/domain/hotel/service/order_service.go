package service

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
	"study/internal/domain/hotel/entity"
	"study/internal/domain/hotel/repository"
	mCtx "study/util/context"
	"study/util/errors"
	"time"
)

type OrderService struct {
	OrderRepo    repository.OrderRepository
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

func (o *OrderService) CreateOrder(ctx context.Context, skuId int64, startDate, endDate string, roomNum int, priceType uint8, payType string, contact [][]Contact) error {
	payload, err := mCtx.GetAuthPayloadFromContext(ctx)
	if err != nil {
		return err
	}

	start, _ := time.Parse("2006-01-02", startDate)
	if start.Before(time.Now().Truncate(24 * time.Hour)) {
		return errors.New("xxx", "start date cannot be in the past")
	}

	end, _ := time.Parse("2006-01-02", endDate)
	if !end.After(start) {
		return errors.New("xxx", "end date must be after start date")
	}

	if len(contact) != int(roomNum) || roomNum == 0 {
		return errors.New("xxx", "contact doesn't match room number")
	}

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

	hotelSku, err := o.hotelSkuRepo.GetHotelSku(ctx, skuId)
	if err != nil {
		return err
	}
	datePrices, err := o.hotelSkuRepo.GetPrice(ctx, startDate, endDate)
	if err != nil {
		return err
	}

	//计算数量
	totalDays := len(datePrices)
	totalNum := totalDays * int(roomNum)
	ticketNum := 0

	if priceType == 2 {
		ticketNum = totalNum
	}

	//计算房价单价
	unitPrice := decimal.Zero
	for _, datePrice := range datePrices {
		switch priceType {
		case 1:
			unitPrice = unitPrice.Add(datePrice.SalePrice)
		case 2:
			if !datePrice.TicketStatus {
				return errors.New("xxx", fmt.Sprintf("can't use coupon on %v", datePrice.Date))
			}
			unitPrice = unitPrice.Add(datePrice.TicketPrice)
			break
		default:
			return errors.New("xxxx", "invalid priceType")

		}
	}

	totalPrice := unitPrice.Mul(decimal.New(int64(roomNum), 2))

	order, err := entity.NewOrder(payload.UserId, hotelSku, totalPrice, totalNum, ticketNum, roomNum)

	if err != nil {
		return err
	}

	err = o.OrderRepo.SaveOrderRoom(ctx, order.Id)

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
