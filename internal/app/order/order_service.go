package order

import (
	"study/db/model"
	"study/internal/domain/order/repository"
	repository2 "study/internal/domain/product/repository"
)

type OrderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository2.ProductRepository
	stockRepo   repository
	txManager   model.TxManager
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	productRepo repository2.ProductRepository,
	txManager model.TxManager,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		txManager:   txManager,
	}
}
