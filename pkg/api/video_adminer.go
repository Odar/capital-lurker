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

type UnparsedVideoForAdmin struct {
	ID                uint64     `db:"video_id"`
	Name              string     `db:"video_name"`
	Img               string     `db:"video_img"`
	Video             string     `db:"video_video"`
	YouTubeVideo      string     `db:"video_youtube_video"`
	PositionInCourse  int        `db:"video_position_in_course"`
	AddedAt           time.Time  `db:"video_added_at"`
	UpdatedAt         time.Time  `db:"video_updated_at"`
	UploadedAt        time.Time  `db:"video_uploaded_at"`
	YouTubedAt        time.Time  `db:"video_youtubed_at"`
	CourseID          *uint64    `db:"course_id"`
	CourseName        *string    `db:"course_name"`
	CourseDescription *string    `db:"course_description"`
	CoursePosition    *uint64    `db:"course_position"`
	CourseAddedAt     *time.Time `db:"course_added_at"`
	CourseUpdatedAt   *time.Time `db:"course_updated_at"`
	ThemeID           *uint64    `db:"theme_id"`
	ThemeName         *string    `db:"theme_name"`
	ThemeSlug         *string    `db:"theme_slug"`
	ThemeOnMainPage   *bool      `db:"theme_on_main_page"`
	ThemeAddedAt      *time.Time `db:"theme_added_at"`
	ThemeUpdatedAt    *time.Time `db:"theme_updated_at"`
	ThemePosition     *uint64    `db:"theme_position"`
	ThemeImg          *string    `db:"theme_img"`
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
