package repositories

import (
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
)

type SpeakerRepo interface {
	GetSpeakersOnMain(limit int64) ([]api.SpeakerOnMain, error)
	GetSpeakersForAdminFromDB(limit int64, page int64, sortBy string, filter *api.SpeakerForAdminFilter) (
		[]models.Speaker, uint64, error)
}
