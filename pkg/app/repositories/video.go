package repositories

import "github.com/Odar/capital-lurker/pkg/api"

type VideoAdminerRepo interface {
	GetVideosForAdmin(limit int64, page int64, sortBy string, filter *api.Filter) ([]api.VideoForAdmin, error)
	CountVideosForAdmin(filter *api.Filter) (uint64, error)
}
