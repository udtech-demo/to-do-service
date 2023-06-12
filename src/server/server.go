package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
	"todo-service/src/infrastructure/authentication"
	"todo-service/src/infrastructure/delivery/graphql"
	"todo-service/src/infrastructure/monitoring/logs"
	"todo-service/src/infrastructure/storage"
	"todo-service/src/registry"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

type App struct {
	httpServer *http.Server
	e          *echo.Echo
}

func NewApp() *App {

	// Logger
	logger := logs.NewLogger()

	//Init storages
	db := storage.InitPostgres(logger)

	jc := authentication.NewJwtConfigurator(logger, "", "")

	// Register and create controller
	useCase := registry.NewRegistry(db, jc).NewUseCase()

	// Initialize Echo instance
	e := echo.New()

	graphql.NewGraphqlRouter(e, useCase)

	// Start server
	s := &http.Server{
		Addr:           viper.GetString("http.port"),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &App{
		httpServer: s,
		e:          e,
	}
}

func (a *App) Run() error {
	// Start server
	go func() {
		if err := a.e.StartServer(a.httpServer); err != nil {
			panic(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.e.Shutdown(ctx); err != nil {
		a.e.Logger.Fatal(err)
	}

	return nil
}
