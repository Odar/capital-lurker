package speaker

import (
	"encoding/json"
    "github.com/labstack/echo/v4"
    "net/http"

	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/repositories"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func New(repo repositories.SpeakerRepo) *speaker {
	return &speaker{
		repo: repo,
	}
}

type speaker struct {
	repo repositories.SpeakerRepo
}

func (s *speaker) GetSpeakersOnMain(ctx echo.Context) error {
	var request api.GetSpeakersOnMainRequest
	err := ctx.Bind(&request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	speakersOnMain, err := s.getSpeakersOnMain(request)
	if err != nil {
		log.Error().Err(err).Msgf("can not get speakers on main with request %+v", request)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(api.GetSpeakersOnMainResponse{
		Speakers: speakersOnMain,
	})
}

func (s *speaker) getSpeakersOnMain(request api.GetSpeakersOnMainRequest) ([]api.SpeakerOnMain, error) {
	speakersOnMain, err := s.repo.GetSpeakersOnMainFromDB(request.Limit)
	if err != nil {
		return nil, errors.Wrap(err, "can not get from db")
	}

	return speakersOnMain, nil
}