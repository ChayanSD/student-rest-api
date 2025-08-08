package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChayanSD/student-rest-api/internal/config"
	"github.com/ChayanSD/student-rest-api/internal/http/handlers/student"
)

func main() {
	//load config

	cfg := config.MustLoad()

	//set up route

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started at", slog.String("address", server.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGALRM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown the server :", slog.String("error", err.Error()))
	}
	slog.Info("server shutdown successfully")
}
