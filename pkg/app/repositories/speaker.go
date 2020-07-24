package repositories

import (
    "github.com/Odar/capital-lurker/pkg/api"
    "github.com/Odar/capital-lurker/pkg/app/models"
)

type SpeakerRepo interface {
    GetSpeakerOnMainFromDB(limit int64) ([]api.SpeakerOnMain, error)
    GetSpeakerForAdminFromDB(decodedParams *api.GetSpeakerForAdminRequest) ([]models.Speaker, uint64, error)
    DeleteSpeakerForAdminFromDB(ID uint64) (string, error)
}
