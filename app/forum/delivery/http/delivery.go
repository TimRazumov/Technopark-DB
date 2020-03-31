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
	router.GET("api/forum/:slug/users", handler.GetUsers)
	router.GET("api/forum/:slug/threads", handler.GetThreads)
}

func (handler *Handler) Create(ctx echo.Context) error {
	var frm models.Forum
	if err := ctx.Bind(&frm); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	err := handler.useCase.Create(&frm)
	if err == nil {
		return ctx.JSON(http.StatusCreated, frm)
	} else if err.Code == http.StatusConflict {
		existForum := handler.useCase.GetBySlug(frm.Slug)
		if existForum == nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSON(http.StatusConflict, existForum)
	}
	return ctx.JSON(err.Code, err)
}

func (handler *Handler) Get(ctx echo.Context) error {
	slug := ctx.Param("slug")
	frm := handler.useCase.GetBySlug(slug)
	if frm == nil {
		err := models.CreateNotFoundForum(slug)
		return ctx.JSON(err.Code, err)
	}
	return ctx.JSON(http.StatusOK, frm)
}

func (handler *Handler) GetUsers(ctx echo.Context) error {
	queryString := models.CreateQueryString()
	if err := ctx.Bind(&queryString); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	slug := ctx.Param("slug")
	usrs := handler.useCase.GetUsersBySlug(slug, queryString)
	if usrs == nil {
		err := models.CreateNotFoundForum(slug)
		return ctx.JSON(err.Code, err)
	}
	return ctx.JSON(http.StatusOK, usrs)
}

func (handler *Handler) GetThreads(ctx echo.Context) error {
	queryString := models.CreateQueryString()
	if err := ctx.Bind(&queryString); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	slug := ctx.Param("slug")
	thrds := handler.useCase.GetThreadsBySlug(slug, queryString)
	if thrds == nil {
		err := models.CreateNotFoundForum(slug)
		return ctx.JSON(err.Code, err)
	}
	return ctx.JSON(http.StatusOK, thrds)
}
