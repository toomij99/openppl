package main

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

func TestRunWeb_UsesDefaultsAndStartsServer(t *testing.T) {
	restore := stubWebModeDeps(t)
	defer restore()

	called := false
	runWebServerFn = func(_ *gorm.DB, host string, port int) error {
		called = true
		if host != "127.0.0.1" {
			t.Fatalf("expected default hostname, got %q", host)
		}
		if port != 5016 {
			t.Fatalf("expected default port 5016, got %d", port)
		}
		return nil
	}

	if err := runWeb(nil); err != nil {
		t.Fatalf("runWeb failed: %v", err)
	}
	if !called {
		t.Fatal("expected web server to be invoked")
	}
}

func TestRunWeb_UsesCustomHostAndPort(t *testing.T) {
	restore := stubWebModeDeps(t)
	defer restore()

	runWebServerFn = func(_ *gorm.DB, host string, port int) error {
		if host != "0.0.0.0" {
			t.Fatalf("expected hostname 0.0.0.0, got %q", host)
		}
		if port != 7000 {
			t.Fatalf("expected port 7000, got %d", port)
		}
		return nil
	}

	if err := runWeb([]string{"--hostname", "0.0.0.0", "--port", "7000"}); err != nil {
		t.Fatalf("runWeb failed: %v", err)
	}
}

func TestRunWeb_InvalidHostname(t *testing.T) {
	restore := stubWebModeDeps(t)
	defer restore()

	err := runWeb([]string{"--hostname", ""})
	if err == nil {
		t.Fatal("expected hostname validation error")
	}
}

func TestRunWeb_InvalidPortLow(t *testing.T) {
	restore := stubWebModeDeps(t)
	defer restore()

	err := runWeb([]string{"--port", "0"})
	if err == nil {
		t.Fatal("expected low port validation error")
	}
}

func TestRunWeb_InvalidPortHigh(t *testing.T) {
	restore := stubWebModeDeps(t)
	defer restore()

	err := runWeb([]string{"--port", "65536"})
	if err == nil {
		t.Fatal("expected high port validation error")
	}
}

func TestRunWeb_UnknownFlag(t *testing.T) {
	restore := stubWebModeDeps(t)
	defer restore()

	err := runWeb([]string{"--nope"})
	if err == nil {
		t.Fatal("expected unknown flag error")
	}
}

func TestRunWeb_OnboardingGateRunsBeforeServer(t *testing.T) {
	restore := stubWebModeDeps(t)
	defer restore()

	called := false
	runOnboardingFn = func(force bool) error {
		called = true
		if !force {
			t.Fatal("expected forced onboarding")
		}
		return nil
	}
	needsSetupCheck = func() (bool, error) {
		return true, nil
	}

	if err := runWeb(nil); err != nil {
		t.Fatalf("runWeb failed: %v", err)
	}
	if !called {
		t.Fatal("expected onboarding to run")
	}
}

func TestNeedsSetup_Error(t *testing.T) {
	original := initDatabaseFn
	defer func() { initDatabaseFn = original }()

	initDatabaseFn = func() (*gorm.DB, error) {
		return nil, errors.New("boom")
	}

	_, err := needsSetup()
	if err == nil {
		t.Fatal("expected needsSetup error")
	}
}

func stubWebModeDeps(t *testing.T) func() {
	t.Helper()
	origNeeds := needsSetupCheck
	origOnboard := runOnboardingFn
	origInit := initDatabaseFn
	origRun := runWebServerFn

	needsSetupCheck = func() (bool, error) {
		return false, nil
	}
	runOnboardingFn = func(bool) error {
		return nil
	}
	initDatabaseFn = func() (*gorm.DB, error) {
		return &gorm.DB{}, nil
	}
	runWebServerFn = func(*gorm.DB, string, int) error {
		return nil
	}

	return func() {
		needsSetupCheck = origNeeds
		runOnboardingFn = origOnboard
		initDatabaseFn = origInit
		runWebServerFn = origRun
	}
}
