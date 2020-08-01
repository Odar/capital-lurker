package university

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	universities, count, err := a.getUniversitiesList(request)
	if err != nil {
		log.Error().Err(err).Msgf("can not get universities list with request %+v", request)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	ctx.Response().WriteHeader(http.StatusOK)
	if universities == nil {
		return json.NewEncoder(ctx.Response()).Encode(api.PostResponse{
			Universities: []models.University{},
			Count:        0,
		})
	}
	return json.NewEncoder(ctx.Response()).Encode(api.PostResponse{
		Universities: universities,
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

	count, err := a.repo.CountUniversities(request.Filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "can not get universities count")
	}
	if universities == nil && count != 0 {
		return nil, count, errors.Errorf("no such page")
	}
	return universities, count, err
}

func (a *adminer) AddUniversity(ctx echo.Context) error {
	var request api.PutRequest
	err := ctx.Bind(&request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	model, err := a.addUniversity(request)
	if err != nil {
		log.Error().Err(err).Msgf("can not get universities list with request %+v", request)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	ctx.Response().WriteHeader(http.StatusOK)
	if model == nil {
		return json.NewEncoder(ctx.Response()).Encode(api.PutResponse{
			University: models.University{},
		})
	}
	return json.NewEncoder(ctx.Response()).Encode(api.PutResponse{
		University: *model,
	})
}

func (a *adminer) addUniversity(request api.PutRequest) (*models.University, error) {
	uni, err := a.repo.AddUniversity(request)
	if uni == nil {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrap(err, "can not add university to list")
	}
	return uni, nil
}

func (a *adminer) DeleteUniversity(ctx echo.Context) error {
	idString := ctx.Param("id")
	if idString != "" {
		idInt, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}

		resp, err := a.deleteUniversity(idInt)
		ctx.Response().WriteHeader(http.StatusOK)
		if resp != nil {
			if err != nil {
				return json.NewEncoder(ctx.Response()).Encode(api.DeleteResponse{
					Whdb:  *resp,
					Error: err.Error(),
				})
			} else {
				return json.NewEncoder(ctx.Response()).Encode(api.DeleteResponse{
					Whdb:  *resp,
					Error: "",
				})
			}
		} else {
			return ctx.String(http.StatusInternalServerError, "something gone wrong")
		}
	}
	return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
}

func (a *adminer) deleteUniversity(id uint64) (*string, error) {
	return a.repo.DeleteUniversity(id)
}

func (a *adminer) PostIdAdmin(ctx echo.Context) error {
	id := ctx.ParamValues()
	if len(id) > 0 && id[0] != "" {
		idInt, err := strconv.ParseUint(id[0], 10, 64)
		if err != nil {
			return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}

		var request api.PostIdRequest
		err = ctx.Bind(&request)
		if err != nil {
			return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}

		resp, err := a.postIdAdmin(request, idInt)

		if err != nil {
			log.Error().Err(err).Msgf("can not update university with request %+v and id = %d", request, idInt)
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		ctx.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(ctx.Response()).Encode(api.PutResponse{
			University: *resp,
		})
	}
	return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
}

func (a *adminer) postIdAdmin(request api.PostIdRequest, id uint64) (*models.University, error) {
	return a.repo.UpdateUniversity(request, id)
}
