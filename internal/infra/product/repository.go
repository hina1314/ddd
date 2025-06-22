package product

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"study/db/model"
	"study/internal/domain/product"
	"study/util"
)

type ProductRepository struct {
	db model.TxManager
}

func NewProductRepository(db model.TxManager) product.Repository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*product.Product, error) {
	row, err := r.db.Querier(ctx).GetProductWithSkus(ctx, id)
	if err != nil {
		return nil, err
	}

	var skus, temp []product.SKU // 用于传给 Product
	if err := json.Unmarshal(row.Skus, &temp); err != nil {
		return nil, err
	}
	for _, e := range temp {
		if e.ID > 0 {
			skus = append(skus, e)
		}
	}

	return &product.Product{
		ID:          row.ProductID,
		Name:        row.ProductName,
		Description: util.NullStringToString(row.Description),
		Price:       row.ProductPrice,
		Images:      util.NullStringToString(row.ProductImages),
		SKUs:        skus,
		CreatedAt:   row.ProductCreatedAt,
	}, nil
}

func (r *ProductRepository) GetSKUByID(ctx context.Context, skuID int64) (*product.SKU, error) {
	//sku, err := p.db.GetProductSkuByID(ctx, skuID)
	//if err != nil {
	//	return nil, err
	//}
	//stock, err := p.db.GetSkuStock(ctx, skuID)
	//if err != nil {
	//	return nil, err
	//}
	//return &product.SKU{
	//	ID:        sku.ID,
	//	ProductID: sku.ProductID,
	//	Name:      sku.Name,
	//	Specs:     sku.Specs,
	//	Price:     sku.Price,
	//	Images:    sku.Images,
	//	CreatedAt: sku.CreatedAt,
	//	Stock: product.Stock{
	//		SKUID: sku.ID,
	//		Stock: stock,
	//	},
	//}, nil
	panic("implement me")
}

func (r *ProductRepository) UpdateSKUStock(ctx context.Context, skuID int64, stock int) error {
	//	return r.db.UpdateSkuStock(ctx, model.UpdateSkuStockParams{
	//		SkuID: int64(skuID),
	//		Stock: int32(stock),
	//	})
	//}
	//
	//func (r *ProductRepository) GetProductsWithPagination(ctx context.Context, offset, limit int) ([]*product.Product, error) {
	//	rows, err := p.db.ListProductsWithPagination(ctx, model.ListProductsWithPaginationParams{
	//		Offset: int32(offset),
	//		Limit:  int32(limit),
	//	})
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	var products []*product.Product
	//	for _, row := range rows {
	//		products = append(products, &product.Product{
	//			ID:          row.ID,
	//			Name:        row.Name,
	//			Description: row.Description,
	//			Price:       row.Price,
	//			Images:      row.Images,
	//			CreatedAt:   row.CreatedAt,
	//		})
	//	}
	//	return products, nil
	panic("implement me")
}

func (r *ProductRepository) GetProductsWithPagination(ctx context.Context, offset, limit int32) ([]*product.Product, error) {
	arg := model.ListProductsParams{
		Offset: offset,
		Limit:  limit,
	}
	list, err := r.db.Querier(ctx).ListProducts(ctx, arg)
	if err != nil {
		return nil, err
	}

	var products = make([]*product.Product, len(list))
	for i, row := range list {
		priceStr, _ := row.Price.([]uint8)
		price, err := decimal.NewFromString(string(priceStr))
		if err != nil {
			return nil, fmt.Errorf("invalid price format: %v", err)
		}

		products[i] = &product.Product{
			ID:          row.ID,
			Name:        row.Name,
			Description: util.NullStringToString(row.Description),
			Price:       price,
			Images:      util.NullStringToString(row.Images),
			CreatedAt:   row.CreatedAt,
		}
	}
	return products, nil
}
