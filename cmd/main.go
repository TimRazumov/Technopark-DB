package main

import (
	"github.com/TimRazumov/Technopark-DB/app/server"
)

func main() {
	srv := server.Server{
		IP:   "0.0.0.0", //"localhost"
		Port: 5000,
	}
	srv.Run()
}
