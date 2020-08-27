package models

import "time"

type Course struct {
	ID          uint64    `db:"id"`
	Name        string    `db:"name"`
	ThemeID     uint64    `db:"theme_id"`
	Description string    `db:"description"`
	Position    uint64    `db:"position"`
	AddedAt     time.Time `db:"added_at"`
	UpdatedAt   time.Time `db:"updated_at"`
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
