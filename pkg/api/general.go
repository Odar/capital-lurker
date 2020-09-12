package api

import "time"

type DateRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Filter struct {
	ID             *uint64    `json:"id"`
	Name           *string    `json:"name"`
	OnMainPage     *bool      `json:"on_main_page"`
	InFilter       *bool      `json:"in_filter"`
	AddedAtRange   *DateRange `json:"added_at_range"`
	UpdatedAtRange *DateRange `json:"updated_at_range"`
	Position       *uint64    `json:"position"`
	Img            *string    `json:"img"`
}
