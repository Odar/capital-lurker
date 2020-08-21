package api

import (
	"time"

	"github.com/Odar/capital-lurker/pkg/app/models"
)

// API for getting courses for administration
type GetCoursesForAdminRequest struct {
	Limit  int64  `json:"limit"`
	Page   int64  `json:"page"`
	SortBy string `json:"sortBy"`
	Filter Filter `json:"filter"`
}

type UnparsedCourseForAdmin struct {
	ID              uint64     `db:"course_id"`
	Name            string     `db:"course_name"`
	Description     string     `db:"course_description"`
	Position        uint64     `db:"course_position"`
	AddedAt         time.Time  `db:"course_added_at"`
	UpdatedAt       time.Time  `db:"course_updated_at"`
	ThemeID         *uint64    `db:"theme_id"`
	ThemeName       *string    `db:"theme_name"`
	ThemeSlug       *string    `db:"theme_slug"`
	ThemeOnMainPage *bool      `db:"theme_on_main_page"`
	ThemeAddedAt    *time.Time `db:"theme_added_at"`
	ThemeUpdatedAt  *time.Time `db:"theme_updated_at"`
	ThemePosition   *uint64    `db:"theme_position"`
	ThemeImg        *string    `db:"theme_img"`
}

type CourseForAdmin struct {
	ID          uint64        `json:"id"`
	Name        string        `json:"name"`
	Theme       *models.Theme `json:"theme"`
	Description string        `json:"description"`
	Position    uint64        `json:"position"`
	AddedAt     time.Time     `json:"added_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type GetCoursesForAdminResponse struct {
	Courses []CourseForAdmin `json:"courses"`
	Count   uint64           `json:"count"`
}

// API for deleting course for administration
type DeleteCourseForAdminRequest struct {
	ID uint64
}

type DeleteCourseForAdminResponse struct {
	WHBD  string `json:"whbd"`
	Error string `json:"error"`
}

// API for updating course for administration
type UpdateCourseForAdminRequest struct {
	Name        *string `json:"name"`
	ThemeID     *uint64 `json:"theme_id"`
	Description *string `json:"description"`
	Position    *uint64 `json:"position"`
}

// API for adding course for administration
type AddCourseForAdminRequest struct {
	Name        string `json:"name"`
	ThemeID     uint64 `json:"theme_id"`
	Description string `json:"description"`
	Position    uint64 `json:"position"`
}
