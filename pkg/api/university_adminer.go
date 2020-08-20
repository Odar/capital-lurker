package api

import (
	"github.com/Odar/capital-lurker/pkg/app/models"
)

type PostRequest struct {
	Limit  int     `query:"limit" json:"limit"`
	Page   int     `query:"page" json:"page"`
	SortBy string  `query:"sortBy" json:"sortBy"`
	Filter *Filter `query:"filter" json:"filter"`
}

type PostResponse struct {
	Universities []models.University `json:"universities"`
	Count        uint64              `json:"count"`
}

type PutRequest struct {
	Name       string `query:"name" json:"name"`
	OnMainPage bool   `query:"on_main_page" json:"on_main_page"`
	InFilter   bool   `query:"in_filter" json:"in_filter"`
	Position   uint64 `query:"position" json:"position"`
	Img        string `query:"img" json:"img"`
}

type PutResponse struct {
	University models.University `json:"university"`
}

type DeleteUniversityResponse struct {
	WHBD  string `json:"whbd"`
	Error string `json:"error"`
}

type UpdateUniversityRequest struct {
	Name       *string `query:"name" json:"name"`
	OnMainPage *bool   `query:"on_main_page" json:"on_main_page"`
	InFilter   *bool   `query:"in_filter" json:"in_filter"`
	Position   *uint64 `query:"position" json:"position"`
	Img        *string `query:"img" json:"img"`
}

type UpdateUniversityResponse struct {
	University models.University `json:"university"`
}
