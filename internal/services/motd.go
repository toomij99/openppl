package services

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

type MOTDQuizOption struct {
	Label string
	Text  string
}

type MOTDDailyQuiz struct {
	Date         string
	Entry        MOTDEntry
	Prompt       string
	Options      []MOTDQuizOption
	CorrectLabel string
	Explanation  string
}

// MOTDAnswer is the GORM model that records a user's daily quiz attempt.
type MOTDAnswer struct {
	ID             uint   `gorm:"primaryKey"`
	Date           string `gorm:"uniqueIndex;size:10"` // "2026-03-13"
	ACSCode        string `gorm:"size:20"`
	Prompt         string `gorm:"type:text"`
	SelectedOption string `gorm:"size:1"`
	CorrectOption  string `gorm:"size:1"`
	SelectedText   string `gorm:"type:text"`
	CorrectText    string `gorm:"type:text"`
	Answer         string `gorm:"type:text"`
	IsCorrect      bool
	Skipped        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type MOTDReadinessArea struct {
	Area     string
	Attempts int
	Correct  int
	Accuracy float64
}

type MOTDReadinessStats struct {
	TotalAttempts      int
	AnsweredAttempts   int
	CorrectAttempts    int
	SkippedAttempts    int
	OverallAccuracy    float64
	Last14Accuracy     float64
	CoverageScore      float64
	ReadinessScore     float64
	ReadinessLabel     string
	Areas              []MOTDReadinessArea
	WeakAreas          []MOTDReadinessArea
	TotalDistinctAreas int
}

type MOTDConfig struct {
	QuizMode bool `json:"quiz_mode"`
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

func BuildDailyQuiz(now time.Time) (MOTDDailyQuiz, error) {
	entry, err := TodaysACSCode(now)
	if err != nil {
		return MOTDDailyQuiz{}, err
	}

	tasks := loadMOTDTasks()
	if len(tasks) < 4 {
		return MOTDDailyQuiz{}, fmt.Errorf("insufficient ACS data to build quiz")
	}

	seed := quizSeed(now, entry.Code)
	correctText := nonEmptyObjective(entry)

	candidatePool := make([]MOTDEntry, 0, len(tasks))
	for _, t := range tasks {
		if t.Code == entry.Code {
			continue
		}
		if t.Section == entry.Section {
			candidatePool = append(candidatePool, t)
		}
	}
	if len(candidatePool) < 3 {
		candidatePool = candidatePool[:0]
		for _, t := range tasks {
			if t.Code != entry.Code {
				candidatePool = append(candidatePool, t)
			}
		}
	}

	rng := rand.New(rand.NewSource(seed))
	rng.Shuffle(len(candidatePool), func(i, j int) {
		candidatePool[i], candidatePool[j] = candidatePool[j], candidatePool[i]
	})

	optionTexts := []string{correctText}
	seen := map[string]struct{}{correctText: {}}
	for _, candidate := range candidatePool {
		text := nonEmptyObjective(candidate)
		if _, ok := seen[text]; ok {
			continue
		}
		seen[text] = struct{}{}
		optionTexts = append(optionTexts, text)
		if len(optionTexts) == 4 {
			break
		}
	}
	if len(optionTexts) < 4 {
		return MOTDDailyQuiz{}, fmt.Errorf("insufficient unique ACS objectives for quiz options")
	}

	rng.Shuffle(len(optionTexts), func(i, j int) {
		optionTexts[i], optionTexts[j] = optionTexts[j], optionTexts[i]
	})

	options := make([]MOTDQuizOption, 0, 4)
	correctLabel := ""
	for i, text := range optionTexts {
		label := string(rune('A' + i))
		if text == correctText {
			correctLabel = label
		}
		options = append(options, MOTDQuizOption{Label: label, Text: text})
	}

	return MOTDDailyQuiz{
		Date:         now.Format("2006-01-02"),
		Entry:        entry,
		Prompt:       fmt.Sprintf("Which objective best matches ACS %s?", entry.Code),
		Options:      options,
		CorrectLabel: correctLabel,
		Explanation:  fmt.Sprintf("%s focuses on: %s", entry.Code, correctText),
	}, nil
}

func NormalizeQuizChoice(input string) string {
	trimmed := strings.TrimSpace(strings.ToUpper(input))
	if trimmed == "" {
		return ""
	}
	for _, r := range trimmed {
		if r >= 'A' && r <= 'D' {
			return string(r)
		}
	}
	return ""
}

func IsCorrectQuizChoice(quiz MOTDDailyQuiz, choice string) bool {
	return choice != "" && choice == quiz.CorrectLabel
}

func QuizOptionText(quiz MOTDDailyQuiz, choice string) string {
	for _, option := range quiz.Options {
		if option.Label == choice {
			return option.Text
		}
	}
	return ""
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

func quizSeed(now time.Time, code string) int64 {
	h := fnv.New64a()
	day := now.UTC().Format("2006-01-02")
	_, _ = h.Write([]byte(day + ":" + code))
	return int64(h.Sum64())
}

func nonEmptyObjective(entry MOTDEntry) string {
	objective := strings.TrimSpace(entry.Text)
	if objective == "" || strings.EqualFold(objective, "[Archived]") {
		return "This ACS item is archived. Review current FAA ACS guidance and explain what changed from the previous standard."
	}
	return objective
}

func DefaultMOTDConfig() MOTDConfig {
	return MOTDConfig{QuizMode: true}
}

func LoadMOTDConfig() (MOTDConfig, error) {
	path, err := motdConfigPath()
	if err != nil {
		return DefaultMOTDConfig(), err
	}
	return loadMOTDConfigFromPath(path)
}

func SaveMOTDConfig(cfg MOTDConfig) error {
	path, err := motdConfigPath()
	if err != nil {
		return err
	}
	return saveMOTDConfigToPath(path, cfg)
}

func motdConfigPath() (string, error) {
	dir, err := motdDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "motd_config.json"), nil
}

func motdDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("motd: get home dir: %w", err)
	}
	dir := filepath.Join(home, ".openppl")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("motd: create dir: %w", err)
	}
	return dir, nil
}

