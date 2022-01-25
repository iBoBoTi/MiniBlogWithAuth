package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type Storage struct {
	DB *sql.DB
}

func DataBaseConnection() (*sql.DB, error) {
	//sets up the database connection
	DBUser := os.Getenv("DB_USER")
	DBPass := os.Getenv("DB_PASS")
	DBHost := os.Getenv("DB_HOST")
	DBName := os.Getenv("DB_NAME")
	DBPort := os.Getenv("DB_PORT")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPass, DBName)

	//dbURL := os.Getenv("DATABASE_URL")

	//DSN := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", DBUser, DBPass, DBHost, DBPort, DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Println(err)
	}
	return db, nil
}
