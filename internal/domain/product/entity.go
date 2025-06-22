package product

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"time"
)

// Product 商品聚合根
type Product struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Images      string          `json:"images"`
	SKUs        []SKU           `json:"skus,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

// SKU 商品规格
type SKU struct {
	ID        int64           `json:"id"`
	ProductID int64           `json:"product_id"`
	Name      string          `json:"name"`
	Specs     json.RawMessage `json:"specs"`
	Price     float64         `json:"price"`
	Images    string          `json:"images"`
	Stock     int64           `json:"stock"`
	CreatedAt time.Time       `json:"created_at"`
}

// IsAvailable 检查商品是否可用
func (p *Product) IsAvailable() bool {
	return p.ID > 0 && p.Name != "" && p.Price.IsPositive()
}
