package services_test

import (
	"strings"
	"testing"

	"ppl-study-planner/internal/services"
)

func TestBuildUpdateRecommendation_NewerVersion(t *testing.T) {
	msg, ok := services.BuildUpdateRecommendation("v0.1.10", "v0.1.11")
	if !ok {
		t.Fatalf("expected recommendation for newer latest version")
	}
	if !strings.Contains(msg, "v0.1.11") || !strings.Contains(msg, "v0.1.10") {
		t.Fatalf("recommendation message missing versions: %q", msg)
	}
}

func TestBuildUpdateRecommendation_NotNewer(t *testing.T) {
	cases := []struct {
		current string
		latest  string
	}{
		{current: "v0.1.11", latest: "v0.1.11"},
		{current: "v0.1.12", latest: "v0.1.11"},
		{current: "(devel)", latest: "v0.1.11"},
		{current: "v0.1.11", latest: "latest"},
	}

	for _, tc := range cases {
		msg, ok := services.BuildUpdateRecommendation(tc.current, tc.latest)
		if ok {
			t.Fatalf("expected no recommendation for current=%q latest=%q, got %q", tc.current, tc.latest, msg)
		}
	}
}
