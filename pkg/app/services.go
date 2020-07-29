package app

import "github.com/labstack/echo/v4"

type Receiver interface {
	Reverse(ctx echo.Context) error
}

type UniversityAdminer interface {
	GetUniversitiesList(ctx echo.Context) error
	PutAdmin(ctx echo.Context) error
}

type Speaker interface {
	GetSpeakersOnMain(ctx echo.Context) error
}
