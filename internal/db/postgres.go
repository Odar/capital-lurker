package db

import (
	"fmt"

	"github.com/Odar/capital-lurker/pkg/db"
	"github.com/jmoiron/sqlx"
)

func New(cfg *db.Config) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", cfg.User, cfg.Password, cfg.Addr, cfg.Port, cfg.DB)
	conn, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	//conn.Exec("")

	return conn, nil
}
