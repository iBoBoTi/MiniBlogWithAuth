package app

import (
	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/CMS-BlogPost/pkg/repository"
	"log"
	"os"
)

type Application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

type Server struct {
	router *gin.Engine
	DB     repository.Storage
}

func (s *Server) RunServer() error {
	// run function that calls routes method
	// runs the app at port 8081
	port := os.Getenv("PORT")

	r := s.Routes()

	// run the app through the router
	err := r.Run(":" + port)

	if err != nil {
		log.Printf("Error calling Run on router: %v", err)
		return err
	}

	return nil
}

func NewServer(router *gin.Engine, db repository.Storage) *Server {
	// initializes the app struct
	return &Server{
		router: router,
		DB:     db,
	}
}
