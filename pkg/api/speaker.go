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

// API for getting speakers for administration
type DateRangeForFilter struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type SpeakerForAdminFilter struct {
	ID             uint64             `json:"id"`
	Name           string             `json:"name"`
	OnMainPage     *bool              `json:"on_main_page"`
	InFilter       *bool              `json:"in_filter"`
	AddedAtRange   DateRangeForFilter `json:"added_at_range"`
	UpdatedAtRange DateRangeForFilter `json:"updated_at_range"`
	Position       uint64             `json:"position"`
	Img            string             `json:"img"`
}

type GetSpeakersForAdminRequest struct {
	Limit  int64                 `json:"limit"`
	Page   int64                 `json:"page"`
	SortBy string                `json:"sortBy"`
	Filter SpeakerForAdminFilter `json:"filter"`
}

type GetSpeakersForAdminResponse struct {
	Speakers []models.Speaker `json:"speakers"`
	Count    uint64           `json:"count"`
}
