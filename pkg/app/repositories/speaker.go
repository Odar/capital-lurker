package repositories

import (
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
)

type SpeakerRepo interface {
	GetSpeakersOnMain(limit int64) ([]api.SpeakerOnMain, error)
	GetSpeakersForAdmin(limit int64, page int64, sortBy string, filter *api.SpeakerForAdminFilter) (
		[]api.SpeakerForAdmin, error)
	CountSpeakersForAdmin(filter *api.SpeakerForAdminFilter) (uint64, error)
	DeleteSpeakerForAdminFromDB(ID uint64) (string, error)
	UpdateSpeakerForAdminInDB(ID uint64, request *api.UpdateSpeakerForAdminRequest) (*models.Speaker, error)
	AddSpeakerForAdminInDB(request *api.AddSpeakerForAdminRequest) (*models.Speaker, error)
}