func loadMOTDConfigFromPath(path string) (MOTDConfig, error) {
	cfg := DefaultMOTDConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("motd: read config: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return cfg, nil
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultMOTDConfig(), fmt.Errorf("motd: parse config: %w", err)
	}
	return cfg, nil
}

func saveMOTDConfigToPath(path string, cfg MOTDConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("motd: encode config: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("motd: write config: %w", err)
	}
	return nil
}

// InitMOTDDB opens (or creates) the per-user SQLite database used to store
// MOTD recall answers. It is stored at ~/.openppl/motd_answers.db so it never
// requires root privileges or a shared data directory.
func InitMOTDDB() (*gorm.DB, error) {
	dir, err := motdDataDir()
	if err != nil {
		return nil, err
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

// SaveMOTDAttempt upserts a quiz attempt for the given date.
func SaveMOTDAttempt(db *gorm.DB, date string, quiz MOTDDailyQuiz, selected string, skipped bool) error {
	if strings.TrimSpace(date) == "" {
		date = time.Now().Format("2006-01-02")
	}
	selectedText := QuizOptionText(quiz, selected)
	correctText := QuizOptionText(quiz, quiz.CorrectLabel)
	isCorrect := !skipped && IsCorrectQuizChoice(quiz, selected)
	legacyAnswer := selected
	if skipped {
		legacyAnswer = ""
	}

	var record MOTDAnswer
	result := db.Where(MOTDAnswer{Date: date}).
		Assign(MOTDAnswer{
			ACSCode:        quiz.Entry.Code,
			Prompt:         quiz.Prompt,
			SelectedOption: selected,
			CorrectOption:  quiz.CorrectLabel,
			SelectedText:   selectedText,
			CorrectText:    correctText,
			Answer:         legacyAnswer,
			IsCorrect:      isCorrect,
			Skipped:        skipped,
		}).
		FirstOrCreate(&record)
	if result.Error != nil {
		return fmt.Errorf("motd: save attempt: %w", result.Error)
	}
	// If the record already existed, FirstOrCreate won't apply Assign fields.
	if result.RowsAffected == 0 {
		record.ACSCode = quiz.Entry.Code
		record.Prompt = quiz.Prompt
		record.SelectedOption = selected
		record.CorrectOption = quiz.CorrectLabel
		record.SelectedText = selectedText
		record.CorrectText = correctText
		record.Answer = legacyAnswer
		record.IsCorrect = isCorrect
		record.Skipped = skipped
		if err := db.Save(&record).Error; err != nil {
			return fmt.Errorf("motd: update attempt: %w", err)
		}
	}
	return nil
}

func LoadMOTDAttempts(db *gorm.DB) ([]MOTDAnswer, error) {
	attempts := make([]MOTDAnswer, 0)
	if err := db.Order("date asc").Find(&attempts).Error; err != nil {
		return nil, fmt.Errorf("motd: load attempts: %w", err)
	}
	return attempts, nil
}

func ComputeMOTDReadiness(attempts []MOTDAnswer, now time.Time) MOTDReadinessStats {
	stats := MOTDReadinessStats{}
	stats.TotalAttempts = len(attempts)

	if len(attempts) == 0 {
		stats.ReadinessLabel = "Needs work"
		stats.TotalDistinctAreas = countDistinctAreas(loadMOTDTasks())
		return stats
	}

	byArea := map[string]*MOTDReadinessArea{}
	areaByCode := mapCodeToArea(loadMOTDTasks())
	cutoff14 := now.AddDate(0, 0, -14)
	correctRecentAreas := map[string]struct{}{}
	answered14 := 0
	correct14 := 0

	for _, attempt := range attempts {
		if attempt.Skipped {
			stats.SkippedAttempts++
		} else {
			stats.AnsweredAttempts++
			if attempt.IsCorrect {
				stats.CorrectAttempts++
			}
		}

		area := areaByCode[attempt.ACSCode]
		if area == "" {
			area = "unknown"
		}
		if _, ok := byArea[area]; !ok {
			byArea[area] = &MOTDReadinessArea{Area: area}
		}
		if !attempt.Skipped {
			byArea[area].Attempts++
			if attempt.IsCorrect {
				byArea[area].Correct++
			}
		}

		attemptDate, err := time.Parse("2006-01-02", attempt.Date)
		if err != nil {
			continue
		}
		if attemptDate.Before(cutoff14) {
			continue
		}
		if attempt.Skipped {
			continue
		}
		answered14++
		if attempt.IsCorrect {
			correct14++
			correctRecentAreas[area] = struct{}{}
		}
	}

	if stats.AnsweredAttempts > 0 {
		stats.OverallAccuracy = percent(stats.CorrectAttempts, stats.AnsweredAttempts)
	}
	if answered14 > 0 {
		stats.Last14Accuracy = percent(correct14, answered14)
	}

	totalAreas := countDistinctAreas(loadMOTDTasks())
	stats.TotalDistinctAreas = totalAreas
	if totalAreas > 0 {
		stats.CoverageScore = 100.0 * float64(len(correctRecentAreas)) / float64(totalAreas)
	}

	stats.ReadinessScore = 0.5*stats.OverallAccuracy + 0.3*stats.Last14Accuracy + 0.2*stats.CoverageScore
	switch {
	case stats.ReadinessScore >= 80:
		stats.ReadinessLabel = "Checkride-ready trend"
	case stats.ReadinessScore >= 60:
		stats.ReadinessLabel = "Building consistency"
	default:
		stats.ReadinessLabel = "Needs work"
	}

	areas := make([]MOTDReadinessArea, 0, len(byArea))
	for _, area := range byArea {
		if area.Attempts > 0 {
			area.Accuracy = percent(area.Correct, area.Attempts)
		}
		areas = append(areas, *area)
	}
	sort.Slice(areas, func(i, j int) bool {
		if areas[i].Area == areas[j].Area {
			return areas[i].Attempts > areas[j].Attempts
		}
		return areas[i].Area < areas[j].Area
	})
	stats.Areas = areas

	weak := make([]MOTDReadinessArea, 0, len(areas))
	for _, area := range areas {
		if area.Attempts == 0 {
			continue
		}
		weak = append(weak, area)
	}
	sort.Slice(weak, func(i, j int) bool {
		if weak[i].Accuracy == weak[j].Accuracy {
			return weak[i].Attempts > weak[j].Attempts
		}
		return weak[i].Accuracy < weak[j].Accuracy
	})
	if len(weak) > 5 {
		weak = weak[:5]
	}
	stats.WeakAreas = weak

	return stats
}

func mapCodeToArea(entries []MOTDEntry) map[string]string {
	codeToArea := make(map[string]string, len(entries))
	for _, entry := range entries {
		codeToArea[entry.Code] = entry.Area
	}
	return codeToArea
}

func countDistinctAreas(entries []MOTDEntry) int {
	areas := map[string]struct{}{}
	for _, entry := range entries {
		if strings.TrimSpace(entry.Area) == "" {
			continue
		}
		areas[entry.Area] = struct{}{}
	}
	return len(areas)
}

func percent(part int, total int) float64 {
	if total == 0 {
		return 0
	}
	return 100.0 * float64(part) / float64(total)
}
