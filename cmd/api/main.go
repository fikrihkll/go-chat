package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	chatDeliveryHttp "github.com/fikrihkll/chat-app/application/chat/delivery/http"
	repositories "github.com/fikrihkll/chat-app/application/chat/repositories"
	usecases "github.com/fikrihkll/chat-app/application/chat/usecases"

	middlewares "github.com/fikrihkll/chat-app/common/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/fikrihkll/chat-app/common"
	"github.com/fikrihkll/chat-app/config"
	"github.com/fikrihkll/chat-app/infrastructure"
	"github.com/labstack/echo/v4"
)

func count() {
	tasks := []string{"task1", "task2", "task3", "task4", "task5", "task6"}
	concurrencyLimit := 3
	sem := make(chan struct {}, concurrencyLimit)
	var wg sync.WaitGroup

	for _, task := range tasks {
		wg.Add(1)
		go func (task string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <- sem }()

			fmt.Printf("task: %s\n", task)
			time.Sleep(time.Second * 2)
		}(task)
	}
	wg.Wait()
	fmt.Println("count finish")
}

func initInfrastructure(cfg config.ApplicationConfig) (pgConn *sql.DB) {

	pgConn, err := infrastructure.NewPgConnection(cfg)
	common.LogExit(err, common.LOG_LEVEL_ERROR)
	
	return pgConn
}

func initApplication(httpServer *echo.Echo, cfg config.ApplicationConfig) {
	pgConn := initInfrastructure(cfg)

	// Repositories
	chatPersistRepo := repositories.NewChatRepositoryPostgree(pgConn)
	userPersistRepo := repositories.NewUserRepositoryPostgree(pgConn)

	// usecases
	chatUsecases := usecases.NewChatApplication(chatPersistRepo, userPersistRepo)
	authUsecases := usecases.NewUserApplication(userPersistRepo)

	
	httpApi := chatDeliveryHttp.NewChatHttpApi(chatUsecases, authUsecases)
	
	// handle http request response
	httpApi.HandleAuthRoute(httpServer)
	httpApi.HandleChatRoute(httpServer)
	httpApi.HandleRootRoute(httpServer)
}

// @title Chat API
// @version 1.0
// @description Go implemented api.
// @contact.name Fikri Haikal
// @contact.url https://github.com/fikrihkll
// @contact.email fkrihkl@gmail.com
// @license.name Apache 2.0
// @BasePath /
// @SecurityDefinitions.BearerAuth type: apiKey, in: header, name: Authorization
// @Security BearerAuth
func main() {
	cfg := config.Load()
	e := echo.New()
	e.Use(middlewares.MiddlewaresRegistry...)

	initApplication(e, cfg)

	// swagger api docs
	url := echoSwagger.URL(fmt.Sprintf("http://localhost:%s/docs/swagger.yaml", cfg.HTTPApiPort))
	e.GET("swagger/*", echoSwagger.EchoWrapHandler(url))
	e.Static("/docs", "docs")

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.HTTPApiPort)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server", err.Error())
		}
	}()

	// graceful shutdown
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 60 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}