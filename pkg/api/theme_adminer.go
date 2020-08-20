package api

import "time"

// API for getting themes for administration
type GetThemesForAdminRequest struct {
	Limit  int64  `json:"limit"`
	Page   int64  `json:"page"`
	SortBy string `json:"sortBy"`
	Filter Filter `json:"filter"`
}

type ThemeForAdmin struct {
	ID         uint64    `db:"id"`
	Name       string    `db:"name"`
	Slug       string    `db:"slug"`
	OnMainPage bool      `db:"on_main_page"`
	AddedAt    time.Time `db:"added_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Position   uint64    `db:"position"`
	Img        string    `db:"img"`
}

type GetThemesForAdminResponse struct {
	Themes []ThemeForAdmin `json:"themes"`
	Count  uint64          `json:"count"`
}

// API for deleting theme for administration
type DeleteThemeForAdminRequest struct {
	ID uint64
}

type DeleteThemeForAdminResponse struct {
	WHBD  string `json:"whbd"`
	Error string `json:"error"`
}

// API for updating theme for administration
type UpdateThemeForAdminRequest struct {
	Name       *string `json:"name"`
	OnMainPage *bool   `json:"on_main_page"`
	Position   *uint64 `json:"position"`
	Img        *string `json:"img"`
}

// API for adding theme for administration
type AddThemeForAdminRequest struct {
	Name       string `json:"name"`
	OnMainPage bool   `json:"on_main_page"`
	Position   uint64 `json:"position"`
	Img        string `json:"img"`
}
