package handler

import (
	"net/http"
	"strconv"

	"github.com/Anhbman/microservice-server-cake/rpc/service"
	"github.com/labstack/echo/v4"
)

type (
	createOrderRequest struct {
		Items []orderItem `json:"items"`
	}

	orderItem struct {
		CakeID   uint64 `json:"cake_id"`
		Quantity int64  `json:"quantity"`
	}
)

func (h *Handler) CreateOrder(ctx echo.Context) error {
	userID := userIdFromToken(ctx)
	req := &createOrderRequest{}
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid request format")
	}

	// if req.UserID == 0 {
	// 	return ctx.JSON(http.StatusBadRequest, "User ID is required")
	// }

	if len(req.Items) == 0 {
		return ctx.JSON(http.StatusBadRequest, "Order items are required")
	}

	payload := &service.CreateOrderRequest{
		UserId: userID,
		Items:  make([]*service.OrderItem, len(req.Items)),
	}

	for i, item := range req.Items {
		if item.CakeID == 0 || item.Quantity <= 0 {
			return ctx.JSON(http.StatusBadRequest, "Invalid cake ID or quantity")
		}
		payload.Items[i] = &service.OrderItem{
			CakeId:   item.CakeID,
			Quantity: item.Quantity,
		}
	}

	resp, err := h.serviceClient.CreateOrder(ctx.Request().Context(), payload)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Failed to create order")
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) GetOrderByID(ctx echo.Context) error {
	idStr := ctx.Param("id")
	if idStr == "" {
		return ctx.JSON(http.StatusBadRequest, "Order ID is required")
	}

	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid order ID format")
	}

	order, err := h.serviceClient.GetOrderById(ctx.Request().Context(), &service.GetOrderByIdRequest{Id: orderID})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Failed to retrieve order")
	}

	return ctx.JSON(http.StatusOK, order)
}
