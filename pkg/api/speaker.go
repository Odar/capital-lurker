package api

import (
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
