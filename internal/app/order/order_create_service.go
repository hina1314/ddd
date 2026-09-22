package order

//
//import (
//	"context"
//	"study/internal/api/handler/dto"
//	"study/internal/app/assemble"
//	"study/internal/domain/order/entity"
//	"study/util/errors"
//)
//
//func (s *OrderService) CreateOrder(ctx context.Context, cmd *assemble.CreateOrderCommand) (*dto.CreateOrderResponse, error) {
//
//	// 获取数据
//	productSku, err := s.productRepo.FindSkuByID(ctx, cmd.SkuID)
//	if err != nil {
//		return nil, err
//	}
//
//	// 检查库存
//
//	// 创建订单
//	order, err := entity.NewOrder(cmd.UserID, productSku, totalPrice, totalNum, ticketNum, cmd.Number)
//	if err != nil {
//		return nil, err
//	}
//
//	// 持久化（事务）
//	tx, txCtx, err := s.txManager.Begin(ctx)
//	if err != nil {
//		return nil, errors.Wrap(err, errors.ErrInternalError, "Error txManager.Begin")
//	}
//	defer func() {
//		if err != nil {
//			_ = tx.Rollback()
//		}
//	}()
//
//	// 分配房间（在事务中）
//	roomItemIDs, err := s.stockService.AllocateRoom(txCtx, hotelSku.HotelID, hotelSku.RoomTypeID, cmd.Number, cmd.Start, cmd.End)
//	if err != nil {
//		return nil, err
//	}
//
//	// 锁定房态
//	if err = s.stockService.LockRoomDates(txCtx, order.ID, order.HotelID, roomItemIDs, datePrices); err != nil {
//		return nil, err
//	}
//
//	// 保存订单
//	if err = s.orderRepo.Save(txCtx, order); err != nil {
//		return nil, err
//	}
//
//	// 提交事务
//	if err = tx.Commit(); err != nil {
//		return nil, errors.Wrap(err, errors.ErrInternalError, "Error txManager.Commit")
//	}
//
//	return &dto.CreateOrderResponse{
//		OrderID: order.ID,
//		Total:   order.TotalPrice,
//	}, nil
//}
