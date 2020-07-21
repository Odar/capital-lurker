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

	log.Info().Msgf("%+v", request)
	model, err := a.postAdmin(request)

	if err != nil {
		log.Error().Err(err).Msgf("can not get universities list with request %+v", request)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	ctx.Response().WriteHeader(http.StatusOK)

	return json.NewEncoder(ctx.Response()).Encode(api.PostResponse{
		Universities: model,
		Count:        uint64(len(model)),
	})
}

func (a *adminer) postAdmin(request api.PostRequest) ([]models.University, error) {
	content, err := a.repo.GetUniversities(request.Filter, request.SortBy) //nil Filter case

	if err != nil {
		return nil, errors.Wrap(err, "can not get from universities list")
	}
	//Можно ли так делать? Я просто думал про рефлексию, но данные здесь статические, поэтомурешил просто все кейсы разобрать
	/*
		switch request.SortBy {
		case "id":
			sort.Slice(content, func (i, j int) bool {return content[i].ID < content[j].ID})
		case "id-DESC":
			sort.Slice(content, func (i, j int) bool {return content[i].ID >= content[j].ID})
		case "name":
			sort.Slice(content, func (i, j int) bool {return content[i].Name < content[j].Name})
		case "name-DESC":
			sort.Slice(content, func (i, j int) bool {return content[i].Name >= content[j].Name})
		case "on_main_page":
			sort.Slice(content, func (i, j int) bool {
				if content[i].OnMainPage == content[j].OnMainPage {
					return true
				} else {
					return content[j].OnMainPage
				}
			})
		case "on_main_page-DESC":
			sort.Slice(content, func (i, j int) bool {
				if content[i].OnMainPage == content[j].OnMainPage {
					return true
				} else {
					return content[i].OnMainPage
				}
			})
		case "in_filter":
			sort.Slice(content, func (i, j int) bool {
				if content[i].OnMainPage == content[j].OnMainPage {
					return true
				} else {
					return content[j].OnMainPage
				}
			})
		case "in_filter-DESC":
			sort.Slice(content, func (i, j int) bool {
				if content[i].OnMainPage == content[j].OnMainPage {
					return true
				} else {
					return content[i].OnMainPage
				}
			})
		case "added_at":
			sort.Slice(content, func (i, j int) bool {return content[i].AddedAt.Before( content[j].AddedAt)})
		case "added_at-DESC":
			sort.Slice(content, func (i, j int) bool {return content[i].AddedAt.After(content[j].AddedAt) || content[i].AddedAt.Equal(content[j].AddedAt) })
		case "updated_at":
			sort.Slice(content, func (i, j int) bool {return content[i].UpdatedAt.Before( content[j].UpdatedAt)})
		case "updated_at-DESC":
			sort.Slice(content, func (i, j int) bool {return content[i].UpdatedAt.After(content[j].UpdatedAt) || content[i].UpdatedAt.Equal(content[j].UpdatedAt) })
		case "position":
			sort.Slice(content, func (i, j int) bool {return content[i].Position < content[j].Position})
		case "position-DESC":
			sort.Slice(content, func (i, j int) bool {return content[i].Position >= content[j].Position})
		}
	*/
	limit, page := request.Limit, request.Page
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	if page > len(content)/limit {
		return nil, errors.New("no such page")
	}

	model := make([]models.University, limit, limit)
	for i := 0; i < limit; i++ {
		model[i] = content[(page-1)*limit+i]
	}

	return model, err
}
