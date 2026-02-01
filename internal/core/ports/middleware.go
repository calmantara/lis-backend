package ports

import (
	"github.com/labstack/echo/v4"
)

type (
	Middleware interface {
		BasicApplicationKey(next echo.HandlerFunc) echo.HandlerFunc
		Logger(next echo.HandlerFunc) echo.HandlerFunc
	}
)
