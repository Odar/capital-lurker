package speaker

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

func (s *speaker) DeleteSpeakerForAdmin(ctx echo.Context) error {
	var request api.DeleteSpeakerForAdminRequest
	var err error
	request.ID, err = strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msgf("can not retrieve id:%+v", ctx.Param("id"))
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	var WHBD string
	WHBD, err = s.deleteSpeakerForAdmin(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not delete speaker for admin with request %+v", request) // Maybe some problemsх
		ctx.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(ctx.Response()).Encode(api.DeleteSpeakerForAdminResponse{
			WHBD:  "error",
			Error: err.Error(),
		})
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(api.DeleteSpeakerForAdminResponse{
		WHBD:  WHBD,
		Error: "",
	})
}

func (s *speaker) deleteSpeakerForAdmin(request *api.DeleteSpeakerForAdminRequest) (string, error) {
	count, err := s.repo.DeleteSpeaker(request.ID)
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

func (s *speaker) UpdateSpeakerForAdmin(ctx echo.Context) error {
	var request api.UpdateSpeakerForAdminRequest
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

	updatedSpeaker, err := s.updateSpeakerForAdmin(requestID, request)
	if err != nil {
		log.Error().Err(err).Msgf("can not update speaker for admin with request %+v", request)
		ctx.Response().WriteHeader(http.StatusInternalServerError)
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(updatedSpeaker)
}

func (s *speaker) updateSpeakerForAdmin(requestID uint64, request api.UpdateSpeakerForAdminRequest) (
	*models.Speaker, error) {
	updatedSpeaker, err := s.repo.UpdateSpeakerForAdmin(requestID, request)
	if err != nil {
		return updatedSpeaker, errors.Wrap(err, "can not update in db")
	}

	return updatedSpeaker, nil
}

func (s *speaker) AddSpeakerForAdmin(ctx echo.Context) error {
	var request api.AddSpeakerForAdminRequest
	err := ctx.Bind(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not retrieve data from JSON:%+v", request)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	addedSpeaker, err := s.addSpeakerForAdmin(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not add speaker for admin with request %+v", request)
		ctx.Response().WriteHeader(http.StatusInternalServerError)
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(addedSpeaker)
}

func (s *speaker) addSpeakerForAdmin(request *api.AddSpeakerForAdminRequest) (*models.Speaker, error) {
	addedSpeaker, err := s.repo.AddSpeakerForAdmin(request)
	if err != nil {
		return addedSpeaker, errors.Wrap(err, "can not add into db")
	}

	return addedSpeaker, nil
}
