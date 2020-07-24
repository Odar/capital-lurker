package db

type Config struct {
    Addr         string
    Port         uint16
    User         string
    Password     string
    DB           string
    MigrationDir string
}
