package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Odar/capital-lurker/internal/config"
)

func init() {
	flag.String("db", "settings", "db that we will use for migrations")
}

func main() {
	cfg := initConfig()

	dbConfig := cfg.CapitalDB
	postgres := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Addr, dbConfig.Port, dbConfig.DB)

	fmt.Printf(`-dir=%s postgres "%s" up`, dbConfig.MigrationDir, postgres)
}

func initConfig() config.Config {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return *cfg
}
