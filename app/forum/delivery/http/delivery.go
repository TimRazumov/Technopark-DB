package http

import (
	"net/http"

	"github.com/TimRazumov/Technopark-DB/app/forum"
	"github.com/TimRazumov/Technopark-DB/app/models"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	useCase forum.UseCase
}

func CreateHandler(router *fasthttprouter.Router, useCase forum.UseCase) {
	handler := &Handler{
		useCase: useCase,
	}
	router.POST("/api/forum/:slug", handler.Create)
	router.GET("/api/forum/:slug/details", handler.Get)
	router.GET("/api/forum/:slug/users", handler.GetUsers)
	router.GET("/api/forum/:slug/threads", handler.GetThreads)
}

func (handler *Handler) Create(ctx *fasthttp.RequestCtx) {
	var frm models.Forum
	if err := frm.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	err := handler.useCase.Create(&frm)
	if err == nil {
		res, _ := frm.MarshalJSON()
		ctx.SetStatusCode(http.StatusCreated)
		ctx.SetBody(res)
		return
	} else if err.Code == http.StatusConflict {
		existForum := handler.useCase.GetBySlug(frm.Slug)
		if existForum == nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		res, _ := existForum.MarshalJSON()
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetBody(res)
		return
	}
	ctx.SetStatusCode(err.Code)
	ctx.SetBody(err.GetMessage())
}

func (handler *Handler) Get(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	frm := handler.useCase.GetBySlug(slug)
	if frm == nil {
		err := models.CreateNotFoundForum(slug)
		ctx.SetStatusCode(err.Code)
		ctx.SetBody(err.GetMessage())
		return
	}
	res, _ := frm.MarshalJSON()
	ctx.SetBody(res)
}

func (handler *Handler) GetUsers(ctx *fasthttp.RequestCtx) {
	queryString := models.CreateQueryString(ctx.URI().QueryArgs())
	slug := ctx.UserValue("slug").(string)
	usrs := handler.useCase.GetUsersBySlug(slug, queryString)
	if usrs == nil {
		err := models.CreateNotFoundForum(slug)
		ctx.SetStatusCode(err.Code)
		ctx.SetBody(err.GetMessage())
		return
	}
	res, _ := usrs.MarshalJSON()
	ctx.SetBody(res)
}

func (handler *Handler) GetThreads(ctx *fasthttp.RequestCtx) {
	queryString := models.CreateQueryString(ctx.URI().QueryArgs())
	slug := ctx.UserValue("slug").(string)
	thrds := handler.useCase.GetThreadsBySlug(slug, queryString)
	if thrds == nil {
		err := models.CreateNotFoundForum(slug)
		ctx.SetStatusCode(err.Code)
		ctx.SetBody(err.GetMessage())
		return
	}
	res, _ := thrds.MarshalJSON()
	ctx.SetBody(res)
}
