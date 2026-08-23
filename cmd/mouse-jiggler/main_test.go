package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type fakeMouse struct {
	trusted bool
	err     error
	calls   int
}

func (m *fakeMouse) Trusted() bool {
	return m.trusted
}

func (m *fakeMouse) Jiggle() error {
	m.calls++
	return m.err
}

func TestRunPrintsVersionWithoutAccessibilityAccess(t *testing.T) {
	var stdout, stderr bytes.Buffer
	status := run([]string{"--version"}, &stdout, &stderr, &fakeMouse{})

	if status != 0 {
		t.Fatalf("run() status = %d, want 0; stderr = %q", status, stderr.String())
	}
	if got := strings.TrimSpace(stdout.String()); got != version {
		t.Fatalf("version = %q, want %q", got, version)
	}
}

func TestRunRejectsNonPositiveInterval(t *testing.T) {
	var stdout, stderr bytes.Buffer
	status := run([]string{"--interval=0s"}, &stdout, &stderr, &fakeMouse{})

	if status != 2 {
		t.Fatalf("run() status = %d, want 2", status)
	}
	if !strings.Contains(stderr.String(), "must be positive") {
		t.Fatalf("stderr = %q, want positive-interval message", stderr.String())
	}
}

func TestRunExplainsMissingAccessibilityAccess(t *testing.T) {
	var stdout, stderr bytes.Buffer
	status := run(nil, &stdout, &stderr, &fakeMouse{})

	if status != 1 {
		t.Fatalf("run() status = %d, want 1", status)
	}
	if !strings.Contains(stderr.String(), "Accessibility") {
		t.Fatalf("stderr = %q, want accessibility instructions", stderr.String())
	}
}

func TestRunReturnsJiggleError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	pointer := &fakeMouse{trusted: true, err: errors.New("event failed")}
	status := run(nil, &stdout, &stderr, pointer)

	if status != 1 {
		t.Fatalf("run() status = %d, want 1", status)
	}
	if pointer.calls != 1 {
		t.Fatalf("Jiggle() calls = %d, want 1", pointer.calls)
	}
	if !strings.Contains(stderr.String(), "event failed") {
		t.Fatalf("stderr = %q, want jiggle error", stderr.String())
	}
}
