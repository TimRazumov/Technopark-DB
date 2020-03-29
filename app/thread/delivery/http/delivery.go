package http

import (
	"net/http"
	"strconv"

	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/thread"

	"github.com/labstack/echo"
)

type Handler struct {
	useCase thread.UseCase
}

func CreateHandler(router *echo.Echo, useCase thread.UseCase) {
	handler := &Handler{
		useCase: useCase,
	}
	router.POST("api/forum/:slug/create", handler.Create)
	router.GET("api/thread/:slug_or_id/details", handler.Get)
	router.POST("api/thread/:slug_or_id/details", handler.Update)
	router.POST("api/thread/:slug_or_id/vote", handler.UpdateVote)
}

func (handler *Handler) Create(ctx echo.Context) error {
	thrd := models.Thread{Forum: ctx.Param("slug")}
	if err := ctx.Bind(&thrd); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	err := handler.useCase.Create(&thrd)
	if err == nil {
		return ctx.JSON(http.StatusCreated, thrd)
	} else if err.Code == http.StatusConflict {
		existThread := handler.useCase.GetBySlug(thrd.Slug)
		if existThread == nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSON(http.StatusConflict, existThread)
	}
	return ctx.JSON(err.Code, err.Message)
}

func (handler *Handler) Get(ctx echo.Context) error {
	thrdKey := ctx.Param("slug_or_id")
	var thrd *models.Thread
	if id, err := strconv.Atoi(thrdKey); err == nil {
		thrd = handler.useCase.GetByID(id)
	} else {
		thrd = handler.useCase.GetBySlug(thrdKey)
	}
	if thrd == nil {
		return ctx.JSON(http.StatusNotFound, "")
	}
	return ctx.JSON(http.StatusOK, thrd)
}

func (handler *Handler) Update(ctx echo.Context) error {
	thrdKey := ctx.Param("slug_or_id")
	var thrd models.Thread
	if id, err := strconv.Atoi(thrdKey); err == nil {
		thrd.ID = id
	} else {
		thrd.ID = -1
		thrd.Slug = thrdKey
	}
	if err := ctx.Bind(&thrd); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	err := handler.useCase.Update(&thrd)
	if err != nil {
		return ctx.JSON(err.Code, err.Message)
	}
	return ctx.JSON(http.StatusOK, thrd)
}

func (handler *Handler) UpdateVote(ctx echo.Context) error {
	thrdKey := ctx.Param("slug_or_id")
	var vt models.Vote
	if err := ctx.Bind(&vt); err != nil || vt.Voice*vt.Voice != 1 {
		return ctx.NoContent(http.StatusBadRequest)
	}
	thrd := handler.useCase.UpdateVote(thrdKey, vt)
	if thrd == nil {
		return ctx.JSON(http.StatusNotFound, "")
	}
	return ctx.JSON(http.StatusOK, thrd)
}
