package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cinemaabyss/microservices/events/config"
	"github.com/cinemaabyss/microservices/events/handler"
	"github.com/cinemaabyss/microservices/events/kafka"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()
	log.Printf("Configs are: %v", cfg)

	producer := kafka.NewProducer(cfg.KafkaBroker)
	defer producer.Close()

	topics := []string{
		cfg.KafkaMovieTopic,
		cfg.KafkaPaymentTopic,
		cfg.KafkaUserTopic,
	}

	consumer := kafka.NewConsumer(cfg.KafkaBroker, topics)
	defer consumer.Close()

	// Start consuming events in a separate goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go consumer.ConsumeEvents(ctx)

	// Create event handler
	eventHandler := handler.NewEventHandler(producer, cfg)

	// Setup routes
	router := mux.NewRouter()
	handler.SetupRoutes(router, eventHandler)

	// Start HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")

}
