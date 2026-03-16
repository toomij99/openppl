package services

import (
	"path/filepath"
	"testing"
)

func TestLoadMOTDConfig_DefaultWhenMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "motd_config.json")
	cfg, err := loadMOTDConfigFromPath(path)
	if err != nil {
		t.Fatalf("loadMOTDConfigFromPath returned error: %v", err)
	}
	if !cfg.QuizMode {
		t.Fatalf("default config should enable quiz mode")
	}
}

func TestSaveAndLoadMOTDConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "motd_config.json")
	if err := saveMOTDConfigToPath(path, MOTDConfig{QuizMode: false}); err != nil {
		t.Fatalf("saveMOTDConfigToPath returned error: %v", err)
	}

	cfg, err := loadMOTDConfigFromPath(path)
	if err != nil {
		t.Fatalf("loadMOTDConfigFromPath returned error: %v", err)
	}
	if cfg.QuizMode {
		t.Fatalf("expected quiz mode to be disabled after reload")
	}
}
