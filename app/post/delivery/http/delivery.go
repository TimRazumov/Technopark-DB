package http

import (
	"net/http"
	"strconv"

	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/post"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	useCase post.UseCase
}

func CreateHandler(router *fasthttprouter.Router, useCase post.UseCase) {
	handler := &Handler{
		useCase: useCase,
	}
	router.POST("/api/thread/:slug_or_id/create", handler.Create)
	router.GET("/api/post/:id/details", handler.Get)
	router.POST("/api/post/:id/details", handler.Update)
}

func (handler *Handler) Create(ctx *fasthttp.RequestCtx) {
	thrdKey := ctx.UserValue("slug_or_id").(string)
	posts := models.Posts{}
	if err := posts.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	err := handler.useCase.Create(thrdKey, &posts)
	if err != nil {
		ctx.SetStatusCode(err.Code)
		ctx.SetBody(err.GetMessage())
		return
	}
	res, _ := posts.MarshalJSON()
	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetBody(res)
}

func (handler *Handler) Get(ctx *fasthttp.RequestCtx) {
	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	options := models.CreateRelated(ctx.URI().QueryArgs())
	pst := handler.useCase.GetByID(id, options)
	if pst == nil {
		err := models.CreateNotFoundThreadPost(id)
		ctx.SetStatusCode(err.Code)
		ctx.SetBody(err.GetMessage())
		return
	}
	res, _ := pst.MarshalJSON()
	ctx.SetBody(res)
}

func (handler *Handler) Update(ctx *fasthttp.RequestCtx) {
	id, er := strconv.Atoi(ctx.UserValue("id").(string))
	if er != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	var pst models.Post
	if err := pst.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	pst.ID = id
	err := handler.useCase.Update(&pst)
	if err != nil {
		ctx.SetStatusCode(err.Code)
		ctx.SetBody(err.GetMessage())
		return
	}
	res, _ := pst.MarshalJSON()
	ctx.SetBody(res)
}
