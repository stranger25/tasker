package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tasker/internal/repository"
	"tasker/internal/service"
	"time"
)

// @title Tasker
// @version 1.0
// @description Tasker make http requests to 3rd-party services
// @host http://localhost:9090
// @BasePath /

func main() {
	var taskHandler *repository.TaskHandler
	c, err := service.InitConfig()

	if err != nil {
		log.Fatalf("Could not read config.yaml %v: ", err)
	}

	s := service.NewService(c, taskHandler)

	err = s.InitTaskHandler()
	if err != nil {
		log.Fatalf("Could not init TaskHandler %v: ", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/task", s.CreateTask)
	mux.HandleFunc("/task/", s.GetTaskStatus)
	mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	server := &http.Server{
		Addr:    ":" + c.Server.Port,
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
