package speaker

import (
    "encoding/json"
    "net/http"

    "github.com/pkg/errors"
    "github.com/rs/zerolog/log"

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
