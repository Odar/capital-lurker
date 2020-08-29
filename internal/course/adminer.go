package course

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Odar/capital-lurker/pkg/app/models"

	"github.com/pkg/errors"

	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/repositories"
	echo "github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func New(repo repositories.CourseAdminerRepo) *adminer {
	return &adminer{
		repo: repo,
	}
}

type adminer struct {
	repo repositories.CourseAdminerRepo
}

func (a *adminer) GetCoursesForAdmin(ctx echo.Context) error {
	var request api.GetCoursesForAdminRequest
	err := ctx.Bind(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not retrieve data from JSON:%+v", request)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	coursesForAdmin, count, err := a.getCoursesForAdmin(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not get courses for admin with request %+v", request)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(api.GetCoursesForAdminResponse{
		Courses: coursesForAdmin,
		Count:   count,
	})
}

func (a *adminer) getCoursesForAdmin(request *api.GetCoursesForAdminRequest) ([]api.CourseForAdmin, uint64, error) {
	if request.Limit <= 0 {
		request.Limit = 10
	}
	if request.Page <= 0 {
		request.Page = 1
	}

	courses, err := a.repo.GetCoursesForAdmin(request.Limit, request.Page, request.SortBy, &request.Filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "can not get courses from db")
	}

	count, err := a.repo.CountCoursesForAdmin(&request.Filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "can not count courses from db")
	}

	if len(courses) == 0 {
		return nil, 0, nil
	}

	return courses, count, nil
}

func (a *adminer) DeleteCourseForAdmin(ctx echo.Context) error {
	var request api.DeleteCourseForAdminRequest
	var err error
	request.ID, err = strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msgf("can not retrieve id:%+v", ctx.Param("id"))
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	var WHBD string
	WHBD, err = a.deleteCourseForAdmin(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not delete course for admin with request %+v", request)
		ctx.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(ctx.Response()).Encode(api.DeleteCourseForAdminResponse{
			WHBD:  "error",
			Error: err.Error(),
		})
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(api.DeleteCourseForAdminResponse{
		WHBD:  WHBD,
		Error: "",
	})
}

func (a *adminer) deleteCourseForAdmin(request *api.DeleteCourseForAdminRequest) (string, error) {
	count, err := a.repo.DeleteCourse(request.ID)
	if err != nil {
		return "error", errors.Wrap(err, "can not delete from db")
	}
	if count == 1 {
		return "deleted", nil
	}
	if count == 0 {
		return "nothing", nil
	}
	return "error", errors.New("something went wrong")
}

func (a *adminer) UpdateCourseForAdmin(ctx echo.Context) error {
	var request api.UpdateCourseForAdminRequest
	err := ctx.Bind(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not retrieve data from JSON:%+v", request)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	requestID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msgf("can not retrieve id:%+v", ctx.Param("id"))
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	updatedCourse, err := a.updateCourseForAdmin(requestID, &request)
	if err != nil {
		log.Error().Err(err).Msgf("can not update course for admin with request %+v", request)
		ctx.Response().WriteHeader(http.StatusInternalServerError)
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(updatedCourse)
}

func (a *adminer) updateCourseForAdmin(requestID uint64, request *api.UpdateCourseForAdminRequest) (
	*models.Course, error) {
	updatedCourse, err := a.repo.UpdateCourseForAdmin(requestID, request)
	if err != nil {
		return updatedCourse, errors.Wrap(err, "can not update in db")
	}

	return updatedCourse, nil
}

func (a *adminer) AddCourseForAdmin(ctx echo.Context) error {
	var request api.AddCourseForAdminRequest
	err := ctx.Bind(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not retrieve data from JSON:%+v", request)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	addedCourse, err := a.addCourseForAdmin(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not add course for admin with request %+v", request)
		ctx.Response().WriteHeader(http.StatusInternalServerError)
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(addedCourse)
}

func (a *adminer) addCourseForAdmin(request *api.AddCourseForAdminRequest) (*models.Course, error) {
	addedCourse, err := a.repo.AddCourseForAdmin(request)
	if err != nil {
		return addedCourse, errors.Wrap(err, "can not add into db")
	}

	return addedCourse, nil
}
