package dto

import "github.com/shopspring/decimal"

type ProductInfoRequest struct {
	ID int64 `json:"id" validate:"required,numeric"`
}

type ProductInfoResponse struct {
	OrderID int64           `json:"order_id" validate:"required,numeric"`
	Total   decimal.Decimal `json:"total" validate:"required,gt=0"`
}

type ProductListRequest struct {
	Pagination
}

type ProductListResponse struct {
	Data  []ProductList
	Total int `json:"total" validate:"required,numeric,gt=0"`
}

type ProductList struct {
	ID          int64           `json:"id" validate:"required,numeric"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price" validate:"required,gt=0"`
}

type Pagination struct {
	Page     int `json:"page" validate:"required,numeric,gt=0"`
	PageSize int `json:"page_size" validate:"required,numeric,gt=0"`
}
