package handler

import (
	"net/http"

	"github.com/Anhbman/microservice-server-cake/rpc/service"
	"github.com/labstack/echo/v4"
	"github.com/twitchtv/twirp"
)

func (h *Handler) RegisterUser(ctx echo.Context) error {
	req := &service.RegisterUserRequest{}
	if err := ctx.Bind(req); err != nil {
		return err
	}

	newUser, err := h.serviceClient.RegisterUser(ctx.Request().Context(), req)
	if err != nil {
		twerr, ok := err.(twirp.Error)
		if ok {
			return ctx.JSON(twirp.ServerHTTPStatusFromErrorCode(twerr.Code()), twerr.Msg())
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, newUser)
}
