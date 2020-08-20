package api

import (
	"time"
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

type UnparsedSpeakerOnMain struct {
	ID             uint64  `db:"speaker_id"`
	Name           string  `db:"speaker_name"`
	Position       uint64  `db:"speaker_position"`
	Img            string  `db:"speaker_img"`
	UniversityID   *uint64 `db:"university_id"`
	UniversityName *string `db:"university_name"`
	UniversityImg  *string `db:"university_img"`
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

type UnparsedSpeakerForAdmin struct {
	ID                   uint64     `db:"speaker_id"`
	Name                 string     `db:"speaker_name"`
	OnMainPage           bool       `db:"speaker_on_main_page"`
	InFilter             bool       `db:"speaker_in_filter"`
	AddedAt              time.Time  `db:"speaker_added_at"`
	UpdatedAt            time.Time  `db:"speaker_updated_at"`
	Position             uint64     `db:"speaker_position"`
	Img                  string     `db:"speaker_img"`
	UniversityID         *uint64    `db:"university_id"`
	UniversityName       *string    `db:"university_name"`
	UniversityOnMainPage *bool      `db:"university_on_main_page"`
	UniversityInFilter   *bool      `db:"university_in_filter"`
	UniversityAddedAt    *time.Time `db:"university_added_at"`
	UniversityUpdatedAt  *time.Time `db:"university_updated_at"`
	UniversityPosition   *uint64    `db:"university_position"`
	UniversityImg        *string    `db:"university_img"`
}

type UniversityForAdmin struct {
	ID         uint64    `json:"id"`
	Name       string    `json:"name"`
	OnMainPage bool      `json:"on_main_page"`
	InFilter   bool      `json:"in_filter"`
	AddedAt    time.Time `json:"added_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Position   uint64    `json:"position"`
	Img        string    `json:"img"`
}

type SpeakerForAdmin struct {
	ID         uint64              `json:"id"`
	Name       string              `json:"name"`
	OnMainPage bool                `json:"on_main_page"`
	InFilter   bool                `json:"in_filter"`
	AddedAt    time.Time           `json:"added_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
	Position   uint64              `json:"position"`
	Img        string              `json:"img"`
	University *UniversityForAdmin `json:"university"`
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
