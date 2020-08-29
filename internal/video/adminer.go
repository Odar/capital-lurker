package video

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Odar/capital-lurker/pkg/app/models"

	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/repositories"
	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func New(repo repositories.VideoAdminerRepo) *adminer {
	return &adminer{
		repo: repo,
	}
}

type adminer struct {
	repo repositories.VideoAdminerRepo
}

func (a *adminer) GetVideosForAdmin(ctx echo.Context) error {
	var request api.GetVideosForAdminRequest
	err := ctx.Bind(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not retrieve data from JSON:%+v", request)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	videosForAdmin, count, err := a.getVideosForAdmin(&request)
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

func (a *adminer) getVideosForAdmin(request *api.GetVideosForAdminRequest) ([]api.VideoForAdmin, uint64, error) {
	if request.Limit <= 0 {
		request.Limit = 10
	}
	if request.Page <= 0 {
		request.Page = 1
	}

	videos, err := a.repo.GetVideosForAdmin(request.Limit, request.Page, request.SortBy, &request.Filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "can not get videos from db")
	}

	count, err := a.repo.CountVideosForAdmin(&request.Filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "can not count videos from db")
	}

	if len(videos) == 0 {
		return nil, 0, nil
	}

	return videos, count, nil
}

func (a *adminer) DeleteVideoForAdmin(ctx echo.Context) error {
	var request api.DeleteVideoForAdminRequest
	var err error
	request.ID, err = strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msgf("can not retrieve id:%+v", ctx.Param("id"))
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	var WHBD string
	WHBD, err = a.deleteVideoForAdmin(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not delete video for admin with request %+v", request)
		ctx.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(ctx.Response()).Encode(api.DeleteVideoForAdminResponse{
			WHBD:  "error",
			Error: err.Error(),
		})
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(api.DeleteVideoForAdminResponse{
		WHBD:  WHBD,
		Error: "",
	})
}

func (a *adminer) deleteVideoForAdmin(request *api.DeleteVideoForAdminRequest) (string, error) {
	count, err := a.repo.DeleteVideo(request.ID)
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

func (a *adminer) UpdateVideoForAdmin(ctx echo.Context) error {
	var request api.UpdateVideoForAdminRequest
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

	updatedVideo, err := a.updateVideoForAdmin(requestID, &request)
	if err != nil {
		log.Error().Err(err).Msgf("can not update video for admin with request %+v", request)
		ctx.Response().WriteHeader(http.StatusInternalServerError)
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(updatedVideo)
}

func (a *adminer) updateVideoForAdmin(requestID uint64, request *api.UpdateVideoForAdminRequest) (
	*models.Video, error) {
	updatedVideo, err := a.repo.UpdateVideoForAdmin(requestID, request)
	if err != nil {
		return updatedVideo, errors.Wrap(err, "can not update in db")
	}

	return updatedVideo, nil
}

func (a *adminer) AddVideoForAdmin(ctx echo.Context) error {
	var request api.AddVideoForAdminRequest
	err := ctx.Bind(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not retrieve data from JSON:%+v", request)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	addedVideo, err := a.addVideoForAdmin(&request)
	if err != nil {
		log.Error().Err(err).Msgf("can not add video for admin with request %+v", request)
		ctx.Response().WriteHeader(http.StatusInternalServerError)
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(addedVideo)
}

func (a *adminer) addVideoForAdmin(request *api.AddVideoForAdminRequest) (*models.Video, error) {
	addedVideo, err := a.repo.AddVideoForAdmin(request)
	if err != nil {
		return addedVideo, errors.Wrap(err, "can not add into db")
	}

	return addedVideo, nil
}
