package onboarding

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"ppl-study-planner/internal/model"
	"ppl-study-planner/internal/services"
)

const (
	defaultPlanDays   = 90
	defaultPlaneRate  = 150.0
	defaultCfiRate    = 60.0
	defaultTravelCost = 500.0
	defaultBudget     = 10000.0
	defaultAirport    = "KFXE"

	modeByDate = "checkride_date"
	modeByDays = "planning_days"
)

type SetupValues struct {
	PlanningMode string

	CheckrideDate time.Time
	PlanDays      int

	SchoolAirport string
	WeatherNote   string
	BusyNote      string
	RiskNotes     []string

	PlaneRate   float64
	CfiRate     float64
	TravelCost  float64
	BudgetLimit float64
}

func NeedsOnboarding(db *gorm.DB) (bool, error) {
	var count int64
	if err := db.Model(&model.StudyPlan{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

func RunInteractiveSetup(db *gorm.DB, in io.Reader, out io.Writer, force bool) error {
	scanner := bufio.NewScanner(in)

	found, err := hasExistingStudyPlan(db)
	if err != nil {
		return err
	}

	if found && !force {
		fmt.Fprint(out, "An existing study plan was found. Reconfigure it now? [y/N]: ")
		if !scanner.Scan() {
			return scanner.Err()
		}
		answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if answer != "y" && answer != "yes" {
			fmt.Fprintln(out, "Onboarding cancelled.")
			return nil
		}
	}

	fmt.Fprintln(out, "\nopenppl onboarding")
	fmt.Fprintln(out, "Choose planning mode and answer setup questions. Press Enter to accept defaults.")

	values, err := collectValues(scanner, out)
	if err != nil {
		return err
	}

	if err := applySetup(db, values); err != nil {
		return err
	}

	fmt.Fprintln(out, "\nSetup complete.")
	if values.PlanningMode == modeByDate {
		fmt.Fprintf(out, "- Planning mode: by checkride date (%s)\n", values.CheckrideDate.Format("2006-01-02"))
	} else {
		fmt.Fprintf(out, "- Planning mode: by days (%d days -> %s)\n", values.PlanDays, values.CheckrideDate.Format("2006-01-02"))
	}
	fmt.Fprintf(out, "- Study window: today -> %s (%d days)\n", values.CheckrideDate.Format("2006-01-02"), values.PlanDays)
	fmt.Fprintf(out, "- School airport: %s\n", values.SchoolAirport)
	fmt.Fprintf(out, "- Weather outlook: %s\n", values.WeatherNote)
	fmt.Fprintf(out, "- Airport traffic outlook: %s\n", values.BusyNote)
	if len(values.RiskNotes) == 0 {
		fmt.Fprintln(out, "- Plan risk: Low (timeline looks healthy).")
	} else {
		fmt.Fprintln(out, "- Plan risk:")
		for _, risk := range values.RiskNotes {
			fmt.Fprintf(out, "  * %s\n", risk)
		}
	}
	fmt.Fprintf(out, "- Plane rate: $%.2f/hr | CFI rate: $%.2f/hr\n", values.PlaneRate, values.CfiRate)
	fmt.Fprintf(out, "- Budget limit: $%.2f\n", values.BudgetLimit)

	return nil
}

func hasExistingStudyPlan(db *gorm.DB) (bool, error) {
	var count int64
	if err := db.Model(&model.StudyPlan{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func collectValues(scanner *bufio.Scanner, out io.Writer) (SetupValues, error) {
	var v SetupValues

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	mode, err := promptPlanningMode(scanner, out)
	if err != nil {
		return v, err
	}
	v.PlanningMode = mode

	if mode == modeByDate {
		date, dateErr := promptFutureDate(scanner, out, "Checkride date (YYYY-MM-DD)", today.AddDate(0, 3, 0), today)
		if dateErr != nil {
			return v, dateErr
		}
		v.CheckrideDate = date
		v.PlanDays = planDaysFrom(today, date)
	} else {
		days, daysErr := promptInt(scanner, out, "Planning horizon in days", defaultPlanDays, 7, 730)
		if daysErr != nil {
			return v, daysErr
		}
		v.PlanDays = days
		v.CheckrideDate = today.AddDate(0, 0, days)
	}

	airport, err := promptAirport(scanner, out, "School airport ICAO", defaultAirport)
	if err != nil {
		return v, err
	}
	v.SchoolAirport = airport
	v.WeatherNote = weatherOutlook(airport, today)
	v.BusyNote = trafficOutlook(airport)

	planeRate, err := promptFloat(scanner, out, "Plane rental rate ($/hr)", defaultPlaneRate, 0, 1000)
	if err != nil {
		return v, err
	}

	cfiRate, err := promptFloat(scanner, out, "CFI rate ($/hr)", defaultCfiRate, 0, 500)
	if err != nil {
		return v, err
	}

	travelCost, err := promptFloat(scanner, out, "Travel setup cost ($)", defaultTravelCost, 0, 100000)
	if err != nil {
		return v, err
	}

	budget, err := promptFloat(scanner, out, "Total budget limit ($)", defaultBudget, 100, 1000000)
	if err != nil {
		return v, err
	}

	v.PlaneRate = planeRate
	v.CfiRate = cfiRate
	v.TravelCost = travelCost
	v.BudgetLimit = budget
	v.RiskNotes = riskAssessment(v.PlanDays)

	return v, nil
}

func promptPlanningMode(scanner *bufio.Scanner, out io.Writer) (string, error) {
	for {
		fmt.Fprintln(out, "Planning mode:")
		fmt.Fprintln(out, "  1) I know my checkride date")
		fmt.Fprintln(out, "  2) I want a fixed plan window (e.g. 90 days)")
		fmt.Fprint(out, "Select [1]: ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", err
			}
			return "", io.EOF
		}

		choice := strings.TrimSpace(scanner.Text())
		switch choice {
		case "", "1":
			return modeByDate, nil
		case "2":
			return modeByDays, nil
		default:
			fmt.Fprintln(out, "  Choose 1 or 2.")
		}
	}
}

func promptFutureDate(scanner *bufio.Scanner, out io.Writer, label string, defaultDate time.Time, notBefore time.Time) (time.Time, error) {
	for {
		fmt.Fprintf(out, "%s [%s]: ", label, defaultDate.Format("2006-01-02"))
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return time.Time{}, err
			}
			return time.Time{}, io.EOF
		}

		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			if defaultDate.Before(notBefore) {
				fmt.Fprintln(out, "  Date must be today or in the future.")
				continue
			}
			return defaultDate, nil
		}

		parsed, err := time.Parse("2006-01-02", text)
		if err != nil {
			fmt.Fprintln(out, "  Invalid date format. Use YYYY-MM-DD.")
			continue
		}
		if parsed.Before(notBefore) {
			fmt.Fprintln(out, "  Date must be today or in the future.")
			continue
		}
		return parsed, nil
	}
}

func promptAirport(scanner *bufio.Scanner, out io.Writer, label string, defaultValue string) (string, error) {
	for {
		fmt.Fprintf(out, "%s [%s]: ", label, defaultValue)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", err
			}
			return "", io.EOF
		}

		value := strings.ToUpper(strings.TrimSpace(scanner.Text()))
		if value == "" {
			return defaultValue, nil
		}

		if len(value) < 3 || len(value) > 4 {
			fmt.Fprintln(out, "  Airport code should be 3-4 letters (e.g. KFXE).")
			continue
		}

		valid := true
		for _, r := range value {
			if r < 'A' || r > 'Z' {
				valid = false
				break
			}
		}
		if !valid {
			fmt.Fprintln(out, "  Use letters only for airport code.")
			continue
		}

		return value, nil
	}
}

