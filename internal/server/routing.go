package server

import "github.com/labstack/echo/v4/middleware"

func (s *server) setRoutes() {
	s.echo.GET("/just-for-test", s.receiver.Reverse)

	// Speaker entity
	s.echo.GET("/speaker/on_main", s.speaker.GetSpeakersOnMain)
	s.echo.POST("/admin/speaker", s.speaker.GetSpeakersForAdmin)
	s.echo.DELETE("/admin/speaker/:id", s.speaker.DeleteSpeakerForAdmin)
	s.echo.POST("/admin/speaker/:id", s.speaker.UpdateSpeakerForAdmin)
	s.echo.PUT("/admin/speaker", s.speaker.AddSpeakerForAdmin)

	// University entity
	s.echo.POST("/admin/university", s.universityAdminer.GetUniversitiesList)
	s.echo.PUT("/admin/university", s.universityAdminer.AddUniversity)
	s.echo.DELETE("/admin/university/:id", s.universityAdminer.DeleteUniversity)
	s.echo.POST("/admin/university/:id", s.universityAdminer.UpdateUniversity)

	// Regular auth
	s.echo.Any("/login", s.authenticator.Login)
	s.echo.Any("/signup", s.authenticator.SignUp)

	// Vk auth
	s.echo.Any("/loginvk", s.authenticator.LoginVkInitOauth)
	s.echo.GET("/tmppageforvkoauth", s.authenticator.LoginVkCheckRegistration) //set this url in callback url setting in vk/dev

	// Restricted
	jwtGroup := s.echo.Group("/signedinonly")
	jwtGroup.Use(middleware.JWT([]byte("Please, change me!")))
	jwtGroup.GET("/test", s.authenticator.TestPage)
	jwtGroup.POST("/logout", s.authenticator.Logout)

	// Theme entity
	s.echo.POST("/admin/theme", s.themeAdminer.GetThemesForAdmin)
	s.echo.DELETE("/admin/theme/:id", s.themeAdminer.DeleteThemeForAdmin)
	s.echo.POST("/admin/theme/:id", s.themeAdminer.UpdateThemeForAdmin)
	s.echo.PUT("/admin/theme", s.themeAdminer.AddThemeForAdmin)

	// Course entity
	s.echo.POST("/admin/course", s.courseAdminer.GetCoursesForAdmin)
	s.echo.DELETE("/admin/course/:id", s.courseAdminer.DeleteCourseForAdmin)
	s.echo.POST("/admin/course/:id", s.courseAdminer.UpdateCourseForAdmin)
	s.echo.PUT("/admin/course", s.courseAdminer.AddCourseForAdmin)
}
