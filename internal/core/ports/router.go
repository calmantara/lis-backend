package ports

import "github.com/labstack/echo/v4"

type RouterV1 interface {
	MountV1(external, internal *echo.Group)
}
