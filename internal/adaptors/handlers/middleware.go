package handlers

import (
	"github.com/Calmantara/lis-backend/internal/core/models"
	"github.com/Calmantara/lis-backend/internal/core/ports"
	"github.com/Calmantara/lis-backend/internal/helpers/errors"
	"github.com/Calmantara/lis-backend/internal/helpers/wrappers"
	"github.com/Calmantara/lis-backend/internal/utils"
	"github.com/labstack/echo/v4"
)

type MiddlewareHdlImpl struct {
}

func NewMiddlewareHandler() ports.Middleware {
	return &MiddlewareHdlImpl{}
}

func (m *MiddlewareHdlImpl) BasicApplicationKey(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// fetch token from header
		authHeader := &models.ApplicationKeyHeader{}
		binder := &echo.DefaultBinder{} // Get an instance of the default binder
		binder.BindHeaders(c, authHeader)

		// validate header fields
		if err := authHeader.Validate(); err != nil {
			err = errors.Wrap(err, errors.ERROR_UNAUTHORIZED.Error())

			return wrappers.ConstructResponseFailure(c, err)
		}

		c.Set(string(utils.EchoClientID), authHeader.ClientID)
		c.Set(string(utils.EchoApplicationID), authHeader.ApplicationID)
		c.Set(string(utils.EchoRequestID), authHeader.RequestID)

		return next(c)
	}
}

func (m *MiddlewareHdlImpl) Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		next(c)

		// TODO: implement logger middleware here

		return nil
	}
}
