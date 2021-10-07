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
	router.GET("/",index)
	router.GET("/login",handleLogin)
	router.GET("/signup",handleSignUp)
	router.POST("/login-auth",handleLoginAuth)
	router.POST("/signup-auth",handleSignUpAuth)
	return router
}
