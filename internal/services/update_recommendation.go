package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const latestReleaseURL = "https://api.github.com/repos/toomij99/openppl/releases/latest"

type releasePayload struct {
	TagName string `json:"tag_name"`
}

func FetchLatestReleaseTag(timeout time.Duration) (string, error) {
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest(http.MethodGet, latestReleaseURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "openppl-update-check")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", fmt.Errorf("unexpected status: %s", res.Status)
	}

	var payload releasePayload
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return "", err
	}

	tag := strings.TrimSpace(payload.TagName)
	if !isSemverTag(tag) {
		return "", fmt.Errorf("invalid latest release tag: %q", tag)
	}

	return tag, nil
}

func BuildUpdateRecommendation(currentVersion string, latestVersion string) (string, bool) {
	currentVersion = strings.TrimSpace(currentVersion)
	latestVersion = strings.TrimSpace(latestVersion)
	if !isSemverTag(currentVersion) || !isSemverTag(latestVersion) {
		return "", false
	}

	cmp, ok := compareSemver(latestVersion, currentVersion)
	if !ok || cmp <= 0 {
		return "", false
	}

	msg := fmt.Sprintf(
		"Update available: %s (installed %s). Update with: curl -fsSL https://openppl.happycloud.ru/install | bash",
		latestVersion,
		currentVersion,
	)
	return msg, true
}

var semverRegex = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

func isSemverTag(v string) bool {
	return semverRegex.MatchString(strings.TrimSpace(v))
}

func compareSemver(a string, b string) (int, bool) {
	pa, ok := parseSemver(a)
	if !ok {
		return 0, false
	}
	pb, ok := parseSemver(b)
	if !ok {
		return 0, false
	}

	for i := 0; i < 3; i++ {
		if pa[i] > pb[i] {
			return 1, true
		}
		if pa[i] < pb[i] {
			return -1, true
		}
	}

	return 0, true
}

func parseSemver(v string) ([3]int, bool) {
	v = strings.TrimSpace(v)
	if !isSemverTag(v) {
		return [3]int{}, false
	}
	parts := strings.Split(strings.TrimPrefix(v, "v"), ".")
	if len(parts) != 3 {
		return [3]int{}, false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return [3]int{}, false
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return [3]int{}, false
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return [3]int{}, false
	}
	return [3]int{major, minor, patch}, true
}
