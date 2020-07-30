package server

func (s *server) setRoutes() {
	s.echo.GET("/just-for-test", s.receiver.Reverse)
	s.echo.GET("/speaker/on_main", s.speaker.GetSpeakersOnMain)
	s.echo.POST("/admin/university", s.universityAdminer.GetUniversitiesList)
	s.echo.PUT("/admin/university", s.universityAdminer.AddUniversity)
	s.echo.DELETE("/admin/university/*", s.universityAdminer.DeleteAdmin)
}