func promptInt(scanner *bufio.Scanner, out io.Writer, label string, defaultValue int, min int, max int) (int, error) {
	for {
		fmt.Fprintf(out, "%s [%d]: ", label, defaultValue)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return 0, err
			}
			return 0, io.EOF
		}

		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			return defaultValue, nil
		}

		value, err := strconv.Atoi(text)
		if err != nil || value < min || value > max {
			fmt.Fprintf(out, "  Enter a number between %d and %d.\n", min, max)
			continue
		}
		return value, nil
	}
}

func promptFloat(scanner *bufio.Scanner, out io.Writer, label string, defaultValue float64, min float64, max float64) (float64, error) {
	for {
		fmt.Fprintf(out, "%s [%.2f]: ", label, defaultValue)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return 0, err
			}
			return 0, io.EOF
		}

		text := strings.TrimSpace(strings.ReplaceAll(scanner.Text(), "$", ""))
		if text == "" {
			return defaultValue, nil
		}

		value, err := strconv.ParseFloat(text, 64)
		if err != nil || value < min || value > max {
			fmt.Fprintf(out, "  Enter a value between %.2f and %.2f.\n", min, max)
			continue
		}
		return value, nil
	}
}

func applySetup(db *gorm.DB, values SetupValues) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var plan model.StudyPlan
		err := tx.Order("id desc").Limit(1).Find(&plan).Error
		if err != nil {
			return err
		}

		if plan.ID == 0 {
			plan = model.StudyPlan{CheckrideDate: values.CheckrideDate}
			if err := tx.Create(&plan).Error; err != nil {
				return err
			}
		} else {
			plan.CheckrideDate = values.CheckrideDate
			if err := tx.Save(&plan).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("study_plan_id = ?", plan.ID).Delete(&model.DailyTask{}).Error; err != nil {
			return err
		}

		tasks := services.GenerateStudyPlan(values.CheckrideDate, values.PlanDays)
		for i := range tasks {
			tasks[i].StudyPlanID = plan.ID
		}
		if len(tasks) > 0 {
			if err := tx.Create(&tasks).Error; err != nil {
				return err
			}
		}

		if err := upsertBudgetItem(tx, model.BudgetPlaneRate, values.PlaneRate); err != nil {
			return err
		}
		if err := upsertBudgetItem(tx, model.BudgetCfiRate, values.CfiRate); err != nil {
			return err
		}
		if err := upsertBudgetItem(tx, model.BudgetLiving, values.TravelCost); err != nil {
			return err
		}
		if err := upsertBudgetItem(tx, model.BudgetLimit, values.BudgetLimit); err != nil {
			return err
		}

		if err := upsertConfig(tx, "planning_mode", values.PlanningMode); err != nil {
			return err
		}
		if err := upsertConfig(tx, "school_airport", values.SchoolAirport); err != nil {
			return err
		}
		if err := upsertConfig(tx, "weather_outlook", values.WeatherNote); err != nil {
			return err
		}
		if err := upsertConfig(tx, "traffic_outlook", values.BusyNote); err != nil {
			return err
		}
		if err := upsertConfig(tx, "risk_notes", strings.Join(values.RiskNotes, " | ")); err != nil {
			return err
		}

		return nil
	})
}

