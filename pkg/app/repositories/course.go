package repositories

import (
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
)

type CourseAdminerRepo interface {
	GetCoursesForAdmin(limit int64, page int64, sortBy string, filter *api.Filter) ([]api.CourseForAdmin, error)
	CountCoursesForAdmin(filter *api.Filter) (uint64, error)
	DeleteCourse(ID uint64) (int64, error)
	UpdateCourseForAdmin(ID uint64, request *api.UpdateCourseForAdminRequest) (*models.Course, error)
	AddCourseForAdmin(request *api.AddCourseForAdminRequest) (*models.Course, error)
}
