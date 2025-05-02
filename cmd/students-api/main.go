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
	"github.com/anwarminst/students-api/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()
	//database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage initialized", slog.String("env", cfg.Env))
	//setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/student", student.New(storage))
	router.HandleFunc("GET /api/student/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetStudents(storage))
	router.HandleFunc("PUT /api/student/{id}", student.UpdateStudent(storage))

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
