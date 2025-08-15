package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcraze/kanban-flow/internal/handlers"
)

func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", handlers.RegisterUser)
		auth.POST("/login", handlers.LoginUser)
	}
}
