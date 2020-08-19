package api

import "time"

// API for getting courses for administration
type GetCoursesForAdminRequest struct {
	Limit  int64  `json:"limit"`
	Page   int64  `json:"page"`
	SortBy string `json:"sortBy"`
	Filter Filter `json:"filter"`
}

type CourseForAdmin struct {
	ID          uint64         `json:"id"`
	Name        string         `json:"name"`
	Theme       *ThemeForAdmin `json:"theme"`
	Description string         `json:"description"`
	AddedAt     time.Time      `json:"added_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
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
	Description *string `json:"description"`
}

// API for adding course for administration
type AddCourseForAdminRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
