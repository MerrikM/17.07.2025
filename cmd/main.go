package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "workmate_test_project/docs"
	"workmate_test_project/internal/config"
	"workmate_test_project/internal/handler"
	"workmate_test_project/internal/service"
)

// @title Auth API Service
// @version 1.0
// @description Тестовое задание на позицию Junior Go-разработчик

// @host localhost:8080
// @BasePath /api-tasks
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("ошибка загрузки конфигурации: %v", err)
	}

	srv, router := config.SetupServer(cfg.Server.Port)

	router.Get("/swagger/*", httpSwagger.WrapHandler)
	taskService := service.NewTaskService()
	taskHandler := handler.NewTaskHandler(taskService)

	router.Route(cfg.Server.BasePath, func(r chi.Router) {
		r.Post("/create-task", taskHandler.CreateTask)
		r.Get("/get", taskHandler.GetTaskStatusById)
		r.Post("/add-file-to-task", taskHandler.AddFileToTask)
	})

	runServer(ctx, srv)
}

func runServer(ctx context.Context, server *http.Server) {
	serverErrors := make(chan error, 1)
	go func() {
		log.Println("сервер запущен на " + server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil {
			log.Fatalf("ошибка работы сервера: %v", err)
		}
	case sig := <-signalChannel:
		log.Printf("получен сигнал %v остановки работы сервера ", sig)
	}

	shutDownCtx, shutDownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutDownCancel()

	if err := server.Shutdown(shutDownCtx); err != nil {
		log.Printf("ошибка при остановке сервера: %v", err)
	} else {
		log.Println("Сервер успешно остановлен")
	}
}
