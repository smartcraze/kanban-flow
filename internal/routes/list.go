package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcraze/kanban-flow/internal/handlers"
	"github.com/smartcraze/kanban-flow/internal/middleware"
)

func RegisterListRoutes(router *gin.Engine) {
	// Protect all list routes with AuthRequired middleware
	router.Use(middleware.AuthRequired())

	// Lists for a specific board
	router.GET("/api/boards/:boardId/lists", handlers.GetListsByBoard)
	router.POST("/api/boards/:boardId/lists", handlers.CreateList)

	// Operations on a single list
	router.PUT("/api/lists/:id", handlers.UpdateList)
	router.DELETE("/api/lists/:id", handlers.DeleteList)
	router.PUT("/api/lists/:id/position", handlers.UpdateListPosition)
}
