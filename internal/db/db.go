package db

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ppl-study-planner/internal/model"
)

// Initialize opens the SQLite database and runs AutoMigrate
func Initialize() (*gorm.DB, error) {
	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, err
	}

	// Open SQLite database
	db, err := gorm.Open(sqlite.Open("data/data.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := db.AutoMigrate(
		&model.StudyPlan{},
		&model.DailyTask{},
		&model.Progress{},
		&model.ChecklistItem{},
		&model.Budget{},
	); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}
