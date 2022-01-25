package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/CMS-BlogPost/pkg/app"
	"github.com/iBoBoTi/CMS-BlogPost/pkg/repository"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// color logger
	gin.ForceConsoleColor()
	er := godotenv.Load()
	if er != nil {
		log.Println(er.Error())
	}

	router := gin.Default()

	// database connection initiated
	db, errr := repository.DataBaseConnection()
	if errr != nil {
		panic(errr)
	}
	defer db.Close()
	DB := repository.Storage{DB: db}
	server := app.NewServer(router, DB)

	// starts server at 8080 default for gin
	err := server.RunServer()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

}
