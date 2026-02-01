package handlers

import (
	"github.com/Calmantara/lis-backend/internal/core/models"
	"github.com/Calmantara/lis-backend/internal/core/ports"
	"go.uber.org/dig"
)

type InjectorInput struct {
	dig.In
	RouterV1         models.Router `name:"routerV1"`
	DeviceMessageHdl ports.DeviceMessageHdl
}

func NewInjector(digger *dig.Container) {
	digger.Provide(NewMiddlewareHandler)
	digger.Provide(NewDeviceMessageHandler)
}

func Invoke(input InjectorInput) {
	input.DeviceMessageHdl.MountV1(input.RouterV1.External, input.RouterV1.Internal)
}
