package product

import (
	"context"
	"errors"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrSKUNotFound       = errors.New("sku not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)

// Service 商品领域服务
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CheckStockAndReserve 检查库存并预留
func (s *Service) CheckStockAndReserve(ctx context.Context, skuID int64, quantity int) error {
	sku, err := s.repo.GetSKUByID(ctx, skuID)
	if err != nil {
		return err
	}
	if sku == nil {
		return ErrSKUNotFound
	}

	if !sku.HasStock(quantity) {
		return ErrInsufficientStock
	}

	// 扣减库存
	newStock := sku.Stock - quantity
	return s.repo.UpdateSKUStock(ctx, skuID, newStock)
}
