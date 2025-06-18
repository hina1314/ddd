package product

import (
	"context"
	"database/sql"
	"study/db/model"
	"study/internal/domain/product"
)

type ProductRepository struct {
	db model.TxManager
}

func NewProductRepository(db model.TxManager) product.Repository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*product.Product, error) {
	query := `
		SELECT id, name, description, price, images, created_at 
		FROM product 
		WHERE id = $1
	`

	var p product.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.Images, &p.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, product.ErrProductNotFound
		}
		return nil, err
	}

	// 获取SKU信息
	skus, err := r.getSKUsByProductID(ctx, id)
	if err != nil {
		return nil, err
	}
	p.SKUs = skus

	return &p, nil
}

func (r *ProductRepository) getSKUsByProductID(ctx context.Context, productID int64) ([]product.SKU, error) {
	query := `
		SELECT ps.id, ps.product_id, ps.name, ps.specs, ps.price, ps.images, ps.created_at,
		       COALESCE(pss.stock, 0) as stock
		FROM product_sku ps
		LEFT JOIN product_sku_stock pss ON ps.id = pss.sku_id
		WHERE ps.product_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skus []product.SKU
	for rows.Next() {
		var sku product.SKU
		err := rows.Scan(
			&sku.ID, &sku.ProductID, &sku.Name, &sku.Specs, &sku.Price,
			&sku.Images, &sku.CreatedAt, &sku.Stock,
		)
		if err != nil {
			return nil, err
		}
		skus = append(skus, sku)
	}

	return skus, nil
}

func (r *ProductRepository) GetSKUByID(ctx context.Context, skuID int64) (*product.SKU, error) {
	query := `
		SELECT ps.id, ps.product_id, ps.name, ps.specs, ps.price, ps.images, ps.created_at,
		       COALESCE(pss.stock, 0) as stock
		FROM product_sku ps
		LEFT JOIN product_sku_stock pss ON ps.id = pss.sku_id
		WHERE ps.id = $1
	`

	var sku product.SKU
	err := r.db.QueryRowContext(ctx, query, skuID).Scan(
		&sku.ID, &sku.ProductID, &sku.Name, &sku.Specs, &sku.Price,
		&sku.Images, &sku.CreatedAt, &sku.Stock,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, product.ErrSKUNotFound
		}
		return nil, err
	}

	return &sku, nil
}

func (r *ProductRepository) UpdateSKUStock(ctx context.Context, skuID int64, stock int) error {
	query := `
		INSERT INTO product_sku_stock (sku_id, stock) 
		VALUES ($1, $2)
		ON CONFLICT (sku_id) 
		DO UPDATE SET stock = $2
	`

	_, err := r.db.ExecContext(ctx, query, skuID, stock)
	return err
}

func (r *ProductRepository) GetProductsWithPagination(ctx context.Context, offset, limit int) ([]*product.Product, error) {
	query := `
		SELECT id, name, description, price, images, created_at 
		FROM product 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*product.Product
	for rows.Next() {
		var p product.Product
		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Price, &p.Images, &p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	return products, nil
}
