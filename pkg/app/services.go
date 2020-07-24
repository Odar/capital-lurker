package app

import "github.com/labstack/echo/v4"

type Receiver interface {
    Reverse(ctx echo.Context) error
}

type Speaker interface {
    GetSpeakerOnMain(ctx echo.Context) error
    GetSpeakerForAdmin(ctx echo.Context) error
}
