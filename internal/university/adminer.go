package university

import (
	"encoding/json"
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
	"github.com/Odar/capital-lurker/pkg/app/repositories"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

func New(repo repositories.AdminerRepo) *adminer {
	return &adminer{
		repo: repo,
	}
}

type adminer struct {
	repo repositories.AdminerRepo
}

func (a *adminer) PostAdmin(ctx echo.Context) error {
	var request api.PostRequest
	err := ctx.Bind(&request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	model, err := a.postAdmin(request)
	res, _ := a.getNumUniversities() //add error processing

	if err != nil {
		log.Error().Err(err).Msgf("can not reverse with request %+v", request)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(api.PostResponse{
		Universities: model,
		Count:        res,
	})
}

func (a *adminer) postAdmin(request api.PostRequest) ([]models.University, error) {
	model, err := a.repo.GetUniversities(request.Filter)
	//fill
	return model, err
}

func (a *adminer) getNumUniversities() (uint64, error) {
	res, err := a.repo.GetNumUniversities()
	return res, err
}

//3 methods
