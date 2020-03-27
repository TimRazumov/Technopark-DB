package main

import "github.com/TimRazumov/Technopark-DB/app/server"

func main() {
	srv := server.Server{
		IP:   "localhost",
		Port: 5000,
	}
	srv.Run()
}
