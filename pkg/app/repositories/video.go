package repositories

import (
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
)

type VideoAdminerRepo interface {
	GetVideosForAdmin(limit int64, page int64, sortBy string, filter *api.Filter) ([]api.VideoForAdmin, error)
	CountVideosForAdmin(filter *api.Filter) (uint64, error)
	DeleteVideo(ID uint64) (int64, error)
	UpdateVideoForAdmin(ID uint64, request *api.UpdateVideoForAdminRequest) (*models.Video, error)
	AddVideoForAdmin(request *api.AddVideoForAdminRequest) (*models.Video, error)
}
