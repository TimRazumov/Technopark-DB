package http

import (
	"net/http"
	"strconv"

	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/thread"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	useCase thread.UseCase
}

func CreateHandler(router *fasthttprouter.Router, useCase thread.UseCase) {
	handler := &Handler{
		useCase: useCase,
	}
	router.POST("/api/forum/:slug/create", handler.Create)
	router.GET("/api/thread/:slug_or_id/details", handler.Get)
	router.POST("/api/thread/:slug_or_id/details", handler.Update)
	router.POST("/api/thread/:slug_or_id/vote", handler.UpdateVote)
	router.GET("/api/thread/:slug_or_id/posts", handler.GetPosts)
}

func (handler *Handler) Create(ctx *fasthttp.RequestCtx) {
	var thrd models.Thread
	if err := thrd.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	thrd.Forum = ctx.UserValue("slug").(string)
	err := handler.useCase.Create(&thrd)
	if err == nil {
		res, _ := thrd.MarshalJSON()
		ctx.SetStatusCode(http.StatusCreated)
		ctx.SetBody(res)
		return
	} else if err.Code == http.StatusConflict {
		existThread := handler.useCase.GetBySlug(thrd.Slug)
		if existThread == nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		res, _ := existThread.MarshalJSON()
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetBody(res)
		return
	}
	ctx.SetStatusCode(err.Code)
	ctx.SetBody(err.GetMessage())
}

func (handler *Handler) Get(ctx *fasthttp.RequestCtx) {
	thrdKey := ctx.UserValue("slug_or_id").(string)
	var thrd *models.Thread
	if id, err := strconv.Atoi(thrdKey); err == nil {
		thrd = handler.useCase.GetByID(id)
	} else {
		thrd = handler.useCase.GetBySlug(thrdKey)
	}
	if thrd == nil {
		err := models.CreateNotFoundAuthorPost(thrdKey)
		ctx.SetStatusCode(err.Code)
		ctx.SetBody(err.GetMessage())
		return
	}
	res, _ := thrd.MarshalJSON()
	ctx.SetBody(res)
}

func (handler *Handler) Update(ctx *fasthttp.RequestCtx) {
	thrdKey := ctx.UserValue("slug_or_id").(string)
	var thrd models.Thread
	if id, err := strconv.Atoi(thrdKey); err == nil {
		thrd.ID = id
	} else {
		thrd.ID = -1
		thrd.Slug = thrdKey
	}
	if err := thrd.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	err := handler.useCase.Update(&thrd)
	if err != nil {
		ctx.SetStatusCode(err.Code)
		ctx.SetBody(err.GetMessage())
		return
	}
	res, _ := thrd.MarshalJSON()
	ctx.SetBody(res)
}

func (handler *Handler) UpdateVote(ctx *fasthttp.RequestCtx) {
	thrdKey := ctx.UserValue("slug_or_id").(string)
	var vt models.Vote
	if err := vt.UnmarshalJSON(ctx.PostBody()); err != nil || vt.Voice*vt.Voice != 1 {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	thrd := handler.useCase.UpdateVote(thrdKey, vt)
	if thrd == nil {
		err := models.CreateNotFoundAuthorPost(thrdKey)
		ctx.SetStatusCode(err.Code)
		ctx.SetBody(err.GetMessage())
		return
	}
	res, _ := thrd.MarshalJSON()
	ctx.SetBody(res)
}

func (handler *Handler) GetPosts(ctx *fasthttp.RequestCtx) {
	queryString := models.CreateQueryString(ctx.URI().QueryArgs())
	thrdKey := ctx.UserValue("slug_or_id").(string)
	psts := handler.useCase.GetPostsBySlugOrID(thrdKey, queryString)
	if psts == nil {
		err := models.CreateNotFoundForum(thrdKey)
		ctx.SetStatusCode(err.Code)
		ctx.SetBody(err.GetMessage())
		return
	}
	res, _ := psts.MarshalJSON()
	ctx.SetBody(res)
}
