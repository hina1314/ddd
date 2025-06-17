package product

import (
	"context"
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
	return s.productRepo.GetByID(ctx, id)
}

// GetProducts 获取商品列表
func (s *AppService) GetProducts(ctx context.Context, page, pageSize int) ([]*product.Product, error) {
	offset := (page - 1) * pageSize
	return s.productRepo.GetProductsWithPagination(ctx, offset, pageSize)
}

// CheckAndReserveStock 检查并预留库存
func (s *AppService) CheckAndReserveStock(ctx context.Context, skuID int64, quantity int) error {
	return s.productService.CheckStockAndReserve(ctx, skuID, quantity)
}
