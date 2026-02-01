package http

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/Calmantara/lis-backend/internal/core/configurations"
	"github.com/Calmantara/lis-backend/internal/core/ports"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

func TestGracefulShutdown(t *testing.T) {
	t.Run("Graceful shutdown with no errors", func(t *testing.T) {
		configurations.Load()
		e := echo.New()
		configurations.Config.Application.Graceful = 1 // Set a short timeout for testing

		go func() {
			time.Sleep(500 * time.Millisecond) // Simulate some delay
			e.Shutdown(context.Background())   // Trigger shutdown
		}()

		GracefulShutdown(e)
		assert.True(t, true) // If no panic, the test passes
	})

	t.Run("Graceful shutdown with timeout", func(t *testing.T) {
		configurations.Load()
		e := echo.New()
		configurations.Config.Application.Graceful = 1 // Set a short timeout for testing

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
		defer cancel()

		go func() {
			time.Sleep(2 * time.Second) // Simulate a delay longer than the timeout
			e.Shutdown(ctx)             // Trigger shutdown
		}()

		GracefulShutdown(e)
		assert.True(t, true) // If no panic, the test passes
	})
}

func TestRunEcho(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_PORT", "14001")

	cmd := &cobra.Command{}
	args := []string{}

	// Simulate OS signal for shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.WithValue(t.Context(), SignalContext{}, stopChan)
	cmd.SetContext(ctx)

	go func() {
		time.Sleep(time.Second)    // Simulate some delay
		stopChan <- syscall.SIGINT // Send interrupt signal
	}()

	RunEcho(cmd, args)

	assert.True(t, true) // If no panic, the test passes
}

func TestRunEcho_WithDiggerContext(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_PORT", "14002")

	cmd := &cobra.Command{}
	args := []string{}

	// create invalid injection
	digger := dig.New()
	handler := func(str string) ports.DeviceMessageHdl {
		return nil
	}
	digger.Provide(handler)

	ctx := context.WithValue(t.Context(), DiggerContext{}, digger)
	cmd.SetContext(ctx)

	assert.Panics(t, func() {
		RunEcho(cmd, args)
	})
}

func TestRunEcho_InvalidPort(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_PORT", "1000000")

	cmd := &cobra.Command{}
	args := []string{}

	assert.Panics(t, func() {
		RunEcho(cmd, args)
	}, "Expected RunEcho to panic due to invalid port")
}
