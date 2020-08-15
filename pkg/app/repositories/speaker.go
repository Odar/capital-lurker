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
	DeleteSpeaker(ID uint64) (int64, error)
	UpdateSpeakerForAdmin(ID uint64, request api.UpdateSpeakerForAdminRequest) (*models.Speaker, error)
}
