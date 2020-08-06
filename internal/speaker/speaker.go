package speaker

import (
	"encoding/json"
	"net/http"

	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
	"github.com/Odar/capital-lurker/pkg/app/repositories"
	echo "github.com/labstack/echo/v4"
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

	speakers, err := s.getSpeakersOnMain(request)
	if err != nil {
		log.Error().Err(err).Msgf("can not get speakers on main with request %+v", request)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(api.GetSpeakersOnMainResponse{
		Speakers: speakers,
	})
}

func (s *speaker) getSpeakersOnMain(request api.GetSpeakersOnMainRequest) ([]api.SpeakerOnMain, error) {
	speakers, err := s.repo.GetSpeakersOnMain(request.Limit)
	if err != nil {
		return nil, errors.Wrap(err, "can not get from db")
	}

	return speakers, nil
}

func (s *speaker) GetSpeakersForAdmin(ctx echo.Context) error {
	var request api.GetSpeakersForAdminRequest
	err := ctx.Bind(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not retrieve data from JSON:%+v", request)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	speakersForAdmin, count, err := s.getSpeakerForAdmin(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not get speaker for admin with request %+v", request)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(api.GetSpeakersForAdminResponse{
		Speakers: speakersForAdmin,
		Count:    count,
	})
}

func (s *speaker) getSpeakerForAdmin(request *api.GetSpeakersForAdminRequest) ([]models.Speaker, uint64, error) {
	if request.Limit <= 0 {
		request.Limit = 10
	}
	if request.Page <= 0 {
		request.Page = 1
	}

	speakers, err := s.repo.GetSpeakersForAdmin(request.Limit, request.Page, request.SortBy, &request.Filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "can not get speakers from db")
	}

	count, err := s.repo.CountSpeakersForAdmin(&request.Filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "can not count speakers from db")
	}

	return speakers, count, nil
}
