package http

import (
	"net/http"

	"github.com/TimRazumov/Technopark-DB/app/service"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	useCase service.UseCase
}

func CreateHandler(router *fasthttprouter.Router, useCase service.UseCase) {
	handler := &Handler{
		useCase: useCase,
	}
	router.GET("/api/service/status", handler.Get)
	router.POST("/api/service/clear", handler.Clear)
}

func (handler *Handler) Get(ctx *fasthttp.RequestCtx) {
	stat := handler.useCase.Get()
	if stat == nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	res, _ := stat.MarshalJSON()
	ctx.SetBody(res)
}

func (handler *Handler) Clear(ctx *fasthttp.RequestCtx) {
	err := handler.useCase.Clear()
	if err != nil {
		ctx.SetStatusCode(err.Code)
		return
	}
}
