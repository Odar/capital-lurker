package repositories

import (
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
)

type ThemeAdminerRepo interface {
	GetThemesForAdmin(limit int64, page int64, sortBy string, filter *api.Filter) ([]models.Theme, error)
	CountThemesForAdmin(filter *api.Filter) (uint64, error)
	DeleteTheme(ID uint64) (int64, error)
	UpdateThemeForAdmin(ID uint64, request *api.UpdateThemeForAdminRequest) (*models.Theme, error)
	AddThemeForAdmin(request *api.AddThemeForAdminRequest) (*models.Theme, error)
}
