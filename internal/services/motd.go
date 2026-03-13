package services

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	gormlogger "gorm.io/gorm/logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed acs_private_airplane_6c.json
var motdDatasetJSON []byte

var (
	motdTasksOnce sync.Once
	motdTasks     []MOTDEntry
)

// MOTDEntry contains the ACS item fields needed by login-time MOTD flow.
type MOTDEntry struct {
	Code     string `json:"code"`
	Area     string `json:"area"`
	Task     string `json:"task"`
	Title    string `json:"task_title"`
	Section  string `json:"section"`
	Index    int    `json:"index"`
	Text     string `json:"text"`
	Category string `json:"category"`
}

// MOTDAnswer is the GORM model that records a user's daily recall answer.
type MOTDAnswer struct {
	ID        uint   `gorm:"primaryKey"`
	Date      string `gorm:"uniqueIndex;size:10"` // "2026-03-13"
	ACSCode   string `gorm:"size:20"`
	Answer    string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DailyCodeIndex returns a deterministic index for the given date.
// The same calendar day always produces the same index.
// Returns 0 when totalCodes <= 0.
func DailyCodeIndex(date time.Time, totalCodes int) int {
	if totalCodes <= 0 {
		return 0
	}

	dayStartUTC := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dayNumber := int(dayStartUTC.Unix() / 86400)

	step := totalCodes - 1
	for step > 1 && gcd(step, totalCodes) != 1 {
		step--
	}
	if gcd(step, totalCodes) != 1 {
		step = 1
	}

	idx := (dayNumber*step + 131) % totalCodes
	if idx < 0 {
		idx += totalCodes
	}

	return idx
}

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// TodaysACSCode returns the ACS task selected for the given date.
// It is deterministic: the same date always returns the same entry.
func TodaysACSCode(now time.Time) (MOTDEntry, error) {
	tasks := loadMOTDTasks()
	if len(tasks) == 0 {
		return MOTDEntry{}, fmt.Errorf("ACS data is empty")
	}
	idx := DailyCodeIndex(now, len(tasks))
	return tasks[idx], nil
}

func loadMOTDTasks() []MOTDEntry {
	motdTasksOnce.Do(func() {
		parsed := make([]MOTDEntry, 0)
		if err := json.Unmarshal(motdDatasetJSON, &parsed); err != nil {
			motdTasks = nil
			return
		}
		motdTasks = parsed
	})

	return motdTasks
}

// InitMOTDDB opens (or creates) the per-user SQLite database used to store
// MOTD recall answers. It is stored at ~/.openppl/motd_answers.db so it never
// requires root privileges or a shared data directory.
func InitMOTDDB() (*gorm.DB, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("motd: get home dir: %w", err)
	}
	dir := filepath.Join(home, ".openppl")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("motd: create dir: %w", err)
	}
	dbPath := filepath.Join(dir, "motd_answers.db")
	silentLogger := gormlogger.New(
		log.New(os.Stderr, "", 0),
		gormlogger.Config{LogLevel: gormlogger.Silent},
	)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: silentLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("motd: open db: %w", err)
	}
	if err := db.AutoMigrate(&MOTDAnswer{}); err != nil {
		return nil, fmt.Errorf("motd: migrate: %w", err)
	}
	return db, nil
}

// SaveMOTDAnswer upserts a recall answer for today's date. If an answer for
// today already exists it is overwritten; otherwise a new record is created.
func SaveMOTDAnswer(db *gorm.DB, code string, answer string) error {
	date := time.Now().Format("2006-01-02")
	var record MOTDAnswer
	result := db.Where(MOTDAnswer{Date: date}).
		Assign(MOTDAnswer{ACSCode: code, Answer: answer}).
		FirstOrCreate(&record)
	if result.Error != nil {
		return fmt.Errorf("motd: save answer: %w", result.Error)
	}
	// If the record already existed, FirstOrCreate won't apply Assign fields —
	// update it explicitly.
	if result.RowsAffected == 0 {
		record.ACSCode = code
		record.Answer = answer
		if err := db.Save(&record).Error; err != nil {
			return fmt.Errorf("motd: update answer: %w", err)
		}
	}
	return nil
}
