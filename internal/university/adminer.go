package university

import (
	"encoding/json"
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
	"github.com/Odar/capital-lurker/pkg/app/repositories"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
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

func (a *adminer) postAdmin(request api.PostRequest) ([]models.University, error) {
	content, err := a.repo.GetUniversities(request.Filter, request.SortBy) //nil Filter case

	if content == nil {
		return nil, err
	}

	if err != nil {
		return nil, errors.Wrap(err, "can not get from universities list")
	}

	limit, page := request.Limit, request.Page
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	if page-1 > len(content)/limit {
		return nil, errors.New("no such page")
	}

	model := make([]models.University, 0, limit)
	for i := 0; i < limit && (page-1)*limit+i < len(content); i++ {
		model = append(model, content[(page-1)*limit+i])
	}

	return model, err
}

func (a *adminer) PutAdmin(ctx echo.Context) error {
	var request api.PutRequest
	err := ctx.Bind(&request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	model, err := a.putAdmin(request)

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

func (a *adminer) putAdmin(request api.PutRequest) (*models.University, error) {
	uni, err := a.repo.AddUniversity(request) //nil Filter case

	if uni == nil {
		return nil, err
	}

	if err != nil {
		return nil, errors.Wrap(err, "can not add university to list")
	}

	return uni, nil
}

func (a *adminer) DeleteAdmin(ctx echo.Context) error {
	id := ctx.ParamValues()
	if len(id) > 0 && id[0] != "" {
		idInt, err := strconv.ParseUint(id[0], 10, 64)
		if err != nil {
			return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}

		resp, err := a.deleteAdmin(idInt)
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

func (a *adminer) deleteAdmin(id uint64) (*string, error) {
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
