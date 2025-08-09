package handler

import (
	"net/http"

	"github.com/Anhbman/microservice-server-cake/rpc/service"
	"github.com/labstack/echo/v4"
	"github.com/twitchtv/twirp"
)

type (
	createProductRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       int64  `json:"price"`
	}
)

func (h *Handler) CreateProduct(ctx echo.Context) error {
	req := &createProductRequest{}
	if err := ctx.Bind(req); err != nil {
		return err
	}

	if req.Name == "" || req.Price <= 0 {
		return ctx.JSON(http.StatusBadRequest, "name and price are required")
	}

	resp, err := h.serviceClient.CreateProduct(ctx.Request().Context(), &service.CreateProductRequest{
		Name:  req.Name,
		Price: req.Price,
	})
	if err != nil {
		twerr, ok := err.(twirp.Error)
		if ok {
			return ctx.JSON(twirp.ServerHTTPStatusFromErrorCode(twerr.Code()), twerr.Msg())
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, resp)
}
