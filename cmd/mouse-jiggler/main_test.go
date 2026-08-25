package main

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

type fakeMouse struct {
	trusted bool
}

func (m *fakeMouse) Trusted() bool {
	return m.trusted
}

func (m *fakeMouse) Jiggle() error {
	return nil
}

func noTUI(context.Context, time.Duration, bool, func() error) error {
	return nil
}

func TestRunPrintsVersionWithoutAccessibilityAccess(t *testing.T) {
	var stdout, stderr bytes.Buffer
	status := run([]string{"--version"}, &stdout, &stderr, &fakeMouse{}, noTUI)

	if status != 0 {
		t.Fatalf("run() status = %d, want 0; stderr = %q", status, stderr.String())
	}
	if got := strings.TrimSpace(stdout.String()); got != version {
		t.Fatalf("version = %q, want %q", got, version)
	}
}

func TestRunRejectsNonPositiveInterval(t *testing.T) {
	var stdout, stderr bytes.Buffer
	status := run([]string{"--interval=0s"}, &stdout, &stderr, &fakeMouse{}, noTUI)

	if status != 2 {
		t.Fatalf("run() status = %d, want 2", status)
	}
	if !strings.Contains(stderr.String(), "must be positive") {
		t.Fatalf("stderr = %q, want positive-interval message", stderr.String())
	}
}

func TestRunExplainsMissingAccessibilityAccess(t *testing.T) {
	var stdout, stderr bytes.Buffer
	status := run(nil, &stdout, &stderr, &fakeMouse{}, noTUI)

	if status != 1 {
		t.Fatalf("run() status = %d, want 1", status)
	}
	if !strings.Contains(stderr.String(), "Accessibility") {
		t.Fatalf("stderr = %q, want accessibility instructions", stderr.String())
	}
}

func TestRunStartsTUIWithRequestedInterval(t *testing.T) {
	var stdout, stderr bytes.Buffer
	var gotInterval time.Duration
	gotPrompt := true
	started := false
	start := func(_ context.Context, interval time.Duration, prompt bool, _ func() error) error {
		started = true
		gotInterval = interval
		gotPrompt = prompt
		return nil
	}

	status := run([]string{"--interval=90s"}, &stdout, &stderr, &fakeMouse{trusted: true}, start)
	if status != 0 {
		t.Fatalf("run() status = %d, want 0; stderr = %q", status, stderr.String())
	}
	if !started || gotInterval != 90*time.Second || gotPrompt {
		t.Fatalf("TUI started=%t interval=%s prompt=%t, want true, 1m30s, false", started, gotInterval, gotPrompt)
	}
}

func TestRunPromptsWhenIntervalIsOmitted(t *testing.T) {
	var stdout, stderr bytes.Buffer
	gotPrompt := false
	start := func(_ context.Context, _ time.Duration, prompt bool, _ func() error) error {
		gotPrompt = prompt
		return nil
	}

	status := run(nil, &stdout, &stderr, &fakeMouse{trusted: true}, start)
	if status != 0 {
		t.Fatalf("run() status = %d, want 0; stderr = %q", status, stderr.String())
	}
	if !gotPrompt {
		t.Fatal("TUI prompt = false, want true")
	}
}

func TestRunReportsTUIFailure(t *testing.T) {
	var stdout, stderr bytes.Buffer
	start := func(context.Context, time.Duration, bool, func() error) error {
		return errors.New("terminal failure")
	}

	status := run(nil, &stdout, &stderr, &fakeMouse{trusted: true}, start)
	if status != 1 {
		t.Fatalf("run() status = %d, want 1", status)
	}
	if !strings.Contains(stderr.String(), "terminal failure") {
		t.Fatalf("stderr = %q, want TUI error", stderr.String())
	}
}
