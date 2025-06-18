package handler

//import (
//	"github.com/go-playground/validator/v10"
//	"github.com/gofiber/fiber/v3"
//	"study/internal/api/handler/dto"
//	"study/internal/api/response"
//	"study/internal/app/assemble"
//	"study/internal/app/order"
//	"study/util/context"
//)
//
//type OrderHandler struct {
//	res          *response.ResponseHandler
//	orderService *order.OrderService
//	validator    *validator.Validate
//}
//
//func NewOrderHandler(res *response.ResponseHandler, orderService *order.OrderService, v *validator.Validate) *OrderHandler {
//	return &OrderHandler{
//		res:          res,
//		orderService: orderService,
//		validator:    v,
//	}
//}
//
//// CreateOrder 处理创建订单请求。
//func (h *OrderHandler) CreateOrder(c fiber.Ctx) error {
//	var req dto.CreateOrderRequest
//	if err := c.Bind().JSON(&req); err != nil {
//		return h.res.HandleError(c, err)
//	}
//
//	var payload, err = context.GetAuthPayloadFromContext(c.Context())
//	if err != nil {
//		return err
//	}
//
//	if err = h.validator.Struct(req); err != nil {
//		return h.res.HandleError(c, err)
//	}
//
//	cmd, err := assemble.NewCreateOrderCommand(req, payload)
//	if err != nil {
//		return h.res.HandleError(c, err)
//	}
//	res, err := h.orderService.CreateOrder(c.Context(), cmd)
//	if err != nil {
//		return h.res.HandleError(c, err)
//	}
//
//	return h.res.Success(c, "order.create", res)
//}