func upsertBudgetItem(tx *gorm.DB, itemType model.BudgetItemType, amount float64) error {
	var b model.Budget
	if err := tx.Where("item_type = ?", itemType).Order("id asc").Limit(1).Find(&b).Error; err != nil {
		return err
	}

	if b.ID == 0 {
		return tx.Create(&model.Budget{ItemType: itemType, Amount: amount}).Error
	}

	b.Amount = amount
	return tx.Save(&b).Error
}

func upsertConfig(tx *gorm.DB, key string, value string) error {
	var cfg model.AppConfig
	if err := tx.Where("key = ?", key).Limit(1).Find(&cfg).Error; err != nil {
		return err
	}

	if cfg.ID == 0 {
		return tx.Create(&model.AppConfig{Key: key, Value: value}).Error
	}

	cfg.Value = value
	return tx.Save(&cfg).Error
}

func planDaysFrom(today time.Time, checkride time.Time) int {
	start := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	end := time.Date(checkride.Year(), checkride.Month(), checkride.Day(), 0, 0, 0, 0, checkride.Location())
	if end.Before(start) {
		return 0
	}
	return int(end.Sub(start).Hours()/24) + 1
}

func riskAssessment(planDays int) []string {
	risks := make([]string, 0, 4)

	if planDays < 30 {
		risks = append(risks, "Very compressed timeline (<30 days). Expect high daily workload and reduced recovery days.")
	} else if planDays < 60 {
		risks = append(risks, "Compressed timeline (<60 days). Add at least one weekly buffer day for missed lessons/weather.")
	}

	avgTasksPerDay := 2.0
	if avgTasksPerDay > 1.8 && planDays < 75 {
		risks = append(risks, "Task density may feel heavy. Consider extending timeline or narrowing daily focus.")
	}

	if planDays > 180 {
		risks = append(risks, "Long timeline (>180 days) can reduce momentum. Add milestone check-ins every 2-3 weeks.")
	}

	return risks
}

func weatherOutlook(airport string, today time.Time) string {
	month := today.Month()

	florida := map[string]bool{
		"KFXE": true,
		"KFLL": true,
		"KPMP": true,
		"KMIA": true,
		"KOPF": true,
		"KPBI": true,
	}

	if florida[airport] {
		if month >= time.May && month <= time.October {
			return "South Florida wet season: expect afternoon convective storms, gust fronts, and frequent temporary delays."
		}
		return "South Florida dry season: generally flyable mornings with periodic windy afternoons and occasional frontal passage."
	}

	if strings.HasPrefix(airport, "K") {
		if month == time.December || month == time.January || month == time.February {
			return "Winter season: monitor icing ceilings/visibility trends and keep alternate training blocks ready."
		}
		if month >= time.June && month <= time.August {
			return "Summer season: expect heat turbulence and occasional thunderstorm constraints in the afternoon."
		}
	}

	return "Typical mixed conditions expected. Keep a rolling weather backup plan for 2-3 training sessions per week."
}

func trafficOutlook(airport string) string {
	busy := map[string]bool{
		"KFXE": true,
		"KFLL": true,
		"KMIA": true,
		"KJFK": true,
		"KLAX": true,
		"KATL": true,
		"KORD": true,
		"KSFO": true,
		"KLAS": true,
	}

	if busy[airport] {
		return "High traffic airport profile: peak flow often 07:00-10:00 and 16:00-19:00 local. Prefer mid-day blocks for training flights."
	}

	return "Moderate traffic profile expected. Keep flexibility around local events/weekends when pattern activity can spike."
}
