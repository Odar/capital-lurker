package api

import (
	"github.com/Odar/capital-lurker/pkg/app/models"
)

// API for getting themes for administration
type GetThemesForAdminRequest struct {
	Limit  int64  `json:"limit"`
	Page   int64  `json:"page"`
	SortBy string `json:"sortBy"`
	Filter Filter `json:"filter"`
}

type GetThemesForAdminResponse struct {
	Themes []models.Theme `json:"themes"`
	Count  uint64         `json:"count"`
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
