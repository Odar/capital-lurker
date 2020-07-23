package api

import "time"
import "github.com/Odar/capital-lurker/pkg/app/models"

type DateRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Filter struct {
	Id             uint64     `json:"id"`
	Name           string     `json:"name"`
	OnMainPage     *bool      `json:"on_main_page"`
	InFilter       *bool      `json:"in_filter"`
	AddedAtRange   *DateRange `json:"added_at_range"`
	UpdatedAtRange *DateRange `json:"updated_at_range"`
	Position       uint64     `json:"position"`
	Img            string     `json:"img"`
}

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

type DeleteResponse struct {
	Whdb  string `json:"whdb"`
	Error string `json:"error"`
}

type PostIdRequest struct {
	Name       string `query:"name" json:"name"`
	OnMainPAge *bool  `query:"on_main_page" json:"on_main_p_age"`
	InFilter   *bool  `query:"in_filter" json:"in_filter"`
	Position   uint64 `query:"position" json:"position"`
	Img        string `query:"img" json:"img"`
}

type PostIdResponse struct {
	University models.University `json:"university"`
}
