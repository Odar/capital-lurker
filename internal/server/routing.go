package server

func (s *server) setRoutes() {
	s.echo.GET("/just-for-test", s.receiver.Reverse)
	s.echo.GET("/speaker/on_main", s.speaker.GetSpeakersOnMain)
	s.echo.POST("/admin/speaker", s.speaker.GetSpeakersForAdmin)
	s.echo.DELETE("/admin/speaker/:id", s.speaker.DeleteSpeakerForAdmin)
	s.echo.POST("/admin/university", s.universityAdminer.GetUniversitiesList)
	s.echo.PUT("/admin/university", s.universityAdminer.AddUniversity)
	s.echo.DELETE("/admin/university/:id", s.universityAdminer.DeleteUniversity)
	s.echo.POST("/admin/university/:id", s.universityAdminer.UpdateUniversity)
	s.echo.File("/login", "assets/loginPage.html")
	s.echo.Any("/auth", s.authentificator.Login)
	//s.echo.File("/home", "assets/homePage.html")
	s.echo.GET("/home", s.authentificator.GetInfoFromVK)
}
