package main

import (
	"github.com/Odar/capital-lurker/internal/config"
	"github.com/Odar/capital-lurker/internal/course"
	"github.com/Odar/capital-lurker/internal/db"
	"github.com/Odar/capital-lurker/internal/receiver"
	"github.com/Odar/capital-lurker/internal/server"
	"github.com/Odar/capital-lurker/internal/speaker"
	"github.com/Odar/capital-lurker/internal/theme"
	"github.com/Odar/capital-lurker/internal/university"
	"github.com/Odar/capital-lurker/internal/video"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := initConfig()
	setLogger(cfg.Logger)

	capitalDB, err := db.New(cfg.CapitalDB)
	if err != nil {
		log.Fatal().Err(err).Msg("can't connect to capital db")
	}

	receiverRepo := receiver.NewRepo(capitalDB)
	receiverService := receiver.New(receiverRepo)

	universityRepo := university.NewRepo(capitalDB)
	universityAdminerService := university.New(universityRepo)

	speakerRepo := speaker.NewRepo(capitalDB)
	speakerService := speaker.New(speakerRepo)

	themeRepo := theme.NewRepo(capitalDB)
	themeAdminerService := theme.New(themeRepo)

	courseRepo := course.NewRepo(capitalDB)
	courseAdminerService := course.New(courseRepo)

	videoRepo := video.NewRepo(capitalDB)
	videoAdminerService := video.New(videoRepo, *cfg.YouTubeClientSecret, *cfg.YouTubeVideoResource)

	srv := server.New(cfg.Server, receiverService, speakerService, universityAdminerService, themeAdminerService,
		courseAdminerService, videoAdminerService)
	err = srv.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("can not initialize server")
	}

	err = srv.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("server stopped with error...")
	}

	log.Fatal().Msg("server stopped...")
}

func initConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msgf("can not load config")
	}

	return cfg
}

func setLogger(cfg *config.Logger) {
	lvl, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		lvl = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(lvl)
}
