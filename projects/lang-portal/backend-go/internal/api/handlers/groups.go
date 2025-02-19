package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ramsi-K/free-genai-bootcamp-2025/tree/main/projects/lang-portal/backend-go/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GroupHandler struct {
	db *gorm.DB
}

func NewGroupHandler(db *gorm.DB) *GroupHandler {
	return &GroupHandler{db: db}
}

func (h *GroupHandler) List(c *gin.Context) {
	var groups []models.Group
	if err := h.db.Where("deleted_at IS NULL").Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching groups"})
		return
	}

	c.JSON(http.StatusOK, groups)
}

func (h *GroupHandler) Get(c *gin.Context) {
	id := c.Param("id")

	// Validate ID is numeric
	if _, err := strconv.ParseUint(id, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var group models.Group
	if err := h.db.Where("id = ? AND deleted_at IS NULL", id).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching group"})
		return
	}

	c.JSON(http.StatusOK, group)
}

func (h *GroupHandler) GetWords(c *gin.Context) {
	id := c.Param("id")

	// Validate ID is numeric
	if _, err := strconv.ParseUint(id, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var group models.Group
	if err := h.db.Where("id = ? AND deleted_at IS NULL", id).
		Preload("Words", "deleted_at IS NULL").
		Preload("Words.Groups", "deleted_at IS NULL").
		First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching group"})
		return
	}

	// Extract group names for each word
	for i := range group.Words {
		groupNames := make([]string, len(group.Words[i].Groups))
		for j, g := range group.Words[i].Groups {
			groupNames[j] = g.Name
		}
		group.Words[i].WordGroups = groupNames
	}

	c.JSON(http.StatusOK, group.Words)
}

func (h *GroupHandler) GetStudySessions(c *gin.Context) {
	id := c.Param("id")

	// Validate ID is numeric
	if _, err := strconv.ParseUint(id, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// First check if group exists
	var group models.Group
	if err := h.db.Where("id = ? AND deleted_at IS NULL", id).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching group"})
		return
	}

	var sessions []models.StudySession
	if err := h.db.Where("group_id = ? AND deleted_at IS NULL", id).
		Preload("Activity").
		Preload("Reviews").
		Order("completed_at DESC").
		Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching study sessions"})
		return
	}

	// Format response with statistics
	response := make([]gin.H, len(sessions))
	for i, session := range sessions {
		stats := session.GetStats()
		response[i] = gin.H{
			"id":           session.ID,
			"completed_at": session.CompletedAt,
			"activity": gin.H{
				"id":   session.Activity.ID,
				"name": session.Activity.Name,
			},
			"stats": stats,
		}
	}

	c.JSON(http.StatusOK, response)
}
