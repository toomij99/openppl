package db

import (
	"testing"

	"ppl-study-planner/internal/model"
)

func TestInitializeAutoMigrate(t *testing.T) {
	database, err := Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	if !database.Migrator().HasTable(&model.AutomationIdempotency{}) {
		t.Fatal("expected automation_idempotencies table to exist")
	}
}
