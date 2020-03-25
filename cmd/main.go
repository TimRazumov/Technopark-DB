package main

import "github.com/TimRazumov/Technopark-DB/app/server"

func main() {
	srv := server.Server{
		IP:   "127.0.0.1",
		Port: 8080,
	}
	srv.Run()
}
