package app

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Application struct{
	errorLog *log.Logger
	infoLog *log.Logger
}

type Server struct{
	router *gin.Engine
}

func NewServer(router *gin.Engine) *Server{
	// initializes the server struct
	return &Server{
		router: router,
	}
}

func (s *Server) RunServer() error{
	// run function that calls routes method
	// runs the server at port 8081

	r := s.Routes()

	// run the server through the router
	err := r.Run("localhost:8081")

	if err != nil {
		log.Printf("Error calling Run on router: %v", err)
		return err
	}

	return nil
}


