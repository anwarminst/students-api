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

	"github.com/anwarminst/students-api/internal/config"
	"github.com/anwarminst/students-api/internal/http/student"
)

func main() {
	//load config
	cfg := config.MustLoad()
	//database setup
	//setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/student", student.New())

	//server setup
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("server started", slog.String(" ", cfg.HTTPserver.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("server closed")
		}
	}()

	<-done

	slog.Info("server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown the server:", slog.String("error", err.Error()))
	}
	slog.Info("server shutdown successfully")
}
