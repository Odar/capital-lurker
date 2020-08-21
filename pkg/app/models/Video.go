package models

import "time"

type Video struct {
	ID               uint64    `db:"id"`
	Name             string    `db:"name"`
	Img              string    `db:"img"`
	Video            string    `db:"video"`
	YouTubeVideo     string    `db:"youtube_video"`
	CourseID         uint64    `db:"course_id"`
	PositionInCourse int       `db:"position_in_course"`
	AddedAt          time.Time `db:"added_at"`
	UpdatedAt        time.Time `db:"updated_at"`
	UploadedAt       time.Time `db:"uploaded_at"`
	YouTubedAt       time.Time `db:"youtubed_at"`
}
