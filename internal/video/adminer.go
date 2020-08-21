package video

import (
	"encoding/json"
	"net/http"

	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/repositories"
	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func New(repo repositories.VideoAdminerRepo) *videoAdminer {
	return &videoAdminer{
		repo: repo,
	}
}

type videoAdminer struct {
	repo repositories.VideoAdminerRepo
}

func (v *videoAdminer) GetVideosForAdmin(ctx echo.Context) error {
	var request api.GetVideosForAdminRequest
	err := ctx.Bind(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not retrieve data from JSON:%+v", request)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	videosForAdmin, count, err := v.getVideosForAdmin(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not get video for admin with request %+v", request)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(api.GetVideosForAdminResponse{
		Videos: videosForAdmin,
		Count:  count,
	})
}

func (v *videoAdminer) getVideosForAdmin(request *api.GetVideosForAdminRequest) ([]api.VideoForAdmin, uint64, error) {
	if request.Limit <= 0 {
		request.Limit = 10
	}
	if request.Page <= 0 {
		request.Page = 1
	}

	videos, err := v.repo.GetVideosForAdmin(request.Limit, request.Page, request.SortBy, &request.Filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "can not get videos from db")
	}

	count, err := v.repo.CountVideosForAdmin(&request.Filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "can not count videos from db")
	}

	if len(videos) == 0 {
		return nil, 0, nil
	}

	return videos, count, nil
}
