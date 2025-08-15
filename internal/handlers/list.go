package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smartcraze/kanban-flow/internal/db"
	"github.com/smartcraze/kanban-flow/internal/models"
)

// GET /api/boards/:boardId/lists
func GetListsByBoard(c *gin.Context) {
	userID := c.GetUint("userID")
	boardID, _ := strconv.Atoi(c.Param("boardId"))

	// Check if user has access to the board
	var member models.BoardMember
	if err := db.DB.
		Where("board_id = ? AND user_id = ?", boardID, userID).
		First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var lists []models.List
	if err := db.DB.
		Where("board_id = ?", boardID).
		Order("position ASC").
		Preload("Cards").
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve lists"})
		return
	}

	c.JSON(http.StatusOK, lists)
}

// POST /api/boards/:boardId/lists
func CreateList(c *gin.Context) {
	userID := c.GetUint("userID")
	boardID, _ := strconv.Atoi(c.Param("boardId"))

	// Permission check (owner/editor only)
	var member models.BoardMember
	if err := db.DB.
		Where("board_id = ? AND user_id = ?", boardID, userID).
		First(&member).Error; err != nil || (member.Role != "owner" && member.Role != "editor") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var input struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find max position to append new list at the end
	var maxPos int
	db.DB.Model(&models.List{}).Where("board_id = ?", boardID).Select("COALESCE(MAX(position), 0)").Scan(&maxPos)

	list := models.List{
		Title:    input.Title,
		BoardID:  uint(boardID),
		Position: maxPos + 1,
	}

	if err := db.DB.Create(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create list"})
		return
	}

	c.JSON(http.StatusCreated, list)
}

// PUT /api/lists/:id
func UpdateList(c *gin.Context) {
	userID := c.GetUint("userID")
	listID, _ := strconv.Atoi(c.Param("id"))

	// Find list & board
	var list models.List
	if err := db.DB.First(&list, listID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "List not found"})
		return
	}

	// Permission check (owner/editor only)
	var member models.BoardMember
	if err := db.DB.
		Where("board_id = ? AND user_id = ?", list.BoardID, userID).
		First(&member).Error; err != nil || (member.Role != "owner" && member.Role != "editor") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var input struct {
		Title string `json:"title"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Model(&list).Updates(models.List{
		Title: input.Title,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "List updated"})
}

// DELETE /api/lists/:id
func DeleteList(c *gin.Context) {
	userID := c.GetUint("userID")
	listID, _ := strconv.Atoi(c.Param("id"))

	// Find list
	var list models.List
	if err := db.DB.First(&list, listID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "List not found"})
		return
	}

	// Permission check (owner/editor only)
	var member models.BoardMember
	if err := db.DB.
		Where("board_id = ? AND user_id = ?", list.BoardID, userID).
		First(&member).Error; err != nil || (member.Role != "owner" && member.Role != "editor") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	if err := db.DB.Delete(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "List deleted"})
}

// PUT /api/lists/:id/position
func UpdateListPosition(c *gin.Context) {
	userID := c.GetUint("userID")
	listID, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Position int `json:"position" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var list models.List
	if err := db.DB.First(&list, listID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "List not found"})
		return
	}

	// Permission check (owner/editor only)
	var member models.BoardMember
	if err := db.DB.
		Where("board_id = ? AND user_id = ?", list.BoardID, userID).
		First(&member).Error; err != nil || (member.Role != "owner" && member.Role != "editor") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	list.Position = input.Position
	list.UpdatedAt = time.Now()

	if err := db.DB.Save(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update position"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Position updated"})
}
