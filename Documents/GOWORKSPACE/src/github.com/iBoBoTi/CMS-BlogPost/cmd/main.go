package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/CMS-BlogPost/pkg/app"
	"github.com/iBoBoTi/CMS-BlogPost/pkg/repository"
)

func main() {
	// color logger
	gin.ForceConsoleColor()

	router := gin.Default()

	server := app.NewServer(router)

	// database connection initiated
	db,_ := repository.DataBaseConnection()
	defer db.Close()



	// starts server at 8080 default for gin
	err := server.RunServer()
	if err != nil {
		return 
	}

}
