package theme

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
	"github.com/Odar/capital-lurker/pkg/app/repositories"
	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func New(repo repositories.ThemeAdminerRepo) *adminer {
	return &adminer{
		repo: repo,
	}
}

type adminer struct {
	repo repositories.ThemeAdminerRepo
}

func (a *adminer) GetThemesForAdmin(ctx echo.Context) error {
	var request api.GetThemesForAdminRequest
	err := ctx.Bind(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not retrieve data from JSON:%+v", request)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	themesForAdmin, count, err := a.getThemesForAdmin(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not get theme for admin with request %+v", request)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(api.GetThemesForAdminResponse{
		Themes: themesForAdmin,
		Count:  count,
	})
}

func (a *adminer) getThemesForAdmin(request *api.GetThemesForAdminRequest) ([]models.Theme, uint64, error) {
	if request.Limit <= 0 {
		request.Limit = 10
	}
	if request.Page <= 0 {
		request.Page = 1
	}

	themes, err := a.repo.GetThemesForAdmin(request.Limit, request.Page, request.SortBy, &request.Filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "can not get themes from db")
	}

	count, err := a.repo.CountThemesForAdmin(&request.Filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "can not count themes from db")
	}

	if len(themes) == 0 {
		return nil, 0, nil
	}

	return themes, count, nil
}

func (a *adminer) DeleteThemeForAdmin(ctx echo.Context) error {
	var request api.DeleteThemeForAdminRequest
	var err error
	request.ID, err = strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msgf("can not retrieve id:%+v", ctx.Param("id"))
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	var WHBD string
	WHBD, err = a.deleteThemeForAdmin(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not delete theme for admin with request %+v", request)
		ctx.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(ctx.Response()).Encode(api.DeleteThemeForAdminResponse{
			WHBD:  "error",
			Error: err.Error(),
		})
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(api.DeleteThemeForAdminResponse{
		WHBD:  WHBD,
		Error: "",
	})
}

func (a *adminer) deleteThemeForAdmin(request *api.DeleteThemeForAdminRequest) (string, error) {
	count, err := a.repo.DeleteTheme(request.ID)
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

func (a *adminer) UpdateThemeForAdmin(ctx echo.Context) error {
	var request api.UpdateThemeForAdminRequest
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

	updatedTheme, err := a.updateThemeForAdmin(requestID, &request)
	if err != nil {
		log.Error().Err(err).Msgf("can not update theme for admin with request %+v", request)
		ctx.Response().WriteHeader(http.StatusInternalServerError)
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(updatedTheme)
}

func (a *adminer) updateThemeForAdmin(requestID uint64, request *api.UpdateThemeForAdminRequest) (
	*models.Theme, error) {
	updatedTheme, err := a.repo.UpdateThemeForAdmin(requestID, request)
	if err != nil {
		return updatedTheme, errors.Wrap(err, "can not update in db")
	}

	return updatedTheme, nil
}

func (a *adminer) AddThemeForAdmin(ctx echo.Context) error {
	var request api.AddThemeForAdminRequest
	err := ctx.Bind(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not retrieve data from JSON:%+v", request)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	addedTheme, err := a.addThemeForAdmin(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not add theme for admin with request %+v", request)
		ctx.Response().WriteHeader(http.StatusInternalServerError)
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(addedTheme)
}

func (a *adminer) addThemeForAdmin(request *api.AddThemeForAdminRequest) (*models.Theme, error) {
	addedTheme, err := a.repo.AddThemeForAdmin(request)
	if err != nil {
		return addedTheme, errors.Wrap(err, "can not add in db")
	}

	return addedTheme, nil
}
