package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/smartcraze/kanban-flow/internal/db"
	"github.com/smartcraze/kanban-flow/internal/models"
	"github.com/smartcraze/kanban-flow/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment")
	}
	database := db.ConnectDB()

	log.Println("⚙️Connected to DB:", database)

	migrationErr := database.AutoMigrate(&models.User{})

	if migrationErr != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	if database == nil {
		log.Fatal("Database connection failed, exiting application")
	}

	app := gin.Default()
	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Kanban Flow API"})
	})
	//user routes
	routes.UserRoutes(app)
	app.Run(":8000")

}
