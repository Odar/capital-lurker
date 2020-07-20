package api

import "time"
import "github.com/Odar/capital-lurker/pkg/app/models"

type DateRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Filter struct {
	Id             uint64    `json:"id"`
	Name           string    `json:"name"`
	OnMainPage     bool      `json:"on_main_page"`
	InFilter       bool      `json:"in_filter"`
	AddedAtRange   DateRange `json:"added_at_range"`
	UpdatedAtRange DateRange `json:"updated_at_range"`
	Position       uint64    `json:"position"`
	Img            string    `json:"img"`
}

type PostRequest struct {
	Limit  int    `query:"limit"`
	Page   int    `query:"page"`
	SortBy string `query:"sortBy"`
	Filter Filter `query:"filter"`
}

type PostResponse struct {
	Universities []models.University `json:"universities"` //должен ли я добавить теги для University?
	Count        uint64              `json:"count"`
}
