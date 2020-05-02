package http

import (
	"net/http"

	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/user"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	useCase user.UseCase
}

func CreateHandler(router *fasthttprouter.Router, useCase user.UseCase) {
	handler := &Handler{
		useCase: useCase,
	}
	router.POST("/api/user/:nickname/create", handler.Create)
	router.GET("/api/user/:nickname/profile", handler.Get)
	router.POST("/api/user/:nickname/profile", handler.Update)
}

func (handler *Handler) Create(ctx *fasthttp.RequestCtx) {
	var usr models.User
	if err := usr.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	usr.NickName = ctx.UserValue("nickname").(string)
	err := handler.useCase.Create(usr)
	if err == nil {
		res, _ := usr.MarshalJSON()
		ctx.SetStatusCode(http.StatusCreated)
		ctx.SetBody(res)
		return
	}
	var existsUsers models.Users
	existOnNickNameUser := handler.useCase.GetByNickName(usr.NickName)
	if existOnNickNameUser != nil {
		existsUsers = append(existsUsers, *existOnNickNameUser)
	}
	existOnEmailUser := handler.useCase.GetByEmail(usr.Email)
	if existOnEmailUser != nil && (existOnNickNameUser == nil ||
		existOnEmailUser.Email != existOnNickNameUser.Email) {
		existsUsers = append(existsUsers, *existOnEmailUser)
	}
	res, _ := existsUsers.MarshalJSON()
	ctx.SetStatusCode(http.StatusConflict)
	ctx.SetBody(res)
}

func (handler *Handler) Get(ctx *fasthttp.RequestCtx) {
	nickName := ctx.UserValue("nickname").(string)
	usr := handler.useCase.GetByNickName(nickName)
	if usr == nil {
		err := models.CreateNotFoundUser(nickName)
		ctx.SetStatusCode(err.Code)
		ctx.SetBody(err.GetMessage())
		return
	}
	res, _ := usr.MarshalJSON()
	ctx.SetBody(res)
}

func (handler *Handler) Update(ctx *fasthttp.RequestCtx) {
	var usr models.User
	if err := usr.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	usr.NickName = ctx.UserValue("nickname").(string)
	err := handler.useCase.Update(&usr)
	if err != nil {
		ctx.SetStatusCode(err.Code)
		ctx.SetBody(err.GetMessage())
	}
	res, _ := usr.MarshalJSON()
	ctx.SetBody(res)
}
