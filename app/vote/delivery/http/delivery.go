package http

import (
	"net/http"

	"github.com/TimRazumov/Technopark-DB/app/forum"
	"github.com/TimRazumov/Technopark-DB/app/models"

	"github.com/labstack/echo"
)

type Handler struct {
	useCase forum.UseCase
}

func CreateHandler(router *echo.Echo, useCase forum.UseCase) {
	handler := &Handler{
		useCase: useCase,
	}
	router.POST("api/forum/create", handler.Create)
	router.GET("api/forum/:slug/details", handler.Get)
}

func (handler *Handler) Create(ctx echo.Context) error {
}

func (handler *Handler) Get(ctx echo.Context) error {
}
