package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Calmantara/lis-backend/internal/adaptors/handlers"
	"github.com/Calmantara/lis-backend/internal/adaptors/storage/mysql"
	"github.com/Calmantara/lis-backend/internal/core/configurations"
	"github.com/Calmantara/lis-backend/internal/core/models"
	"github.com/Calmantara/lis-backend/internal/core/services"
	"github.com/Calmantara/lis-backend/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

type (
	SignalContext struct{}
	DiggerContext struct{}
)

func RunEcho(cmd *cobra.Command, args []string) {
	configurations.Load()
	utils.NewZap(
		utils.WithAppName(configurations.Config.Application.Name),
		utils.WithEnvironment(configurations.Config.Application.Environment),
	)

	// os channel
	stopChan := getSignalChan(cmd)
	// dig dependency injection
	digger := getDigger(cmd)
	// echo server
	e := echo.New()
	routerV1 := RouterV1(e)
	digger.Provide(func() models.Router {
		return routerV1
	}, dig.Name("routerV1"))

	// digger config
	digger.Provide(func() configurations.LisPlatform {
		return configurations.Config.LisPlatform
	})

	// dependency injection
	mysql.NewInjector(digger)
	services.NewInjector(digger)
	handlers.NewInjector(digger)
	// invoke
	err := digger.Invoke(handlers.Invoke)
	if err != nil {
		panic(err)
	}

	// start server
	go func() {
		// Wait for interrupt signal to gracefully shut down the server
		<-stopChan
		err = GracefulShutdown(e)
		utils.Log.Infow("gracefully shutting down the server", map[string]any{"error": err})
	}()

	port := fmt.Sprintf(":%d", configurations.Config.Application.Port)
	utils.Log.Infow("starting HTTP server",
		map[string]any{
			"port":          port,
			"disable http2": e.DisableHTTP2,
		},
	)
	if err := e.Start(port); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func GracefulShutdown(e *echo.Echo) error {
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(configurations.Config.Application.Graceful*time.Second),
	)
	defer cancel()

	return e.Shutdown(ctx)
}

func RouterV1(e *echo.Echo) models.Router {
	return models.Router{
		External: e.Group("/api/v1"),
		Internal: e.Group("/api/internal/v1"),
	}
}

func getDigger(cmd *cobra.Command) *dig.Container {
	ctx := cmd.Context()
	if ctx != nil {
		// check digger from context
		ctxDigger := ctx.Value(DiggerContext{})
		if dg, ok := ctxDigger.(*dig.Container); ok {
			return dg
		}
	}

	return dig.New()
}

func getSignalChan(cmd *cobra.Command) chan os.Signal {
	ctx := cmd.Context()
	if ctx != nil {
		// check signal channel from context
		ctxValue := ctx.Value(SignalContext{})
		if ch, ok := ctxValue.(chan os.Signal); ok {
			return ch
		}
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	return stopChan
}
