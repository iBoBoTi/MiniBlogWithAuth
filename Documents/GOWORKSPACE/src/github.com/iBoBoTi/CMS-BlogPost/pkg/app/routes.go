package app

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() *gin.Engine {
	// initializes the routes
	router := s.router

	// renders the templates
	router.LoadHTMLGlob("./ui/templates/*")

	// renders all static files
	router.Static("/static", "./ui/static")

	// application routes
	router.GET("/", CheckNotLoginMiddleware, s.index)
	router.GET("/login", CheckNotLoginMiddleware, s.handleLogin)
	router.GET("/signup", CheckNotLoginMiddleware, s.handleSignUp)
	router.POST("/login-auth", CheckNotLoginMiddleware, s.handleLoginAuth)
	router.POST("/signup-auth", CheckNotLoginMiddleware, s.handleSignUpAuth)

	GroupRoutes := router.Group("/blogar")
	{

		GroupRoutes.Use(CheckLoginMiddleware)

		GroupRoutes.GET("/", s.home)
		GroupRoutes.GET("/logout", s.handleLogOut)
		GroupRoutes.GET("/create", s.handlePostCreate)
		GroupRoutes.POST("/add", s.handlePostCreateForm)
		GroupRoutes.GET("/edit/:id", s.handlePostEdit)
		GroupRoutes.POST("/post-edit/:id", s.handlePostEditForm)

		GroupRoutes.GET("/post/:id", s.handlePostRetrieve)
		GroupRoutes.GET("/post/delete/:id", s.handlePostDelete)
		GroupRoutes.GET("/my-post", s.handleUserPost)
		GroupRoutes.POST("/comment/:id", s.handleCommentCreateForm)

		// Comments Routes
	}
	return router
}
