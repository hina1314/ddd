package product

import (
	"context"
	"database/sql"
	"errors"
	"study/internal/domain/product"
)

// AppService  商品应用服务
type AppService struct {
	productService *product.Service
	productRepo    product.Repository
}

func NewAppService(productService *product.Service, productRepo product.Repository) *AppService {
	return &AppService{
		productService: productService,
		productRepo:    productRepo,
	}
}

// GetProduct 获取商品详情
func (s *AppService) GetProduct(ctx context.Context, id int64) (*product.Product, error) {
	found, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, product.ErrProductNotFound
		}
		return nil, err
	}
	return found, nil
}

// GetProducts 获取商品列表
func (s *AppService) GetProducts(ctx context.Context, page, pageSize int32) ([]*product.Product, error) {
	offset := (page - 1) * pageSize
	return s.productRepo.GetProductsWithPagination(ctx, offset, pageSize)
}

// CheckAndReserveStock 检查并预留库存
//func (s *AppService) CheckAndReserveStock(ctx context.Context, skuID int64, quantity int) error {
//	return s.productService.CheckStockAndReserve(ctx, skuID, quantity)
//}
