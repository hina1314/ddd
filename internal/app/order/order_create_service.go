package order

import (
	"context"
	"study/internal/api/handler/dto"
	"study/internal/app/assemble"
	"study/internal/domain/order/entity"
	"study/util/errors"
)

func (s *OrderService) CreateOrder(ctx context.Context, cmd *assemble.CreateOrderCommand) (*dto.CreateOrderResponse, error) {
	// 参数验证（领域服务）
	if err := s.OrderService.ValidateBookingDates(cmd.Start, cmd.End); err != nil {
		return nil, err
	}

	// 获取数据
	hotelSku, err := s.hotelRepo.FindSkuByID(ctx, cmd.SkuID)
	if err != nil {
		return nil, err
	}
	datePrices, err := s.hotelRepo.GetPrice(ctx, hotelSku.ID, cmd.Start, cmd.End)
	if err != nil {
		return nil, err
	}

	// 检查预订冲突（领域服务）
	if err = s.userPlanService.CheckBookingConflicts(ctx, cmd.Start, cmd.End, cmd.Contact); err != nil {
		return nil, err
	}

	// 计算数量（领域服务）
	totalNum, ticketNum, err := s.OrderService.CalculateQuantities(datePrices, cmd.Number, cmd.PayType)
	if err != nil {
		return nil, err
	}

	// 计算价格
	totalPrice, err := s.pricingService.CalculateTotalPrice(datePrices, cmd.PriceType, cmd.Number)
	if err != nil {
		return nil, err
	}

	// 创建订单
	order, err := entity.NewOrder(cmd.UserID, hotelSku, totalPrice, totalNum, ticketNum, cmd.Number)
	if err != nil {
		return nil, err
	}

	// 持久化（事务）
	tx, txCtx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternalError, "Error txManager.Begin")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 分配房间（在事务中）
	roomItemIDs, err := s.stockService.AllocateRoom(txCtx, hotelSku.HotelID, hotelSku.RoomTypeID, cmd.Number, cmd.Start, cmd.End)
	if err != nil {
		return nil, err
	}

	// 锁定房态
	if err = s.stockService.LockRoomDates(txCtx, order.ID, order.HotelID, roomItemIDs, datePrices); err != nil {
		return nil, err
	}

	// 添加房间到订单
	if err = order.AddRoom(hotelSku, roomItemIDs); err != nil {
		return nil, err
	}

	// 保存订单
	if err = s.orderRepo.Save(txCtx, order); err != nil {
		return nil, err
	}

	// 创建用户行程
	if err = s.userPlanService.CreateUserPlans(txCtx, order.ID, roomItemIDs, cmd.Contact, cmd.Start, cmd.End); err != nil {
		return nil, err
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternalError, "Error txManager.Commit")
	}

	return &dto.CreateOrderResponse{
		OrderID: order.ID,
		Total:   order.TotalPrice,
	}, nil
}
