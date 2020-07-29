package repositories

import (
	"github.com/Odar/capital-lurker/pkg/api"
)

type SpeakerRepo interface {
	GetSpeakersOnMain(limit int64) ([]api.SpeakerOnMain, error)
}
