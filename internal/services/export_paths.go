package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var errRepoRootNotFound = errors.New("go.mod not found while resolving repository root")

// ResolveArtifactOutputDir normalizes artifact directories under <repo>/icss.
func ResolveArtifactOutputDir(outputDir string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	return resolveArtifactOutputDirFrom(cwd, outputDir)
}

func resolveArtifactOutputDirFrom(baseDir, outputDir string) (string, error) {
	repoRoot, err := findRepoRoot(baseDir)
	if err != nil {
		return "", err
	}

	artifactRoot := filepath.Join(repoRoot, "icss")
	trimmed := strings.TrimSpace(outputDir)

	var candidate string
	switch {
	case trimmed == "":
		candidate = artifactRoot
	case filepath.IsAbs(trimmed):
		candidate = filepath.Clean(trimmed)
	default:
		normalizedRel := normalizeRelativeArtifactDir(trimmed)
		candidate = filepath.Join(artifactRoot, normalizedRel)
	}

	relToArtifactRoot, err := filepath.Rel(artifactRoot, candidate)
	if err != nil {
		return "", fmt.Errorf("normalize artifact path: %w", err)
	}

	if relToArtifactRoot == ".." || strings.HasPrefix(relToArtifactRoot, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("artifact output must be under %q, got %q", artifactRoot, candidate)
	}

	return candidate, nil
}

func normalizeRelativeArtifactDir(raw string) string {
	cleaned := filepath.Clean(raw)
	if cleaned == "." {
		return ""
	}

	prefix := "icss" + string(filepath.Separator)
	cleaned = strings.TrimPrefix(cleaned, "."+string(filepath.Separator))
	if cleaned == "icss" {
		return ""
	}
	if strings.HasPrefix(cleaned, prefix) {
		return strings.TrimPrefix(cleaned, prefix)
	}

	return cleaned
}

func findRepoRoot(startDir string) (string, error) {
	current := filepath.Clean(startDir)
	for {
		candidate := filepath.Join(current, "go.mod")
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", errRepoRootNotFound
		}
		current = parent
	}
}
