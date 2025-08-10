package handler

import (
	"client/internal/utils"
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
	cake.UserId = userIdFromToken(ctx)
	// file, err := ctx.FormFile("image")
	// if err != nil {
	// 	log.Errorf("Invalid data!")
	// }

	// fileName, err := utils.SaveFile(file)
	// if err != nil {
	// 	log.Errorf("Failed to save file!", err)
	// } else {
	// 	cake.ImageUrl = fileName
	// }

	newCake, err := h.serviceClient.CreateCake(ctx.Request().Context(), &cake)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, newCake)
}

func (h *Handler) Search(ctx echo.Context) error {
	name := ctx.QueryParam("name")
	pageStr := ctx.QueryParam("page")
	pageSizeStr := ctx.QueryParam("page_size")

	var page, pageSize int
	var err error
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, err)
		}
		if page < 0 {
			page = 0
		}
	} else {
		page = 0
	}

	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, err)
		}
		if pageSize < 1 {
			pageSize = 3
		}
	} else {
		pageSize = 3
	}

	cakes, err := h.serviceClient.SearchCake(ctx.Request().Context(), &service.SearchCakeRequest{
		Name:     name,
		UserId:   int64(userIdFromToken(ctx)),
		Page:     int64(page),
		PageSize: int64(pageSize),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, cakes)
}

func (h *Handler) UpdateByID(ctx echo.Context) error {
	idStr := ctx.Param("id")
	if idStr == "" {
		return ctx.JSON(http.StatusBadRequest, "id is required")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "id must be a number")
	}
	cake := service.Cake{}
	cake.Id = int64(id)
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

	newCake, err := h.serviceClient.UpdateCake(ctx.Request().Context(), &cake)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, newCake)
}

func (h *Handler) GetAllCakes(ctx echo.Context) error {
	cakes, err := h.serviceClient.GetAllCakes(ctx.Request().Context(), &service.GetAllCakesRequest{})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, cakes)
}
