package app

import echo "github.com/labstack/echo/v4"

type Receiver interface {
	Reverse(ctx echo.Context) error
}

type UniversityAdminer interface {
	GetUniversitiesList(ctx echo.Context) error
	AddUniversity(ctx echo.Context) error
	DeleteUniversity(ctx echo.Context) error
	UpdateUniversity(ctx echo.Context) error
}

type Speaker interface {
	GetSpeakersOnMain(ctx echo.Context) error
	GetSpeakersForAdmin(ctx echo.Context) error
	DeleteSpeakerForAdmin(ctx echo.Context) error
}

type Authenticator interface {
	Login(ctx echo.Context) error
	SignUp(ctx echo.Context) error
	SignUpViaVk(ctx echo.Context) error
	TestPage(ctx echo.Context) error
}
