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

	model, err := a.getUniversitiesList(request)
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
		Count:        uint64(len(model)),
	})
}

func (a *adminer) getUniversitiesList(request api.PostRequest) ([]models.University, error) {
	universities, err := a.repo.GetUniversities(request.Filter, request.SortBy)
	if err != nil {
		return nil, errors.Wrap(err, "can not get from universities list")
	}
	if universities == nil {
		return nil, err
	}

	limit, page := request.Limit, request.Page
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	if page > 1+len(universities)/limit {
		return nil, errors.New("no such page")
	}

	model := make([]models.University, 0, limit)
	for i := 0; i < limit && (page-1)*limit+i < len(universities); i++ {
		model = append(model, universities[(page-1)*limit+i])
	}
	return model, err
}
