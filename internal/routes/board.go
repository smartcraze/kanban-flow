package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcraze/kanban-flow/internal/handlers"
	"github.com/smartcraze/kanban-flow/internal/middleware"
)

func BoardRoutes(router *gin.Engine) {
	// boardroutes are protected routes, so they should be added after authentication routes
	board := router.Group("/api/boards")
	board.Use(middleware.AuthRequired())
	{
		board.GET("/", handlers.GetAllBoards)
		board.POST("/", handlers.CreateBoard)
		board.GET("/:id", handlers.GetBoardByID)
		board.PUT("/:id", handlers.UpdateBoard)
		board.DELETE("/:id", handlers.DeleteBoard)
		board.POST("/:id/members", handlers.AddMemberToBoard)
	}

}
