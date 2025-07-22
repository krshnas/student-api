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

	"github.com/krshnas/students-api/internal/config"
	"github.com/krshnas/students-api/internal/http/handlers/student"
	"github.com/krshnas/students-api/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// setup custom logger - do it at your own choice
	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students/", student.GetList(storage))
	router.HandleFunc("PATCH /api/students/{id}", student.UpdateById(storage))
	// delete the student
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteById(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // look for it, when these signals triggers

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done
	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}
