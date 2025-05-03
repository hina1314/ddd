package order

import (
	"study/db/model"
	repository2 "study/internal/domain/hotel/repository"
	"study/internal/domain/hotel/service"
	"study/internal/domain/order/repository"
	service3 "study/internal/domain/order/service"
	repository3 "study/internal/domain/user/repository"
	service2 "study/internal/domain/user/service"
)

type OrderService struct {
	orderRepo       repository.OrderRepository
	hotelRepo       repository2.HotelRepository
	userRepo        repository3.UserRepository
	stockService    *service.StockService
	pricingService  *service.PricingService
	userPlanService *service2.UserPlanService
	OrderService    *service3.OrderService
	txManager       model.TxManager
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	hotelRepo repository2.HotelRepository,
	userRepo repository3.UserRepository,
	stockService *service.StockService,
	pricingService *service.PricingService,
	userPlanService *service2.UserPlanService,
	orderService *service3.OrderService,
	txManager model.TxManager,
) *OrderService {
	return &OrderService{
		orderRepo:       orderRepo,
		hotelRepo:       hotelRepo,
		userRepo:        userRepo,
		stockService:    stockService,
		pricingService:  pricingService,
		userPlanService: userPlanService,
		OrderService:    orderService,
		txManager:       txManager,
	}
}
