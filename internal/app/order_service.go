package app

import (
	"context"
	repository2 "study/internal/domain/hotel/repository"
	"study/internal/domain/hotel/service"
	"study/internal/domain/order/repository"
	service2 "study/internal/domain/order/service"
)

type OrderService struct {
	orderRepo      repository.OrderRepository
	hotelRepo      repository2.HotelRepository
	stockDomainSvc service.StockService
	orderDomainSvc service2.OrderService
}

// CreateOrder /service/order_service.go
func (s *OrderService) CreateOrder(ctx context.Context, cmd CreateOrderCommand) error {
	// 参数验证
	if err := s.validateCreateOrderParams(cmd); err != nil {
		return err
	}

	// 获取所需数据
	hotelSku, err := s.hotelRepo.FindSkuByID(ctx, cmd.SkuId)
	if err != nil {
		return err
	}

	// 检查预订冲突
	if err := s.checkBookingConflicts(cmd.Contact, cmd.StartDate, cmd.EndDate); err != nil {
		return err
	}

	// 计算价格和分配房间
	prices, err := s.hotelRepo.GetPrice(ctx, cmd.StartDate, cmd.EndDate)
	roomItemIDs, err := s.stockDomainSvc.AllocateRoom(ctx, hotelSku.HotelID, hotelSku.RoomTypeID, roomNum, start, end)

	// 调用领域服务创建订单
	orderService := s.orderDomainService
	order, err := orderService.CreateOrder(userId, hotelSku, prices, cmd.RoomNum, cmd.PriceType, roomItemIDs)

	// 持久化
	return s.orderRepo.Save(order)
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
