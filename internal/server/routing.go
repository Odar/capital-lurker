package server

func (s *server) setRoutes() {
    s.echo.GET("/just-for-test", s.receiver.Reverse)
    s.echo.GET("/speaker/on_main", s.speaker.GetSpeakerOnMain)
    s.echo.POST("/admin/speaker/", s.speaker.GetSpeakerForAdmin)
    s.echo.DELETE("/admin/speaker/:id", s.speaker.DeleteSpeakerForAdmin)
    s.echo.POST("/admin/speaker/:id", s.speaker.UpdateSpeakerForAdmin)
}
