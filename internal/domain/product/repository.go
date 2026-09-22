package product

import "context"

// Repository 商品仓储接口
type Repository interface {
	GetByID(ctx context.Context, id int64) (*Product, error)
	GetSKUByID(ctx context.Context, skuID int64) (*SKU, error)
	UpdateSKUStock(ctx context.Context, skuID int64, stock int) error
	GetProductsWithPagination(ctx context.Context, offset, limit int32) ([]*Product, error)
}
