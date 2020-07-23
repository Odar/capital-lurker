package app

import "github.com/labstack/echo/v4"

type Receiver interface {
	Reverse(ctx echo.Context) error
}

//heze
type UniversityAdminer interface {
	PostAdmin(ctx echo.Context) error
	PutAdmin(ctx echo.Context) error
	DeleteAdmin(ctx echo.Context) error
	PostIdAdmin(ctx echo.Context) error
}
