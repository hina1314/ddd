package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
	"study/internal/api/handler/dto"
	"study/internal/api/response"
	"study/internal/app/product"
	"study/util/context"
)

type ProductHandler struct {
	productApp *product.AppService
	res        *response.ResponseHandler
	validator  *validator.Validate
}

func NewProductHandler(res *response.ResponseHandler, productApp *product.AppService, v *validator.Validate) *ProductHandler {
	return &ProductHandler{
		res:        res,
		productApp: productApp,
		validator:  v,
	}
}

// ProductInfo 获取商品详情
func (h *ProductHandler) ProductInfo(c fiber.Ctx) error {
	var (
		req dto.ProductInfoRequest
		err error
	)

	if err = c.Bind().JSON(&req); err != nil {
		return h.res.HandleError(c, err)
	}

	if err = h.validator.Struct(req); err != nil {
		return h.res.HandleError(c, err)
	}

	info, err := h.productApp.GetProduct(c.Context(), req.ID)
	if err != nil {
		return h.res.HandleError(c, err)
	}

	return h.res.Success(c, "product.detail", info)
}

var req dto.CreateOrderRequest
if err := c.Bind().JSON(&req); err != nil {
return h.res.HandleError(c, err)
}

var payload, err = context.GetAuthPayloadFromContext(c.Context())
if err != nil {
return err
}


cmd, err := assemble.NewCreateOrderCommand(req, payload)
if err != nil {
return h.res.HandleError(c, err)
}
res, err := h.orderService.CreateOrder(c.Context(), cmd)
if err != nil {
return h.res.HandleError(c, err)
}

return h.res.Success(c, "order.create", res)

// ProductList 获取商品列表
func (h *ProductHandler) ProductList(c fiber.Ctx) error{
	var (
		req dto.ProductListRequest
		err error
	)
	if err = c.Bind().JSON(&req); err != nil {
		return h.res.HandleError(c, err)
	}

	if err = h.validator.Struct(req); err != nil {
		return h.res.HandleError(c, err)
	}


	products, err := h.productApp.GetProducts(c.Context(), req.Page, req.PageSize)
	if err != nil {
		return h.res.HandleError(c, err)
	}


	var list = make([]dto.ProductList, len(products))
	for _, each := range products {
		list = append(list, dto.ProductList{
			ID:          each.ID,
			Description: each.Description,
			Price:       each.Price,
		})
	}

	res := dto.ProductListResponse{
		Data:  list,
		Total: 0,
	}
	return h.res.Success(c, "product.list", res)
}
