package handlers

import (
	"net/http"

	"github.com/Ramsi-K/free-genai-bootcamp-2025/tree/main/projects/lang-portal/backend-go/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SettingsHandler struct {
	db *gorm.DB
}

func NewSettingsHandler(db *gorm.DB) *SettingsHandler {
	return &SettingsHandler{db: db}
}

func (h *SettingsHandler) ResetHistory(c *gin.Context) {
	// Delete all study sessions and word reviews
	if err := h.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.StudySession{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error resetting study sessions"})
		return
	}

	if err := h.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.WordReview{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error resetting word reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Study history reset successfully"})
}

func (h *SettingsHandler) FullReset(c *gin.Context) {
	// Drop all tables
	if err := h.db.Migrator().DropTable(
		&models.WordReview{},
		&models.StudySession{},
		&models.StudyActivity{},
		&models.WordsGroups{},
		&models.Word{},
		&models.Group{},
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error dropping tables"})
		return
	}

	// Auto-migrate to recreate tables
	if err := h.db.AutoMigrate(
		&models.Word{},
		&models.Group{},
		&models.WordsGroups{},
		&models.StudyActivity{},
		&models.StudySession{},
		&models.WordReview{},
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error recreating tables"})
		return
	}

	// Load seed data
	if err := h.loadSeedData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading seed data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database reset and reseeded successfully"})
}

func (h *SettingsHandler) loadSeedData() error {
	// Create default study activities
	activities := []models.StudyActivity{
		{
			Name:         "Flashcards",
			Description:  "Practice words with flashcards",
			Type:         "flashcards",
			ThumbnailURL: "/images/flashcards.png",
		},
		{
			Name:         "Multiple Choice",
			Description:  "Practice with multiple choice questions",
			Type:         "multiple_choice",
			ThumbnailURL: "/images/multiple-choice.png",
		},
		{
			Name:         "Sentence Practice",
			Description:  "Practice constructing sentences",
			Type:         "sentence_practice",
			ThumbnailURL: "/images/sentence-practice.png",
		},
	}

	for _, activity := range activities {
		if err := h.db.Create(&activity).Error; err != nil {
			return err
		}
	}

	// Note: Word and Group seed data should be loaded from external files
	// This is just a placeholder for the basic structure
	return nil
}
