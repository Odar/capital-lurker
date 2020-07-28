package repositories

import (
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
)

type SpeakerRepo interface {
	GetSpeakersOnMain(limit int64) ([]api.SpeakerOnMain, error)
	GetSpeakersForAdmin(limit int64, page int64, sortBy string, filter *api.SpeakerForAdminFilter) ([]models.Speaker,
		error)
	CountSpeakersForAdmin(page int64, sortBy string, filter *api.SpeakerForAdminFilter) (uint64, error)
}
