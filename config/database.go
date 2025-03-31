package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func ConnectDB()*gorm.DB{
	err := godotenv.Load()
	if err != nil{
		log.Printf("Cannot load env file %s", err)
	}

	dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    appEnv := os.Getenv("APP_ENV")


	if dbHost == "" || dbPort== "" || dbUser == "" || dbPassword == "" || dbName == "" || appEnv == ""{
		log.Fatal("Missing required environment variable")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to the database")
	}

	return db
}