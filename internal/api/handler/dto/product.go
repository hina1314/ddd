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
	Data  []ProductList `json:"data"`
	Total int           `json:"total" validate:"required,numeric,gt=0"`
}

type ProductList struct {
	ID          int64           `json:"id" validate:"required,numeric"`
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price" validate:"required,gt=0"`
	Image       string          `json:"image" validate:"required"`
}

type Pagination struct {
	Page     int32 `json:"page" validate:"required,numeric,gt=0"`
	PageSize int32 `json:"page_size" validate:"required,numeric,gt=0"`
}
