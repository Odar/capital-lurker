package server

func (s *server) setRoutes() {
	s.echo.GET("/just-for-test", s.receiver.Reverse)
	s.echo.POST("/admin/university/", s.universityAdminer.PostAdmin)
	s.echo.PUT("/admin/university/", s.universityAdminer.PutAdmin)
}
