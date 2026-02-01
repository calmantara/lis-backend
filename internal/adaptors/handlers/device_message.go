package handlers

import (
	"github.com/Calmantara/lis-backend/internal/core/models"
	"github.com/Calmantara/lis-backend/internal/core/ports"
	"github.com/Calmantara/lis-backend/internal/helpers/errors"
	"github.com/Calmantara/lis-backend/internal/helpers/wrappers"
	"github.com/Calmantara/lis-backend/internal/utils"
	"github.com/labstack/echo/v4"
)

type DeviceMessageHdlImpl struct {
	deviceMessageService ports.DeviceMessageService
	middleware           ports.Middleware
}

func NewDeviceMessageHandler(
	deviceMessageService ports.DeviceMessageService,
	middleware ports.Middleware,
) ports.DeviceMessageHdl {
	return &DeviceMessageHdlImpl{
		deviceMessageService: deviceMessageService,
		middleware:           middleware,
	}
}

func (a *DeviceMessageHdlImpl) MountV1(external, internal *echo.Group) {
	deviceMessageGroup := external.Group("/device-messages")
	{
		deviceMessageGroup.POST("", a.Create, a.middleware.BasicApplicationKey)
	}
}

func (a *DeviceMessageHdlImpl) Create(ctx echo.Context) error {
	c := utils.GetEchoContext(ctx)
	// bind payload
	params := &models.DeviceMessageParam{}
	if err := ctx.Bind(params); err != nil {
		err = errors.Wrap(err, errors.ERROR_BAD_REQUEST.Error())

		return wrappers.ConstructResponseFailure(ctx, err)
	}
	// validate payload
	if err := params.Validate(); err != nil {
		err = errors.Wrap(err, errors.ERROR_BAD_REQUEST.Error())

		return wrappers.ConstructResponseFailure(ctx, err)
	}

	err := a.deviceMessageService.Process(c, params.ToInput())
	if err != nil {
		return wrappers.ConstructResponseFailure(ctx, err)
	}
	// 	construct response
	return wrappers.ConstructResponseSuccess(
		ctx,
		wrappers.SuccessOK("ok", nil),
	)
}
