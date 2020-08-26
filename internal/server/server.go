package server

import (
	"context"
	"fmt"

	"github.com/Odar/capital-lurker/pkg/app"
	echoprometheus "github.com/globocom/echo-prometheus"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(cfg *Config, receiver app.Receiver, speaker app.Speaker, universityAdminer app.UniversityAdminer,
	themeAdminer app.ThemeAdminer) *server {
	e := echo.New()
	e.Server.Addr = fmt.Sprintf(":%d", cfg.Port)

	return &server{
		echo:              e,
		receiver:          receiver,
		speaker:           speaker,
		universityAdminer: universityAdminer,
		themeAdminer:      themeAdminer,
		cfg:               cfg,
	}
}

type server struct {
	echo              *echo.Echo
	receiver          app.Receiver
	universityAdminer app.UniversityAdminer
	speaker           app.Speaker
	themeAdminer      app.ThemeAdminer
	cfg               *Config
}

func (s *server) Init() error {
	s.echo.HideBanner = true
	s.echo.HidePort = true

	var configMetrics = echoprometheus.NewConfig()
	configMetrics.Namespace = "pcj"
	configMetrics.NormalizeHTTPStatus = false
	configMetrics.Buckets = []float64{
		0.001, 0.005,
		0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08, 0.09,
		0.1, 0.15, 0.2, 0.25, 0.3, 0.35, 0.4, 0.45,
		0.5, 0.55, 0.6, 0.65, 0.7, 0.75, 0.8, 0.85, 0.9, 0.95,
		1, 2, 3, 5, 7, 10, 20, 30,
	}

	s.echo.Use(
		s.Logger(),
		echoprometheus.MetricsMiddlewareWithConfig(configMetrics),
		middleware.Recover(),
	)

	s.setRoutes()

	return nil
}

func (s *server) Run() error {
	e := s.echo
	return e.StartServer(e.Server)
}

func (s *server) Stop() error {
	return s.echo.Shutdown(context.Background())
}
