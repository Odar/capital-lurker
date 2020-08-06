package repositories

import (
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
)

type SpeakerRepo interface {
	GetSpeakersOnMain(limit int64) ([]api.SpeakerOnMain, error)
	GetSpeakersForAdmin(limit int64, page int64, sortBy string, filter *api.Filter) ([]models.Speaker,
		error)
	CountSpeakersForAdmin(filter *api.Filter) (uint64, error)
	DeleteSpeakerForAdmin(ID uint64) (string, error)
}
