package ports

import (
	"context"

	"github.com/Calmantara/lis-backend/internal/core/models"
	"github.com/labstack/echo/v4"
)

type (
	DeviceMessageCommand interface {
		Create(ctx context.Context, user *models.DeviceMessage) error
	}

	DeviceMessageService interface {
		Process(ctx context.Context, inputs *models.DeviceMessageInput) (err error)
	}

	DeviceMessageHdl interface {
		RouterV1
		Create(ctx echo.Context) error
	}
)
