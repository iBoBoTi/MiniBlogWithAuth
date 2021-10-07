package repository

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

type Storage struct{
	DB *sql.DB
}


func DataBaseConnection() (*sql.DB, error){
	// sets up the database connection
	pswd := os.Getenv("MYSQL_PASSWORD")
	db, err := sql.Open("mysql","root:"+pswd+"@tcp(localhost:3306)/blog-cms")
	if err != nil{
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil{
		log.Println(err)
	}
	return db, nil
}