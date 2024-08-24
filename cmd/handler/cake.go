package handler

import (
	"net/http"
	"strconv"

	"github.com/Anhbman/microservice-server-cake/rpc/service"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetByID(ctx echo.Context) error {
	idStr := ctx.Param("id")
	if idStr == "" {
		return ctx.JSON(http.StatusBadRequest, "id is required")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "id must be a number")
	}
	cake, err := h.serviceClient.GetCakeById(ctx.Request().Context(), &service.GetCakeByIdRequest{Id: int64(id)})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, cake)
}
