package models

import "time"

type Speaker struct {
	ID           uint64    `db:"id"`
	Name         string    `db:"name"`
	OnMainPage   bool      `db:"on_main_page"`
	InFilter     bool      `db:"in_filter"`
	AddedAt      time.Time `db:"added_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	Position     uint64    `db:"position"`
	Img          string    `db:"img"`
	UniversityID uint64    `db:"university_id"`
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
