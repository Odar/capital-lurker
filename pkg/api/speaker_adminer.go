package api

import (
	"time"

	"github.com/Odar/capital-lurker/pkg/app/models"
)

// API for getting speakers for main page
type GetSpeakersOnMainRequest struct {
	Limit int64 `query:"limit"`
}

type UniversityOnMain struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

type SpeakerOnMain struct {
	ID         uint64            `json:"id" db:"speaker_id"`
	Name       string            `json:"name" db:"speaker_name"`
	Position   uint64            `json:"position" db:"speaker_position"`
	Img        string            `json:"img" db:"speaker_img"`
	University *UniversityOnMain `json:"university"`
}

type GetSpeakersOnMainResponse struct {
	Speakers []SpeakerOnMain `json:"speakers"`
}

// API for getting speakers for administration
type GetSpeakersForAdminRequest struct {
	Limit  int64  `json:"limit"`
	Page   int64  `json:"page"`
	SortBy string `json:"sortBy"`
	Filter Filter `json:"filter"`
}

type SpeakerForAdmin struct {
	ID         uint64             `json:"id"`
	Name       string             `json:"name"`
	OnMainPage bool               `json:"on_main_page"`
	InFilter   bool               `json:"in_filter"`
	AddedAt    time.Time          `json:"added_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	Position   uint64             `json:"position"`
	Img        string             `json:"img"`
	University *models.University `json:"university"`
}

type GetSpeakersForAdminResponse struct {
	Speakers []SpeakerForAdmin `json:"speakers"`
	Count    uint64            `json:"count"`
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
	Name         *string `json:"name"`
	OnMainPage   *bool   `json:"on_main_page"`
	InFilter     *bool   `json:"in_filter"`
	Position     *uint64 `json:"position"`
	Img          *string `json:"img"`
	UniversityID *uint64 `json:"university_id"`
}

// API for adding speaker for administration
type AddSpeakerForAdminRequest struct {
	Name         string `json:"name"`
	OnMainPage   bool   `json:"on_main_page"`
	InFilter     bool   `json:"in_filter"`
	Position     uint64 `json:"position"`
	Img          string `json:"img"`
	UniversityID uint64 `json:"university_id"`
}
