package app

import (
	"context"
	"study/internal/api/handler/dto"
	"study/internal/assemble"
	entity2 "study/internal/domain/hotel/entity"
	repository2 "study/internal/domain/hotel/repository"
	"study/internal/domain/hotel/service"
	"study/internal/domain/order/entity"
	"study/internal/domain/order/repository"
	"study/internal/domain/user"
	service2 "study/internal/domain/user/service"
	"study/util/errors"
	"time"
)

type OrderService struct {
	orderRepo         repository.OrderRepository
	hotelRepo         repository2.HotelRepository
	userRepo          user.UserRepository
	stockDomainSvc    *service.StockService
	pricingDomainSvc  *service.PricingService
	userPlanDomainSvc service2.UserPlanService
}

func NewOrderService(orderRepo repository.OrderRepository, hotelRepo repository2.HotelRepository, userRepo user.UserRepository, stockDomainSvc *service.StockService, pricingDomainSvc *service.PricingService, userPlanDomainSvc service2.UserPlanService) *OrderService {
	return &OrderService{
		orderRepo:         orderRepo,
		hotelRepo:         hotelRepo,
		userRepo:          userRepo,
		stockDomainSvc:    stockDomainSvc,
		pricingDomainSvc:  pricingDomainSvc,
		userPlanDomainSvc: userPlanDomainSvc,
	}
}

// CreateOrder /service/order_service.go
func (s *OrderService) CreateOrder(ctx context.Context, cmd *assemble.CreateOrderCommand) (*dto.CreateOrderResponse, error) {
	// 参数验证
	if cmd.Start.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, errors.New("xxx", "start date cannot be"+
			" in the past")
	}

	if !cmd.End.After(cmd.Start) {
		return nil, errors.New("xxx", "end date must be after start date")
	}

	// 获取所需数据
	hotelSku, err := s.hotelRepo.FindSkuByID(ctx, cmd.SkuID)
	if err != nil {
		return nil, err
	}

	// 检查预订冲突
	for _, roomContact := range cmd.Contact {
		for _, contact := range roomContact.Guests {
			if err = s.userRepo.CheckBookingConflicts(ctx, cmd.Start, cmd.End, contact.Phone); err != nil {
				return nil, err
			}
		}
	}

	// 计算价格和分配房间
	datePrices, err := s.hotelRepo.GetPrice(ctx, cmd.Start, cmd.End)
	if err != nil {
		return nil, err
	}

	//计算数量
	totalDays := len(datePrices)
	totalNum := totalDays * cmd.Number
	ticketNum := 0

	// 房券支付
	if cmd.PayType == "ticket" {
		ticketNum = totalNum
	}

	totalPrice, err := s.pricingDomainSvc.CalculateTotalPrice(datePrices, cmd.PriceType, cmd.Number)
	if err != nil {
		return nil, err
	}

	// 查询库存
	roomItemIDs, err := s.stockDomainSvc.AllocateRoom(ctx, hotelSku.HotelID, hotelSku.RoomTypeID, cmd.Number, cmd.Start, cmd.End)
	if err != nil {
		return nil, err
	}

	// 创建订单
	order, err := entity.NewOrder(cmd.UserID, hotelSku, totalPrice, totalNum, ticketNum, cmd.Number)
	if err != nil {
		return nil, err
	}

	// 添加房间
	if err = order.AddRoom(hotelSku, roomItemIDs); err != nil {
		return nil, err
	}

	// 持久化
	err = s.orderRepo.Save(ctx, order)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	// 添加用户行程
	if len(cmd.Contact) != len(roomItemIDs) {
		return nil, errors.New("xxxx", "room number mismatch")
	}

	userPlans, err := s.userPlanDomainSvc.BuildPlans(order.ID, roomItemIDs, cmd.Contact, cmd.Start, cmd.End)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.AddUserPlan(ctx, userPlans)
	if err != nil {
		return nil, err
	}
	return &dto.CreateOrderResponse{
		OrderID: order.ID,
		Total:   order.TotalPrice,
	}, nil
}
