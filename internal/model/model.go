package model

import (
	"time"
)

// StudyPlan represents a PPL study plan with a target checkride date
type StudyPlan struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	CheckrideDate time.Time   `gorm:"checkride_date" json:"checkride_date"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	DailyTasks    []DailyTask `gorm:"foreignKey:StudyPlanID" json:"daily_tasks,omitempty"`
}

// DailyTask represents a single task in the study plan
type DailyTask struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	StudyPlanID uint       `gorm:"study_plan_id" json:"study_plan_id"`
	Date        time.Time  `gorm:"date" json:"date"`
	Category    string     `gorm:"category" json:"category"`
	Title       string     `gorm:"title" json:"title"`
	Description string     `gorm:"description" json:"description"`
	Completed   bool       `gorm:"completed" json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	Progress    []Progress `gorm:"foreignKey:DailyTaskID" json:"progress,omitempty"`
}

// Progress tracks when a task was completed
type Progress struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	DailyTaskID uint      `gorm:"daily_task_id" json:"daily_task_id"`
	CompletedAt time.Time `gorm:"completed_at" json:"completed_at"`
}

// ChecklistCategory represents FAA ACS categories
type ChecklistCategory string

const (
	CategoryDocuments ChecklistCategory = "Documents"
	CategoryAircraft  ChecklistCategory = "Aircraft"
	CategoryGround    ChecklistCategory = "Ground"
	CategoryFlight    ChecklistCategory = "Flight"
)

// ChecklistItem represents an item in the pre-checkride checklist
type ChecklistItem struct {
	ID        uint              `gorm:"primaryKey" json:"id"`
	Category  ChecklistCategory `gorm:"category" json:"category"`
	Title     string            `gorm:"title" json:"title"`
	Completed bool              `gorm:"completed" json:"completed"`
	CreatedAt time.Time         `json:"created_at"`
}

// BudgetItemType represents types of budget items
type BudgetItemType string

const (
	BudgetFlightHours BudgetItemType = "flight_hours"
	BudgetPlaneRate   BudgetItemType = "plane_rate"
	BudgetCfiRate     BudgetItemType = "cfi_rate"
	BudgetLiving      BudgetItemType = "living"
)

// Budget tracks estimated and actual costs
type Budget struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ItemType  BudgetItemType `gorm:"item_type" json:"item_type"`
	Amount    float64        `gorm:"amount" json:"amount"`
	CreatedAt time.Time      `json:"created_at"`
}
