package receiver

import (
	"encoding/json"
	"net/http"

	"github.com/Odar/capital-lurker/pkg/app/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/Odar/capital-lurker/pkg/app/repositories"

	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/labstack/echo/v4"
)

func New(repo repositories.ReceiverRepo) *receiver {
	return &receiver{
		repo: repo,
	}
}

type receiver struct {
	repo repositories.ReceiverRepo
}

func (r *receiver) Reverse(ctx echo.Context) error {
	var request api.ReverseRequest
	err := ctx.Bind(&request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	model, err := r.reverse(request)
	if err != nil {
		log.Error().Err(err).Msgf("can not reverse with request %+v", request)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(api.ReverseResponse{
		Word: model.Answer,
		Uses: model.Uses,
	})
}

func (r *receiver) reverse(request api.ReverseRequest) (*models.ReceiverCache, error) {
	model, err := r.repo.GetCache(request.Word)
	if err != nil {
		return nil, errors.Wrap(err, "can not get from cache")
	}

	if model == nil {
		reversed := reverseWord(request.Word)

		model = &models.ReceiverCache{
			Word:   request.Word,
			Answer: reversed,
			Uses:   1,
		}

		err = r.repo.SetCache(*model)
		if err != nil {
			return nil, errors.Wrap(err, "can not set in cache")
		}
	} else {
		err = r.repo.AddUsing(model)
		if err != nil {
			return nil, errors.Wrap(err, "can not add using for cache")
		}
	}

	return model, nil
}

func reverseWord(word string) string {
	n := 0
	inRunes := make([]rune, len(word))
	for _, r := range word {
		inRunes[n] = r
		n++
	}
	inRunes = inRunes[0:n]

	// Reverse
	for i := 0; i < n/2; i++ {
		inRunes[i], inRunes[n-1-i] = inRunes[n-1-i], inRunes[i]
	}
	// Convert back to UTF-8.
	return string(inRunes)
}
