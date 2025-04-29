package app

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"study/internal/app/co"
	entity2 "study/internal/domain/hotel/entity"
	repository2 "study/internal/domain/hotel/repository"
	"study/internal/domain/hotel/service"
	"study/internal/domain/order/entity"
	"study/internal/domain/order/repository"
	service2 "study/internal/domain/order/service"
	"study/internal/domain/user"
	mCtx "study/util/context"
	"study/util/errors"
	"time"
)

type OrderService struct {
	orderRepo      repository.OrderRepository
	hotelRepo      repository2.HotelRepository
	userRepo       user.UserRepository
	stockDomainSvc service.StockService
	orderDomainSvc service2.OrderService
}

// CreateOrder /service/order_service.go
func (s *OrderService) CreateOrder(ctx context.Context, command co.CreateOrderCommand) error {
	// 参数验证
	var payload, err = mCtx.GetAuthPayloadFromContext(ctx)
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
	// 获取所需数据
	hotelSku, err := s.hotelRepo.FindSkuByID(ctx, skuID)
	if err != nil {
		return err
	}

	// 检查预订冲突
	for _, one := range contact {
		for _, one2 := range one {
			if err = s.userRepo.CheckBookingConflicts(ctx, start, end, one2); err != nil {
				return err
			}
		}
	}

	// 计算价格和分配房间
	datePrices, err := s.hotelRepo.GetPrice(ctx, start, end)
	if err != nil {
		return err
	}

	//计算数量
	totalDays := len(datePrices)
	totalNum := totalDays * roomNum
	ticketNum := 0

	// 房券支付
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

	// 查询库存
	roomItemIDs, err := s.stockDomainSvc.AllocateRoom(ctx, hotelSku.HotelID, hotelSku.RoomTypeID, roomNum, start, end)
	if err != nil {
		return err
	}

	// 创建订单
	order, err := entity.NewOrder(payload.UserId, hotelSku, totalPrice, totalNum, ticketNum, roomNum)
	if err != nil {
		return err
	}

	// 添加房间
	if err = order.AddRoom(hotelSku, roomItemIDs); err != nil {
		return err
	}

	// 持久化
	err = s.orderRepo.Save(ctx, order)
	if err != nil {
		return err
	}

	// 向房态中添加数据
	var roomDates = make([]entity2.HotelRoomDate, totalNum)
	for _, roomItemID := range roomItemIDs {
		for _ = range datePrices {
			roomDates = append(roomDates, entity2.HotelRoomDate{
				OrderID:   order.ID,
				HotelID:   order.HotelID,
				RoomID:    roomItemID,
				Date:      time.Time{},
				Status:    entity2.RoomDateStatusOrderLock,
				CreatedAt: time.Time{},
			})
		}
	}

	err = s.hotelRepo.AddRoomDate(ctx, roomDates)
	if err != nil {
		return err
	}

	// 添加用户行程
	var userPlans = make([]user.UserPlan, 0)
	for _, one := range contact {
		for _, one2 := range one {
			userPlans = append(userPlans, user.UserPlan{
				ID:         0,
				OrderId:    0,
				RoomItemID: 0,
				Phone:      "",
				Name:       one2,
				Start:      time.Time{},
				End:        time.Time{},
				Status:     user.UserPlanStatusInit,
				CreatedAt:  time.Time{},
				UpdatedAt:  time.Time{},
				DeletedAt:  nil,
			})
		}
	}

	err = s.userRepo.AddUserPlan(ctx, userPlans)
	if err != nil {
		return err
	}
	return nil
}

//func (s *OrderService) CreateOrder(userId int64, sku HotelSku, prices []DatePrice, roomNum int, priceType uint8, roomItemIDs []int64) (*entity.Order, error) {
//	// 计算价格（领域规则）
//	totalPrice, ticketNum := s.calculatePrice(prices, roomNum, priceType)
//
//	// 创建订单实体
//	order, err := entity.NewOrder(userId, sku, totalPrice, totalDays*roomNum, ticketNum, roomNum)
//	if err != nil {
//		return nil, err
//	}
//
//	// 添加房间
//	if err := order.AddRoom(sku, roomItemIDs); err != nil {
//		return nil, err
//	}
//
//	return order, nil
//}
