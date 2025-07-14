package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/cinemaabyss/microservices/proxy/configs"
	proxy "github.com/cinemaabyss/microservices/proxy/service"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg *config.Config) *Server {
	// Создаем прокси
	port := config.GetEnv("PORT", "8000")
	monolithURL := config.GetEnv("MONOLITH_URL", "http://localhost:8080")
	moviesURL := config.GetEnv("MOVIES_SERVICE_URL", "http://localhost:8081")

	rp := proxy.NewReverseProxy(monolithURL, moviesURL)
	// Инициализируем роутер
	router := proxy.NewRouter(cfg, rp)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
	}
}

func (s *Server) Run() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		panic(err)
	}
}
