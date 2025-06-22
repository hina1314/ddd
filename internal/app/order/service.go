package order

import (
	"context"
	"study/internal/app/assemble"
	"study/internal/domain/order"
	"study/internal/domain/product"
)

// CreateOrderRequest 创建订单请求
//type CreateOrderRequest struct {
//	UserID int64             `json:"user_id"`
//	Items  []CreateOrderItem `json:"items"`
//}
//
//type CreateOrderItem struct {
//	ProductID int64 `json:"product_id"`
//	Quantity  int   `json:"quantity"`
//}

// AppService 订单应用服务
type AppService struct {
	orderService   *order.Service
	orderRepo      order.Repository
	productService *product.Service
	productRepo    product.Repository
}

func NewAppService(
	orderService *order.Service,
	orderRepo order.Repository,
	productService *product.Service,
	productRepo product.Repository,
) *AppService {
	return &AppService{
		orderService:   orderService,
		orderRepo:      orderRepo,
		productService: productService,
		productRepo:    productRepo,
	}
}

// AddCart 添加到购物车
func (s *AppService) AddCart(ctx context.Context, cmd *assemble.CreateAddCartCommand) (*order.Cart, error) {

	// todo:
	//  1 查询是否sku存在
	//  2 查询用户是否已有该sku
	//  3 如有则创建, 如无则新增

	cart := order.NewCart(cmd.UserID, cmd.SkuID, cmd.Quantity)
	err := s.orderRepo.AddCart(ctx, cart)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

// CreateOrder 创建订单
//func (s *AppService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*order.Order, error) {
//	var orderItems []order.OrderItem
//
//	// 验证商品并构建订单项
//	for _, item := range req.Items {
//		found, err := s.productRepo.GetByID(ctx, item.ProductID)
//		if err != nil {
//			return nil, err
//		}
//
//		if !found.IsAvailable() {
//			return nil, product.ErrProductNotFound
//		}
//
//		orderItem := order.OrderItem{
//			ProductID: item.ProductID,
//			Quantity:  item.Quantity,
//			UnitPrice: found.Price,
//		}
//		orderItems = append(orderItems, orderItem)
//	}
//
//	// 创建订单
//	newOrder, err := s.orderService.CreateOrder(ctx, req.UserID, orderItems)
//	if err != nil {
//		return nil, err
//	}
//
//	//// 预留库存（这里简化处理，实际应该用SKU）
//	//for _, item := range req.Items {
//	//	// 这里假设每个商品只有一个默认SKU
//	//	// 实际项目中需要前端传递具体的SKU ID
//	//	skus, err := s.getProductSKUs(ctx, item.ProductID)
//	//	if err != nil {
//	//		return nil, err
//	//	}
//	//	if len(skus) > 0 {
//	//		err = s.productService.CheckStockAndReserve(ctx, skus[0].ID, item.Quantity)
//	//		if err != nil {
//	//			return nil, err
//	//		}
//	//	}
//	//}
//
//	return newOrder, nil
//}
//
//func (s *AppService) getProductSKUs(ctx context.Context, productID int64) ([]product.SKU, error) {
//	prod, err := s.productRepo.GetByID(ctx, productID)
//	if err != nil {
//		return nil, err
//	}
//	return prod.SKUs, nil
//}
//
//// PayOrder 支付订单
//func (s *ApplicationService) PayOrder(ctx context.Context, orderNo string) error {
//	order, err := s.orderRepo.GetByOrderNo(ctx, orderNo)
//	if err != nil {
//		return err
//	}
//
//	err = order.Pay()
//	if err != nil {
//		return err
//	}
//
//	return s.orderRepo.Update(ctx, order)
//}

// GetOrder 获取订单详情
func (s *AppService) GetOrder(ctx context.Context, orderNo string) (*order.Order, error) {
	return s.orderRepo.GetByOrderNo(ctx, orderNo)
}

// GetUserOrders 获取用户订单列表
func (s *AppService) GetUserOrders(ctx context.Context, userID int64, page, pageSize int) ([]*order.Order, error) {
	offset := (page - 1) * pageSize
	return s.orderRepo.GetByUserID(ctx, userID, offset, pageSize)
}
