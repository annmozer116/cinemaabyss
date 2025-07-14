package main

import (
	config "github.com/cinemaabyss/microservices/proxy/configs"
	"github.com/cinemaabyss/microservices/proxy/server"
)

func main() {

	cfg := config.Load()

	srv := server.New(cfg)
	srv.Run()
}
