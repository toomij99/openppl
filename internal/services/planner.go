package services

import (
	"time"

	"ppl-study-planner/internal/model"
)

// Task categories based on FAA ACS areas
var Categories = []string{
	"Theory",
	"Chair Flying",
	"Garmin 430",
	"CFI Flights",
}

// ACS Areas mapped to categories
// Theory: Areas 1-4 (Aeronautical Knowledge)
// Chair Flying: Areas 5-8 (Preflight Procedures, Airport Operations, Takeoffs/Landings, Fundamentals of Flight)
// Garmin 430: Areas 9-12 (Navigation, Slow Flight, Stalls, Basic Instrument)
// CFI Flights: Areas 13-16 (Complex Operations, Emergency Operations, Night, Post-flight)

var ACSAreasByCategory = map[string][]string{
	"Theory": {
		"Area 1: Aerodynamics",
		"Area 2: Regulations",
		"Area 3: Weather",
		"Area 4: Cross-Country Planning",
	},
	"Chair Flying": {
		"Area 5: Preflight Procedures",
		"Area 6: Airport Operations",
		"Area 7: Takeoffs and Landings",
		"Area 8: Fundamentals of Flight",
	},
	"Garmin 430": {
		"Area 9: Navigation",
		"Area 10: Slow Flight and Stalls",
		"Area 11: Basic Instrument",
		"Area 12: Attitude Instrument Flying",
	},
	"CFI Flights": {
		"Area 13: Complex Aircraft Operations",
		"Area 14: Emergency Operations",
		"Area 15: Night Operations",
		"Area 16: Post-Flight Procedures",
	},
}

// GenerateStudyPlan creates a backward-scheduled study plan
func GenerateStudyPlan(checkrideDate time.Time, totalDays int) []model.DailyTask {
	var tasks []model.DailyTask

	// Generate tasks for each day leading up to checkride
	for i := 0; i < totalDays; i++ {
		taskDate := checkrideDate.AddDate(0, 0, -i)

		// Skip dates in the past (before today)
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		if taskDate.Before(today) {
			continue
		}

		// Assign 1-2 tasks per day from different categories
		for catIndex, category := range Categories {
			// Distribute: some days have multiple tasks, some have fewer
			if (i+catIndex)%2 == 0 {
				task := createTaskForCategory(taskDate, category, i, totalDays)
				tasks = append(tasks, task)
			}
		}
	}

	return tasks
}

// createTaskForCategory creates a single task for a given date and category
func createTaskForCategory(date time.Time, category string, dayIndex, totalDays int) model.DailyTask {
	areas := ACSAreasByCategory[category]
	areaIndex := dayIndex % len(areas)
	area := areas[areaIndex]

	title := getTaskTitle(category, area)
	description := getTaskDescription(category, area, dayIndex)

	return model.DailyTask{
		Date:        date,
		Category:    category,
		Title:       title,
		Description: description,
		Completed:   false,
	}
}

// getTaskTitle returns a task title for the category and area
func getTaskTitle(category, area string) string {
	switch category {
	case "Theory":
		return area + " - Knowledge Review"
	case "Chair Flying":
		return area + " - Oral Prep"
	case "Garmin 430":
		return area + " - GPS/FMS Practice"
	case "CFI Flights":
		return area + " - Flight Maneuver"
	default:
		return area
	}
}

// getTaskDescription returns a detailed description for the task
func getTaskDescription(category, area string, dayIndex int) string {
	switch category {
	case "Theory":
		return "Review " + area + " concepts. Complete knowledge prep questions. Focus on FAA test prep."
	case "Chair Flying":
		return "Practice " + area + " procedures verbally. Visualize maneuvers. Review CFI teaching points."
	case "Garmin 430":
		return "Practice " + area + " on Garmin 430 simulator. Program flight plans. Review oceanic/continental navigation."
	case "CFI Flights":
		return "Prepare for " + area + " with CFI. Review ACS standards. Discuss common student errors."
	default:
		return "Complete " + area + " study materials."
	}
}

// DistributeTasksByCategory groups tasks by their category
func DistributeTasksByCategory(tasks []model.DailyTask) map[string][]model.DailyTask {
	distribution := make(map[string][]model.DailyTask)

	for _, task := range tasks {
		distribution[task.Category] = append(distribution[task.Category], task)
	}

	return distribution
}

// CalculateProgress calculates overall progress stats
func CalculateProgress(tasks []model.DailyTask) (completed int, total int, percentage float64) {
	total = len(tasks)
	for _, task := range tasks {
		if task.Completed {
			completed++
		}
	}

	if total > 0 {
		percentage = float64(completed) / float64(total) * 100
	}

	return completed, total, percentage
}

// GetProgressByCategory calculates progress per category
func GetProgressByCategory(tasks []model.DailyTask) map[string]ProgressStats {
	stats := make(map[string]ProgressStats)

	for _, category := range Categories {
		stats[category] = ProgressStats{
			Category:  category,
			Completed: 0,
			Total:     0,
		}
	}

	for _, task := range tasks {
		if s, exists := stats[task.Category]; exists {
			s.Total++
			if task.Completed {
				s.Completed++
			}
			stats[task.Category] = s
		}
	}

	return stats
}

// ProgressStats holds progress information for a category
type ProgressStats struct {
	Category   string
	Completed  int
	Total      int
	Percentage float64
}

// CalculatePercentage calculates the percentage for each category
func (ps ProgressStats) CalculatePercentage() float64 {
	if ps.Total > 0 {
		return float64(ps.Completed) / float64(ps.Total) * 100
	}
	return 0
}
