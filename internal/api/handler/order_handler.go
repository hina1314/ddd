package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"study/internal/api/handler/dto"
	"study/internal/api/response"
	"study/internal/app"
	"study/internal/app/co"
	"study/util/errors"
	"time"
)

type OrderHandler struct {
	res          *response.ResponseHandler
	orderService *app.OrderService
	validator    *validator.Validate
}

func NewOrderHandler(base *response.ResponseHandler, v *validator.Validate) *OrderHandler {
	return &OrderHandler{
		res:       base,
		validator: v,
	}
}

// CreateOrder 处理创建订单请求。
func (h *OrderHandler) CreateOrder(c fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := c.Bind().Body(&req); err != nil {
		return h.res.HandleError(c, err)
	}

	if err := h.validator.Struct(req); err != nil {
		return h.res.HandleError(c, err)
	}

	commandObj := co.CreateOrderCommand{
		SkuID:     req.SkuID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Number:    req.Number,
		PriceType: req.PriceType,
		PayType:   req.PayType,
		Contact:   nil,
	}
	order, err := h.orderService.CreateOrder(c.Context(), commandObj)
	if err != nil {
		return h.res.HandleError(c, err)
	}

	return h.res.SuccessResponse(c, "user.create", order)
}
