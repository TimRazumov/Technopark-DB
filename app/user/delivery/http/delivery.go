package http

import (
	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/user"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

type Handler struct {
	useCase user.UseCase
}

func CreateHandler(router *echo.Echo, useCase user.UseCase) {
	handler := &Handler{
		useCase: useCase,
	}
	router.POST("/user/:nickname/create", handler.Create)
	router.GET("/user/:nickname/profile", handler.Get)
	router.POST("/user/:nickname/profile", handler.Update)
}

func (handler *Handler) Create(ctx echo.Context) error {
	// TODO: валидация
	var usr models.User
	if err := ctx.Bind(&usr); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	usr.NickName = ctx.Param("nickname")
	log.Println("read user")
	err := handler.useCase.Create(usr)
	if err == nil {
		return ctx.JSON(http.StatusOK, usr)
	}
	var existsUsers []models.User
	existOnNickNameUser := handler.useCase.GetByNickName(usr.NickName)
	if existOnNickNameUser != nil {
		existsUsers = append(existsUsers, *existOnNickNameUser)
	}
	existOnEmailUser := handler.useCase.GetByEmail(usr.Email)
	if existOnEmailUser != nil && (existOnNickNameUser == nil ||
		existOnEmailUser.Email != existOnNickNameUser.Email) {
		existsUsers = append(existsUsers, *existOnEmailUser)
	}
	if len(existsUsers) == 0 {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusConflict, existsUsers)
}

func (handler *Handler) Get(ctx echo.Context) error {
	usr := handler.useCase.GetByNickName(ctx.Param("nickname"))
	if usr == nil {
		return ctx.NoContent(http.StatusNotFound)
	}
	return ctx.JSON(http.StatusOK, usr)
}

func (handler *Handler) Update(ctx echo.Context) error {
	// TODO: валидация
	var usr models.User
	if err := ctx.Bind(&usr); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	usr.NickName = ctx.Param("nickname")
	err := handler.useCase.Update(&usr)
	if err != nil {
		return ctx.NoContent(err.Code)
	}
	return ctx.JSON(http.StatusOK, usr)
}
