package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tasker/internal/service"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Tasker
// @version 1.0
// @description Tasker make http requests to 3rd-party services

// @host http://localhost:9090
// @BasePath /

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/task", service.CreateTask)
	mux.HandleFunc("/task/", service.GetTaskStatus)
	mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:9090./docs/swagger.json"), //The url pointing to API definition"
	))

	server := &http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}

	log.Println("Server stopped")
}
