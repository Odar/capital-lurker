package video

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
)

const (
	pathToSave = "/home/alexander/videoTest/"
)

type videodisc struct {
}

func New() *videodisc {
	return &videodisc{}
}

func (v *videodisc) UploadVideo(ctx echo.Context) error {
	id := ctx.FormValue("id")

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	src, err := file.Open()
	if err != nil {
		log.Error().Err(err).Msgf("can not open file to upload")
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	dst, err := os.Create(pathToSave + file.Filename + "_" + id)
	if err != nil {
		log.Error().Err(err).Msgf("can not open file to save")
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		log.Error().Err(err).Msgf("can not save file")
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	ctx.Response().WriteHeader(http.StatusOK)
	ctx.Response().Write([]byte("file saved successfully"))
	return nil
}
