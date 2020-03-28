package http

import (
	"net/http"

	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/post"

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
	router.GET("api/forum/:slug_or_id/details", handler.Get)
}

func (handler *Handler) Create(ctx echo.Context) error {
	pstKey := ctx.Param("slug_or_id")
	var posts []models.Post
	if err := ctx.Bind(&posts); err != nil || len(posts) == 0 {
		return ctx.NoContent(http.StatusBadRequest)
	}
	err := handler.useCase.Create(pstKey, &posts)
	if err != nil {
		return ctx.JSON(err.Code, err.Message)
	}
	return ctx.JSON(http.StatusCreated, posts)
}

func (handler *Handler) Get(ctx echo.Context) error {
	return ctx.JSON(http.StatusCreated, nil)
}
