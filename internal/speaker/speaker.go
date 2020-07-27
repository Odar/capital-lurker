package speaker

import (
    "encoding/json"
    "github.com/Odar/capital-lurker/pkg/app/models"
    "github.com/pkg/errors"
    "github.com/rs/zerolog/log"
    "net/http"
    "strconv"

    "github.com/Odar/capital-lurker/pkg/app/repositories"

    "github.com/Odar/capital-lurker/pkg/api"
    "github.com/labstack/echo/v4"
)

func New(repo repositories.SpeakerRepo) *speaker {
    return &speaker{
        repo: repo,
    }
}

type speaker struct {
    repo repositories.SpeakerRepo
}

func (s *speaker) GetSpeakerOnMain(ctx echo.Context) error {
    var request api.GetSpeakerOnMainRequest
    err := ctx.Bind(&request)
    if err != nil {
        return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
    }

    speakersOnMain, err := s.getSpeakerOnMain(request)
    if err != nil {
        log.Error().Err(err).Msgf("can not get speaker on main with request %+v", request)
        return ctx.String(http.StatusInternalServerError, err.Error())
    }

    ctx.Response().WriteHeader(http.StatusOK)
    return json.NewEncoder(ctx.Response()).Encode(api.GetSpeakerOnMainResponse{
        Speakers: speakersOnMain,
    })
}

func (s *speaker) getSpeakerOnMain(request api.GetSpeakerOnMainRequest) ([]api.SpeakerOnMain, error) {
    speakersOnMain, err := s.repo.GetSpeakerOnMainFromDB(request.Limit)
    if err != nil {
        return nil, errors.Wrap(err, "can not get from db")
    }

    if speakersOnMain == nil {
        speakersOnMain = make([]api.SpeakerOnMain, 0)
    }

    return speakersOnMain, nil
}

func (s *speaker) GetSpeakerForAdmin(ctx echo.Context) error {
    var request api.GetSpeakerForAdminRequest
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
    return json.NewEncoder(ctx.Response()).Encode(api.GetSpeakerForAdminResponse{
        Speakers: speakersForAdmin,
        Count:    count,
    })
}

func (s *speaker) getSpeakerForAdmin(request *api.GetSpeakerForAdminRequest) ([]models.Speaker, uint64, error) {
    speakersForAdmin, count, err := s.repo.GetSpeakerForAdminFromDB(request)
    if err != nil {
        return nil, 0, errors.Wrap(err, "can not get from db")
    }

    if speakersForAdmin == nil {
        speakersForAdmin = make([]models.Speaker, 0)
        count = 0
    }

    return speakersForAdmin, count, nil
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
        log.Error().Err(err).Msgf("can not delete speaker for admin with request %+v", request)
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
    WHBD, err := s.repo.DeleteSpeakerForAdminFromDB(request.ID)
    if err != nil {
        return WHBD, errors.Wrap(err, "can not delete from db")
    }

    return WHBD, nil
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

    updateSpeaker, err := s.updateSpeakerForAdmin(requestID, &request)
    if err != nil {
        log.Error().Err(err).Msgf("can not update speaker for admin with request %+v", request)
        ctx.Response().WriteHeader(http.StatusInternalServerError)
    }

    ctx.Response().WriteHeader(http.StatusOK)
    return json.NewEncoder(ctx.Response()).Encode(updateSpeaker)
}

func (s *speaker) updateSpeakerForAdmin(requestID uint64, request *api.UpdateSpeakerForAdminRequest) (
    *models.Speaker, error) {
    updateSpeaker, err := s.repo.UpdateSpeakerForAdminInDB(requestID, request)
    if err != nil {
        return updateSpeaker, errors.Wrap(err, "can not update in db")
    }

    return updateSpeaker, nil
}
