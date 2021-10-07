package app

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() *gin.Engine{
	// initializes the routes
	router := s.router

	// renders the templates
	router.LoadHTMLGlob("./ui/templates/*")

	// renders all static files
	router.Static("/static", "./ui/static")

	// application routes
	router.GET("/",CheckNotLoginMiddleware,index)
	router.GET("/login",CheckNotLoginMiddleware,handleLogin)
	router.GET("/signup",CheckNotLoginMiddleware,handleSignUp)
	router.POST("/login-auth",CheckNotLoginMiddleware,handleLoginAuth)
	router.POST("/signup-auth",CheckNotLoginMiddleware,handleSignUpAuth)

	GroupRoutes := router.Group("/blogar")
	{

		GroupRoutes.Use(CheckLoginMiddleware)

		GroupRoutes.GET("/",home)
		GroupRoutes.GET("/logout",handleLogOut)

	}
	return router
}
