package db

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"ppl-study-planner/internal/model"
)

// Initialize opens the SQLite database and runs AutoMigrate
func Initialize() (*gorm.DB, error) {
	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, err
	}

	// Open SQLite database
	dbLogger := gormlogger.New(
		log.New(os.Stdout, "", log.LstdFlags),
		gormlogger.Config{
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(sqlite.Open("data/data.db"), &gorm.Config{Logger: dbLogger})
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
		&model.AppConfig{},
		&model.AutomationIdempotency{},
	); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}
