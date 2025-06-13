package handler

import (
	"client/jwt"
	"net/http"

	"github.com/Anhbman/microservice-server-cake/rpc/service"
	"github.com/labstack/echo/v4"
	"github.com/twitchtv/twirp"
)

type (
	loginResponse struct {
		Token string `json:"token"`
		Id    uint64 `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

func (h *Handler) RegisterUser(ctx echo.Context) error {
	req := &service.RegisterUserRequest{}
	if err := ctx.Bind(req); err != nil {
		return err
	}

	if req.GetName() == "" {
		return ctx.JSON(http.StatusBadRequest, "name are required")
	}

	if req.GetEmail() == "" || req.GetPassword() == "" {
		return ctx.JSON(http.StatusBadRequest, "email are required")
	}

	if req.GetPassword() == "" {
		return ctx.JSON(http.StatusBadRequest, "password are required")
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

func (h *Handler) LoginUser(ctx echo.Context) error {
	req := &service.LoginUserRequest{}
	if err := ctx.Bind(req); err != nil {
		return err
	}

	if req.GetEmail() == "" || req.GetPassword() == "" {
		return ctx.JSON(http.StatusBadRequest, "email and password are required")

	}

	user, err := h.serviceClient.LoginUser(ctx.Request().Context(), req)
	if err != nil {
		twerr, ok := err.(twirp.Error)
		if ok {
			return ctx.JSON(twirp.ServerHTTPStatusFromErrorCode(twerr.Code()), twerr.Msg())
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	token := jwt.GenerateJWT(uint(user.GetUser().Id))
	return ctx.JSON(http.StatusOK, loginResponse{
		Token: token,
		Id:    user.GetUser().GetId(),
		Name:  user.GetUser().GetName(),
		Email: user.GetUser().GetEmail(),
	})
}

func (h *Handler) CurrentUser(ctx echo.Context) error {
	id := userIdFromToken(ctx)
	user, err := h.serviceClient.GetUserById(ctx.Request().Context(), &service.GetUserByIdRequest{Id: id})
	if err != nil {
		twerr, ok := err.(twirp.Error)
		if ok {
			return ctx.JSON(twirp.ServerHTTPStatusFromErrorCode(twerr.Code()), twerr.Msg())
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, user)
}

func userIdFromToken(ctx echo.Context) uint64 {
	id, ok := ctx.Get("userId").(uint64)
	if !ok {
		return 0
	}
	return id
}
