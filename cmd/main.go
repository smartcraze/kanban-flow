package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/smartcraze/kanban-flow/internal/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment")
	}
	database := db.ConnectDB()

	log.Println("⚙️Connected to DB:", database)

	if database == nil {
		log.Fatal("Database connection failed, exiting application")
	}

	app := gin.Default()
	app.Run(":8000")

}
