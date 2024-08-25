package handler

import (
	"client/cmd/utils"
	"net/http"
	"strconv"

	"github.com/Anhbman/microservice-server-cake/rpc/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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

func (h *Handler) Create(ctx echo.Context) error {
	cake := service.CreateCakeRequest{}
	cake.Name = ctx.FormValue("name")
	if cake.Name == "" {
		return ctx.JSON(http.StatusBadRequest, "name is required")
	}
	cake.Description = ctx.FormValue("description")
	if cake.Description == "" {
		return ctx.JSON(http.StatusBadRequest, "description is required")
	}
	cake.Price, _ = strconv.ParseInt(ctx.FormValue("price"), 10, 64)
	if cake.Price == 0 {
		return ctx.JSON(http.StatusBadRequest, "price is required")
	}
	cake.UserId = 5
	file, err := ctx.FormFile("image")
	if err != nil {
		log.Errorf("Invalid data!")
	}

	fileName, err := utils.SaveFile(file)
	if err != nil {
		log.Errorf("Failed to save file!")
	} else {
		cake.ImageUrl = fileName
	}

	newCake, err := h.serviceClient.CreateCake(ctx.Request().Context(), &cake)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, newCake)
}
