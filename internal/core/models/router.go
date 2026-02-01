package models

import "github.com/labstack/echo/v4"

type Router struct {
	External *echo.Group
	Internal *echo.Group
}
