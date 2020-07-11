package server

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (s *server) Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			contentLength := req.Header.Get(echo.HeaderContentLength)
			if contentLength == "" {
				contentLength = "0"
			}

			var logLevel zerolog.Level
			errorMessage := ""
			if err != nil {
				logLevel = zerolog.ErrorLevel
				errorMessage = err.Error()
			} else {
				logLevel = zerolog.InfoLevel
			}

			log.WithLevel(logLevel).
				Str("id", id).
				Str("remote_ip", c.RealIP()).
				Str("host", req.Host).
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Str("user_agent", req.UserAgent()).
				Int("status", res.Status).
				Str("latency", stop.Sub(start).String()).
				Str("bytes_in", contentLength).
				Int64("bytes_out", res.Size).
				Str("error", errorMessage).
				Msg("Echo server log")

			return
		}
	}
}
