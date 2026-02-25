package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExportPaths_DefaultsToICSSDirectory(t *testing.T) {
	repoRoot := makeTestRepoRoot(t)

	resolved, err := resolveArtifactOutputDirFrom(filepath.Join(repoRoot, "internal", "services"), "")
	if err != nil {
		t.Fatalf("resolveArtifactOutputDirFrom failed: %v", err)
	}

	want := filepath.Join(repoRoot, "icss")
	if resolved != want {
		t.Fatalf("expected %q, got %q", want, resolved)
	}
}

func TestExportPaths_NormalizesRelativeDirectoriesUnderICSS(t *testing.T) {
	repoRoot := makeTestRepoRoot(t)
	base := filepath.Join(repoRoot, "internal")

	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "plain relative", input: "exports", want: filepath.Join(repoRoot, "icss", "exports")},
		{name: "prefixed with icss", input: filepath.Join("icss", "exports"), want: filepath.Join(repoRoot, "icss", "exports")},
		{name: "current dir", input: ".", want: filepath.Join(repoRoot, "icss")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resolved, err := resolveArtifactOutputDirFrom(base, tc.input)
			if err != nil {
				t.Fatalf("resolveArtifactOutputDirFrom failed: %v", err)
			}
			if resolved != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, resolved)
			}
		})
	}
}

func TestExportPaths_RejectsAbsolutePathsOutsideICSS(t *testing.T) {
	repoRoot := makeTestRepoRoot(t)

	outside := filepath.Join(repoRoot, "exports")
	if _, err := resolveArtifactOutputDirFrom(repoRoot, outside); err == nil {
		t.Fatal("expected error for absolute path outside icss")
	}
}

func TestExportPaths_RejectsTraversalOutsideICSS(t *testing.T) {
	repoRoot := makeTestRepoRoot(t)

	if _, err := resolveArtifactOutputDirFrom(repoRoot, "../outside"); err == nil {
		t.Fatal("expected traversal outside icss to be rejected")
	}
}

func makeTestRepoRoot(t *testing.T) string {
	t.Helper()

	repoRoot := t.TempDir()
	if err := os.WriteFile(filepath.Join(repoRoot, "go.mod"), []byte("module test\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	return repoRoot
}
