package api

import (
	"time"

	"github.com/Odar/capital-lurker/pkg/app/models"
)

// API for getting speakers for main page
type GetSpeakersOnMainRequest struct {
	Limit int64 `query:"limit"`
}

type SpeakerOnMain struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Position uint64 `json:"position"`
	Img      string `json:"img"`
}

type GetSpeakersOnMainResponse struct {
	Speakers []SpeakerOnMain `json:"speakers"`
}

type GetSpeakersForAdminRequest struct {
	Limit  int64  `json:"limit"`
	Page   int64  `json:"page"`
	SortBy string `json:"sortBy"`
	Filter Filter `json:"filter"`
}

type GetSpeakersForAdminResponse struct {
	Speakers []models.Speaker `json:"speakers"`
	Count    uint64           `json:"count"`
}

// API for deleting speaker for administration
type DeleteSpeakerForAdminRequest struct {
	ID uint64
}

type DeleteSpeakerForAdminResponse struct {
	WHBD  string `json:"whbd"`
	Error string `json:"error"`
}

// API for updating speaker for administration
type UpdateSpeakerForAdminRequest struct {
	Name       *string `json:"name"`
	OnMainPage *bool   `json:"on_main_page"`
	InFilter   *bool   `json:"in_filter"`
	Position   *uint64 `json:"position"`
	Img        *string `json:"img"`
}

type UpdateSpeakerForAdminResponse struct {
	ID         uint64    `db:"id"`
	Name       string    `db:"name"`
	OnMainPage bool      `db:"on_main_page"`
	InFilter   bool      `db:"in_filter"`
	AddedAt    time.Time `db:"added_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Position   uint64    `db:"position"`
	Img        string    `db:"img"`
}

// API for adding speaker for administration
type AddSpeakerForAdminRequest struct {
	Name       string `json:"name"`
	OnMainPage bool   `json:"on_main_page"`
	InFilter   bool   `json:"in_filter"`
	Position   uint64 `json:"position"`
	Img        string `json:"img"`
}

type AddSpeakerForAdminResponse struct {
	ID         uint64    `db:"id"`
	Name       string    `db:"name"`
	OnMainPage bool      `db:"on_main_page"`
	InFilter   bool      `db:"in_filter"`
	AddedAt    time.Time `db:"added_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Position   uint64    `db:"position"`
	Img        string    `db:"img"`
}
