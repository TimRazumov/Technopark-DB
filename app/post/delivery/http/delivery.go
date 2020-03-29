package http

import (
	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/post"
	"net/http"

	"github.com/labstack/echo"
)

type Handler struct {
	useCase post.UseCase
}

func CreateHandler(router *echo.Echo, useCase post.UseCase) {
	handler := &Handler{
		useCase: useCase,
	}
	router.POST("api/thread/:slug_or_id/create", handler.Create)
	//router.GET("api/post/:id/details", handler.Get)
	//router.POST("api/post/:id/details", handler.Update)
}

func (handler *Handler) Create(ctx echo.Context) error {
	thrdKey := ctx.Param("slug_or_id")
	var posts []models.Post
	if err := ctx.Bind(&posts); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	err := handler.useCase.Create(thrdKey, &posts)
	if err != nil {
		return ctx.JSON(err.Code, err)
	}
	return ctx.JSON(http.StatusCreated, posts)
}

func (handler *Handler) Get(ctx echo.Context) error {
	return ctx.JSON(http.StatusCreated, nil)
}
