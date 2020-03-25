package server

import (
	"fmt"
	"log"

	userRepo "github.com/TimRazumov/Technopark-DB/app/user/repository"

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
	Database: "forums_db",
	User:     "admin",
	Password: "admin1234",
}

func (server *Server) Run() {
	postgeClient, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     config,
			MaxConnections: 55,
		})
	if err != nil {
		log.Fatal(err)
	}
	userRepo.CreateRepository(postgeClient)
	// start
	router := echo.New()
	if err := router.Start(server.GetAddr()); err != nil {
		log.Fatal(err)
	}
}
