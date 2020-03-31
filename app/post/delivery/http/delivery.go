package http

import (
	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/post"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	useCase post.UseCase
}

func CreateHandler(router *echo.Echo, useCase post.UseCase) {
	handler := &Handler{
		useCase: useCase,
	}
	router.POST("api/thread/:slug_or_id/create", handler.Create)
	router.GET("api/post/:id/details", handler.Get)
	router.POST("api/post/:id/details", handler.Update)
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
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	var options models.Related
	query := strings.Split(ctx.QueryParam("related"), ",")
	for _, param := range query {
		if param == "user" {
			options.User = true
			continue
		}
		if param == "forum" {
			options.Forum = true
			continue
		}
		if param == "thread" {
			options.Thread = true
			continue
		}
	}
	pst := handler.useCase.GetByID(id, options)
	if pst == nil {
		return ctx.JSON(http.StatusNotFound, nil)
	}
	return ctx.JSON(http.StatusOK, pst)
}

func (handler *Handler) Update(ctx echo.Context) error {
	id, er := strconv.Atoi(ctx.Param("id"))
	if er != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	var pst models.Post
	if err := ctx.Bind(&pst); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	pst.ID = id
	err := handler.useCase.Update(&pst)
	if err != nil {
		return ctx.JSON(err.Code, err)
	}
	return ctx.JSON(http.StatusOK, pst)
}
