package ports

import (
	"context"

	"github.com/labstack/echo/v4"
)

type (
	HealthHdl interface {
		HealthCheck(ctx echo.Context) error
	}

	HealthSvc interface {
		CheckMaster(ctx context.Context) (err error)
		CheckSlave(ctx context.Context) (err error)
	}

	HealthRepo interface {
		CheckMaster(ctx context.Context) (err error)
		CheckSlave(ctx context.Context) (err error)
	}
)
