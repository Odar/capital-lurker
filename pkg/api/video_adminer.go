package api

import (
	"time"
)

// API for getting videos for administration
type GetVideosForAdminRequest struct {
	Limit  int64  `json:"limit"`
	Page   int64  `json:"page"`
	SortBy string `json:"sortBy"`
	Filter Filter `json:"filter"`
}

type VideoForAdmin struct {
	ID               uint64          `json:"id"`
	Name             string          `json:"name"`
	Img              string          `json:"img"`
	Video            string          `json:"video"`
	YouTubeVideo     string          `json:"youtube_video"`
	Course           *CourseForAdmin `json:"course"`
	PositionInCourse int             `json:"position_in_course"`
	AddedAt          time.Time       `json:"added_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	UploadedAt       time.Time       `json:"uploaded_at"`
	YouTubedAt       time.Time       `json:"youtubed_at"`
}

type GetVideosForAdminResponse struct {
	Videos []VideoForAdmin `json:"courses"`
	Count  uint64          `json:"count"`
}

// API for deleting video for administration
type DeleteVideoForAdminRequest struct {
	ID uint64
}

type DeleteVideoForAdminResponse struct {
	WHBD  string `json:"whbd"`
	Error string `json:"error"`
}

// API for updating video for administration
type UpdateVideoForAdminRequest struct {
	Name             *string    `json:"name"`
	Img              *string    `json:"img"`
	Video            *string    `json:"video"`
	YouTubeVideo     *string    `json:"youtube_video"`
	CourseID         *uint64    `json:"course_id"`
	PositionInCourse *int       `json:"position_in_course"`
	UploadedAt       *time.Time `json:"uploaded_at"`
	YouTubedAt       *time.Time `json:"youtubed_at"`
}

// API for adding course for administration
type AddVideoForAdminRequest struct {
	Name             string    `json:"name"`
	Img              string    `json:"img"`
	Video            string    `json:"video"`
	YouTubeVideo     string    `json:"youtube_video"`
	CourseID         uint64    `json:"course_id"`
	PositionInCourse int       `json:"position_in_course"`
	UploadedAt       time.Time `json:"uploaded_at"`
	YouTubedAt       time.Time `json:"youtubed_at"`
}
