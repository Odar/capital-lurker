package app

type Server interface {
	Init() error
	Run() error
	Stop() error
}
