package server

import (
	"fmt"
	"log"

	userHandler "github.com/TimRazumov/Technopark-DB/app/user/delivery/http"
	userRepo "github.com/TimRazumov/Technopark-DB/app/user/repository"
	userUseCase "github.com/TimRazumov/Technopark-DB/app/user/usecase"

	forumHandler "github.com/TimRazumov/Technopark-DB/app/forum/delivery/http"
	forumRepo "github.com/TimRazumov/Technopark-DB/app/forum/repository"
	forumUseCase "github.com/TimRazumov/Technopark-DB/app/forum/usecase"

	threadHandler "github.com/TimRazumov/Technopark-DB/app/thread/delivery/http"
	threadRepo "github.com/TimRazumov/Technopark-DB/app/thread/repository"
	threadUseCase "github.com/TimRazumov/Technopark-DB/app/thread/usecase"

	serviceHandler "github.com/TimRazumov/Technopark-DB/app/service/delivery/http"
	serviceRepo "github.com/TimRazumov/Technopark-DB/app/service/repository"
	serviceUseCase "github.com/TimRazumov/Technopark-DB/app/service/usecase"

	"github.com/jackc/pgx"
	"github.com/labstack/echo"
)

type Server struct {
	IP   string
	Port uint
}

func (server *Server) GetAddr() string {
	return fmt.Sprintf("%s:%d", server.IP, server.Port)
}

// TODO: в отдельный файл
var config = pgx.ConnConfig{
	Host:     "localhost",
	Port:     5432,
	Database: "forum_db",
	User:     "forum_user",
	Password: "forum1234",
}

func (server *Server) Run() {
	router := echo.New()
	//repo
	postgeClient, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     config,
			MaxConnections: 55,
		})
	if err != nil {
		log.Fatal(err)
	}
	usrRepo := userRepo.CreateRepository(postgeClient)
	frmRepo := forumRepo.CreateRepository(postgeClient)
	thrdRepo := threadRepo.CreateRepository(postgeClient)
	servRepo := serviceRepo.CreateRepository(postgeClient)
	// usecase
	usrUseCase := userUseCase.CreateUseCase(usrRepo)
	frmUseCase := forumUseCase.CreateUseCase(usrRepo, frmRepo)
	thrdUseCase := threadUseCase.CreateUseCase(usrRepo, frmRepo, thrdRepo)
	servUseCase := serviceUseCase.CreateUseCase(servRepo)
	// delivery
	userHandler.CreateHandler(router, usrUseCase)
	forumHandler.CreateHandler(router, frmUseCase)
	threadHandler.CreateHandler(router, thrdUseCase)
	serviceHandler.CreateHandler(router, servUseCase)
	// start
	if err := router.Start(server.GetAddr()); err != nil {
		log.Fatal(err)
	}
}
