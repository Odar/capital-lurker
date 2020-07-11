package server

func (s *server) setRoutes() {
	s.echo.GET("/just-for-test", s.receiver.Reverse)
}
