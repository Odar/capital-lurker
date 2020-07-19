package app

import "github.com/labstack/echo/v4"

type Receiver interface {
	Reverse(ctx echo.Context) error
}

//heze
type Adminer interface {
	PostAdmin(ctx echo.Context) error
}
