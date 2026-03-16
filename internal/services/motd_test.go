package services_test

import (
	"math"
	"strings"
	"testing"
	"time"

	"ppl-study-planner/internal/services"
)

// TestDailyCodeIndex_Deterministic verifies that the same date always produces
// the same index — the PRNG must be deterministic.
func TestDailyCodeIndex_Deterministic(t *testing.T) {
	date := time.Date(2026, 3, 13, 0, 0, 0, 0, time.UTC)
	idx1 := services.DailyCodeIndex(date, 960)
	idx2 := services.DailyCodeIndex(date, 960)
	if idx1 != idx2 {
		t.Errorf("DailyCodeIndex not deterministic: got %d then %d for same date", idx1, idx2)
	}
}

// TestDailyCodeIndex_DifferentDays verifies that different calendar days produce
// different indices. Note: there is a 1/960 chance of a false failure when two
// consecutive days happen to hash to the same index — this is acceptable.
func TestDailyCodeIndex_DifferentDays(t *testing.T) {
	day1 := time.Date(2026, 3, 13, 0, 0, 0, 0, time.UTC)
	day2 := time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC)
	idx1 := services.DailyCodeIndex(day1, 960)
	idx2 := services.DailyCodeIndex(day2, 960)
	if idx1 == idx2 {
		t.Errorf("DailyCodeIndex returned same index %d for different days 2026-03-13 and 2026-03-14 (1/960 chance of false failure)", idx1)
	}
}

// TestDailyCodeIndex_ZeroTotal verifies that a zero totalCodes returns 0 without panicking.
func TestDailyCodeIndex_ZeroTotal(t *testing.T) {
	date := time.Date(2026, 3, 13, 0, 0, 0, 0, time.UTC)
	result := services.DailyCodeIndex(date, 0)
	if result != 0 {
		t.Errorf("DailyCodeIndex(date, 0) = %d; want 0", result)
	}
}

// TestTodaysACSCode_ReturnsEntry verifies that TodaysACSCode returns a valid entry.
func TestTodaysACSCode_ReturnsEntry(t *testing.T) {
	entry, err := services.TodaysACSCode(time.Now())
	if err != nil {
		t.Fatalf("TodaysACSCode returned unexpected error: %v", err)
	}
	if entry.Code == "" {
		t.Error("TodaysACSCode returned empty Code")
	}
	if entry.Text == "" {
		t.Error("TodaysACSCode returned empty Text")
	}
}

func TestBuildDailyQuiz_DeterministicForSameDay(t *testing.T) {
	now := time.Date(2026, 3, 16, 8, 30, 0, 0, time.UTC)
	quiz1, err := services.BuildDailyQuiz(now)
	if err != nil {
		t.Fatalf("BuildDailyQuiz returned error: %v", err)
	}
	quiz2, err := services.BuildDailyQuiz(now)
	if err != nil {
		t.Fatalf("BuildDailyQuiz returned error on second call: %v", err)
	}

	if quiz1.Entry.Code != quiz2.Entry.Code {
		t.Fatalf("quiz code mismatch for same day: %s vs %s", quiz1.Entry.Code, quiz2.Entry.Code)
	}
	if quiz1.CorrectLabel != quiz2.CorrectLabel {
		t.Fatalf("correct label mismatch for same day: %s vs %s", quiz1.CorrectLabel, quiz2.CorrectLabel)
	}
	if len(quiz1.Options) != 4 || len(quiz2.Options) != 4 {
		t.Fatalf("expected 4 options, got %d and %d", len(quiz1.Options), len(quiz2.Options))
	}
	for i := range quiz1.Options {
		if quiz1.Options[i].Label != quiz2.Options[i].Label || quiz1.Options[i].Text != quiz2.Options[i].Text {
			t.Fatalf("option %d mismatch for same day", i)
		}
	}
}

func TestNormalizeQuizChoice(t *testing.T) {
	cases := map[string]string{
		"a":        "A",
		"  B  ":    "B",
		"option c": "C",
		"D)":       "D",
		"":         "",
		"skip":     "",
		"  9   ":   "",
	}

	for input, want := range cases {
		got := services.NormalizeQuizChoice(input)
		if got != want {
			t.Fatalf("NormalizeQuizChoice(%q) = %q; want %q", input, got, want)
		}
	}
}

func TestComputeMOTDReadiness(t *testing.T) {
	attempts := []services.MOTDAnswer{
		{Date: "2026-03-05", ACSCode: "PA.I.A.K1", IsCorrect: true, Skipped: false},
		{Date: "2026-03-06", ACSCode: "PA.I.A.K2", IsCorrect: false, Skipped: false},
		{Date: "2026-03-10", ACSCode: "PA.II.A.K1", IsCorrect: true, Skipped: false},
		{Date: "2026-03-12", ACSCode: "PA.II.A.K2", IsCorrect: true, Skipped: false},
		{Date: "2026-03-13", ACSCode: "PA.III.A.K1", IsCorrect: false, Skipped: true},
	}

	now := time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC)
	stats := services.ComputeMOTDReadiness(attempts, now)

	if stats.TotalAttempts != 5 {
		t.Fatalf("TotalAttempts = %d; want 5", stats.TotalAttempts)
	}
	if stats.AnsweredAttempts != 4 {
		t.Fatalf("AnsweredAttempts = %d; want 4", stats.AnsweredAttempts)
	}
	if stats.CorrectAttempts != 3 {
		t.Fatalf("CorrectAttempts = %d; want 3", stats.CorrectAttempts)
	}
	if stats.SkippedAttempts != 1 {
		t.Fatalf("SkippedAttempts = %d; want 1", stats.SkippedAttempts)
	}

	if math.Abs(stats.OverallAccuracy-75.0) > 0.001 {
		t.Fatalf("OverallAccuracy = %.3f; want 75.0", stats.OverallAccuracy)
	}
	if math.Abs(stats.Last14Accuracy-75.0) > 0.001 {
		t.Fatalf("Last14Accuracy = %.3f; want 75.0", stats.Last14Accuracy)
	}
	if stats.ReadinessScore <= 0 {
		t.Fatalf("ReadinessScore should be > 0, got %.3f", stats.ReadinessScore)
	}
	if strings.TrimSpace(stats.ReadinessLabel) == "" {
		t.Fatalf("ReadinessLabel should not be empty")
	}
}

// TestDailyCodeIndex_BoundsCheck verifies that for 365 consecutive days all
// returned indices are in [0, 960) and that at least 300 distinct values appear
// (distribution check — confirms the PRNG actually varies across days).
func TestDailyCodeIndex_BoundsCheck(t *testing.T) {
	const totalCodes = 960
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	seen := make(map[int]struct{})

	for i := 0; i < 365; i++ {
		date := start.AddDate(0, 0, i)
		idx := services.DailyCodeIndex(date, totalCodes)
		if idx < 0 || idx >= totalCodes {
			t.Errorf("day %d: index %d out of range [0, %d)", i, idx, totalCodes)
		}
		seen[idx] = struct{}{}
	}

	if len(seen) < 300 {
		t.Errorf("only %d distinct indices across 365 days; want at least 300 (distribution too poor)", len(seen))
	}
}
