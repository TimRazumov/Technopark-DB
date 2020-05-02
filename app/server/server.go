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

	postHandler "github.com/TimRazumov/Technopark-DB/app/post/delivery/http"
	postRepo "github.com/TimRazumov/Technopark-DB/app/post/repository"
	postUseCase "github.com/TimRazumov/Technopark-DB/app/post/usecase"

	serviceHandler "github.com/TimRazumov/Technopark-DB/app/service/delivery/http"
	serviceRepo "github.com/TimRazumov/Technopark-DB/app/service/repository"
	serviceUseCase "github.com/TimRazumov/Technopark-DB/app/service/usecase"

	"github.com/buaazp/fasthttprouter"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
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
	//repo
	postgeClient, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     config,
			MaxConnections: 100,
		})
	if err != nil {
		log.Fatal(err)
	}
	usrRepo := userRepo.CreateRepository(postgeClient)
	frmRepo := forumRepo.CreateRepository(postgeClient)
	thrdRepo := threadRepo.CreateRepository(postgeClient)
	pstRepo := postRepo.CreateRepository(postgeClient)
	servRepo := serviceRepo.CreateRepository(postgeClient)
	// usecase
	usrUseCase := userUseCase.CreateUseCase(usrRepo)
	frmUseCase := forumUseCase.CreateUseCase(usrRepo, frmRepo)
	thrdUseCase := threadUseCase.CreateUseCase(usrRepo, frmRepo, thrdRepo)
	pstUseCase := postUseCase.CreateUseCase(usrRepo, frmRepo, thrdRepo, pstRepo)
	servUseCase := serviceUseCase.CreateUseCase(servRepo)
	// delivery
	router := fasthttprouter.New()
	userHandler.CreateHandler(router, usrUseCase)
	threadHandler.CreateHandler(router, thrdUseCase)
	forumHandler.CreateHandler(router, frmUseCase)
	postHandler.CreateHandler(router, pstUseCase)
	serviceHandler.CreateHandler(router, servUseCase)
	// start
	log.Println("server started on address:", server.GetAddr())
	middleware := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			ctx.Response.Header.Set("Content-Type", "application/json")
			next(ctx)
		}
	}
	err = fasthttp.ListenAndServe(server.GetAddr(), middleware(router.Handler))
	if err != nil {
		log.Fatal(err)
	}
}
