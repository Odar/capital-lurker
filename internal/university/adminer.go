package university

import (
	"encoding/json"
	"net/http"

	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
	"github.com/Odar/capital-lurker/pkg/app/repositories"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func New(repo repositories.AdminerRepo) *adminer {
	return &adminer{
		repo: repo,
	}
}

type adminer struct {
	repo repositories.AdminerRepo
}

func (a *adminer) GetUniversitiesList(ctx echo.Context) error {
	var request api.PostRequest
	err := ctx.Bind(&request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	model, count, err := a.getUniversitiesList(request)
	if err != nil {
		log.Error().Err(err).Msgf("can not get universities list with request %+v", request)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	ctx.Response().WriteHeader(http.StatusOK)
	if model == nil {
		return json.NewEncoder(ctx.Response()).Encode(api.PostResponse{
			Universities: []models.University{},
			Count:        0,
		})
	}
	return json.NewEncoder(ctx.Response()).Encode(api.PostResponse{
		Universities: model,
		Count:        count,
	})
}

func (a *adminer) getUniversitiesList(request api.PostRequest) ([]models.University, uint64, error) {
	if request.Limit <= 0 {
		request.Limit = 10
	}
	if request.Page <= 0 {
		request.Page = 1
	}
	universities, err := a.repo.GetUniversitiesList(request.Filter, request.SortBy, request.Limit, request.Page)
	if err != nil {
		return nil, 0, errors.Wrap(err, "can not get from universities list")
	}

	count, err := a.repo.GetUniversitiesCount(request.Filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "can not get universities count")
	}
	if universities == nil && count != 0 {
		return nil, count, errors.Errorf("no such page")
	}
	return universities, count, err
}
