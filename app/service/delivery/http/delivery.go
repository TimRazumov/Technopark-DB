package http

import (
	"net/http"

	"github.com/TimRazumov/Technopark-DB/app/service"

	"github.com/labstack/echo"
)

type Handler struct {
	useCase service.UseCase
}

func CreateHandler(router *echo.Echo, useCase service.UseCase) {
	handler := &Handler{
		useCase: useCase,
	}
	router.GET("api/service/status", handler.Get)
	router.POST("api/service/clear", handler.Clear)
}

func (handler *Handler) Get(ctx echo.Context) error {
	stat := handler.useCase.Get()
	if stat == nil {
		ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, stat)
}

func (handler *Handler) Clear(ctx echo.Context) error {
	err := handler.useCase.Clear()
	if err != nil {
		ctx.NoContent(err.Code)
	}
	return ctx.NoContent(http.StatusOK)
}
