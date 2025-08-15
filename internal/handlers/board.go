package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smartcraze/kanban-flow/internal/db"
	"github.com/smartcraze/kanban-flow/internal/models"
)

// GetAllBoards retrieves all boards where the user is a member
func GetAllBoards(c *gin.Context) {
	userID := c.GetUint("userID")

	var boards []models.Board
	if err := db.DB.
		Joins("JOIN board_members bm ON bm.board_id = boards.id").
		Where("bm.user_id = ?", userID).
		Find(&boards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve boards"})
		return
	}

	c.JSON(http.StatusOK, boards)
}

// CreateBoard creates a new board and assigns the user as owner
func CreateBoard(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	board := models.Board{
		Title:       input.Title,
		Description: input.Description,
		OwnerID:     userID,
	}

	if err := db.DB.Create(&board).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create board"})
		return
	}

	// Add owner to board_members
	member := models.BoardMember{
		BoardID:  board.ID,
		UserID:   userID,
		Role:     "owner",
		JoinedAt: time.Now(),
	}
	db.DB.Create(&member)

	c.JSON(http.StatusCreated, board)
}

// GetBoardByID retrieves a specific board with its lists & cards
func GetBoardByID(c *gin.Context) {
	userID := c.GetUint("userID")
	boardID, _ := strconv.Atoi(c.Param("id"))

	// Check membership
	var member models.BoardMember
	if err := db.DB.
		Where("board_id = ? AND user_id = ?", boardID, userID).
		First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var board models.Board
	if err := db.DB.
		Preload("Lists.Cards").
		First(&board, boardID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Board not found"})
		return
	}

	c.JSON(http.StatusOK, board)
}

// UpdateBoard updates board details (owner & editors only)
func UpdateBoard(c *gin.Context) {
	userID := c.GetUint("userID")
	boardID, _ := strconv.Atoi(c.Param("id"))

	// Check role
	var member models.BoardMember
	if err := db.DB.
		Where("board_id = ? AND user_id = ?", boardID, userID).
		First(&member).Error; err != nil || (member.Role != "owner" && member.Role != "editor") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Model(&models.Board{}).
		Where("id = ?", boardID).
		Updates(models.Board{
			Title:       input.Title,
			Description: input.Description,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update board"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Board updated"})
}

// DeleteBoard deletes a board (owner only)
func DeleteBoard(c *gin.Context) {
	userID := c.GetUint("userID")
	boardID, _ := strconv.Atoi(c.Param("id"))

	// Check role
	var member models.BoardMember
	if err := db.DB.
		Where("board_id = ? AND user_id = ?", boardID, userID).
		First(&member).Error; err != nil || member.Role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only owner can delete"})
		return
	}

	if err := db.DB.Delete(&models.Board{}, boardID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete board"})
		return
	}

	// Also delete memberships
	db.DB.Where("board_id = ?", boardID).Delete(&models.BoardMember{})

	c.JSON(http.StatusOK, gin.H{"message": "Board deleted"})
}

// AddMemberToBoard adds a member to a board (owner only)
func AddMemberToBoard(c *gin.Context) {
	userID := c.GetUint("userID")
	boardID, _ := strconv.Atoi(c.Param("id"))

	// Check role
	var member models.BoardMember
	if err := db.DB.
		Where("board_id = ? AND user_id = ?", boardID, userID).
		First(&member).Error; err != nil || member.Role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only owner can add members"})
		return
	}

	var input struct {
		UserID uint   `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required,oneof=editor viewer"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newMember := models.BoardMember{
		BoardID:  uint(boardID),
		UserID:   input.UserID,
		Role:     input.Role,
		JoinedAt: time.Now(),
	}
	if err := db.DB.Create(&newMember).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added"})
}
