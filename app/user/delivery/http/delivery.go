package http

import (
	"net/http"
	"regexp"

	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/user"
	"github.com/labstack/echo"
)

var nickNamePattern = regexp.MustCompile("^[A-Za-z0-9_.]+$")

type Handler struct {
	useCase user.UseCase
}

func CreateHandler(router *echo.Echo, useCase user.UseCase) {
	handler := &Handler{
		useCase: useCase,
	}
	router.POST("api/user/:nickname/create", handler.Create)
	router.GET("api/user/:nickname/profile", handler.Get)
	router.POST("api/user/:nickname/profile", handler.Update)
}

func (handler *Handler) Create(ctx echo.Context) error {
	usr := models.User{NickName: ctx.Param("nickname")}
	if err := ctx.Bind(&usr); err != nil || !nickNamePattern.MatchString(usr.NickName) {
		return ctx.NoContent(http.StatusBadRequest)
	}
	err := handler.useCase.Create(usr)
	if err == nil {
		return ctx.JSON(http.StatusCreated, usr)
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
	nickName := ctx.Param("nickname")
	err := models.CreateNotFoundUser(nickName)
	if !nickNamePattern.MatchString(nickName) {
		return ctx.JSON(err.Code, err)
	}
	usr := handler.useCase.GetByNickName(nickName)
	if usr == nil {
		return ctx.JSON(err.Code, err)
	}
	return ctx.JSON(http.StatusOK, usr)
}

func (handler *Handler) Update(ctx echo.Context) error {
	nickName := ctx.Param("nickname")
	var usr models.User
	if err := ctx.Bind(&usr); !nickNamePattern.MatchString(nickName) || err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	usr.NickName = nickName
	err := handler.useCase.Update(&usr)
	if err != nil {
		return ctx.JSON(err.Code, err)
	}
	return ctx.JSON(http.StatusOK, usr)
}
