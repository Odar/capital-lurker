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
	UpdateSpeakerForAdmin(ctx echo.Context) error
	AddSpeakerForAdmin(ctx echo.Context) error
}

type ThemeAdminer interface {
	GetThemesForAdmin(ctx echo.Context) error
	DeleteThemeForAdmin(ctx echo.Context) error
	UpdateThemeForAdmin(ctx echo.Context) error
	AddThemeForAdmin(ctx echo.Context) error
}

type CourseAdminer interface {
	GetCoursesForAdmin(ctx echo.Context) error
	DeleteCourseForAdmin(ctx echo.Context) error
	UpdateCourseForAdmin(ctx echo.Context) error
	AddCourseForAdmin(ctx echo.Context) error
}
