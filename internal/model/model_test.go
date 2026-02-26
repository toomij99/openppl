package model

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAutomationIdempotencyModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&AutomationIdempotency{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}

	first := AutomationIdempotency{
		ActionName:   "remind",
		RequestID:    "req-1",
		ArgsHash:     "hash-a",
		ActorScope:   "telegram:user:1",
		ResultState:  "executed",
		ResponseJSON: `{"result_state":"executed"}`,
	}
	if err := db.Create(&first).Error; err != nil {
		t.Fatalf("create first row: %v", err)
	}

	duplicate := AutomationIdempotency{
		ActionName:   "remind",
		RequestID:    "req-1",
		ArgsHash:     "hash-a",
		ActorScope:   "telegram:user:1",
		ResultState:  "executed",
		ResponseJSON: `{"result_state":"executed"}`,
	}
	if err := db.Create(&duplicate).Error; err == nil {
		t.Fatal("expected unique index violation for duplicate dedupe key")
	}
}
