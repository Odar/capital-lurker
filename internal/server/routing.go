package server

import "github.com/labstack/echo/v4/middleware"

const secret = ""

func (s *server) setRoutes() {
	s.echo.GET("/just-for-test", s.receiver.Reverse)
	s.echo.GET("/speaker/on_main", s.speaker.GetSpeakersOnMain)
	s.echo.POST("/admin/speaker", s.speaker.GetSpeakersForAdmin)
	s.echo.DELETE("/admin/speaker/:id", s.speaker.DeleteSpeakerForAdmin)
	s.echo.POST("/admin/university", s.universityAdminer.GetUniversitiesList)
	s.echo.PUT("/admin/university", s.universityAdminer.AddUniversity)
	s.echo.DELETE("/admin/university/:id", s.universityAdminer.DeleteUniversity)
	s.echo.POST("/admin/university/:id", s.universityAdminer.UpdateUniversity)

	//regular
	s.echo.Any("/login", s.authenticator.Login)
	s.echo.Any("/signup", s.authenticator.SignUp)

	//vk
	s.echo.Any("/loginvk", s.authenticator.LoginVkInitOauth)
	s.echo.GET("/tmppageforvkoauth", s.authenticator.LoginVkCheckRegistration) //set this url in callback url setting in vk/dev

	//restricted
	jwtGroup := s.echo.Group("/signedinonly")
	jwtGroup.Use(middleware.JWT([]byte("Please, change me!")))
	jwtGroup.GET("/test", s.authenticator.TestPage)
	jwtGroup.POST("/logout", s.authenticator.Logout)

}
