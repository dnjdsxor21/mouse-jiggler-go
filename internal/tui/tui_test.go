package tui

import (
	"errors"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
)

func TestMoveResultPrependsSuccessfulActivity(t *testing.T) {
	at := time.Date(2026, time.August, 25, 9, 30, 45, 0, time.Local)
	m := newModel(time.Minute, false, func() error { return nil }, func() time.Time { return at })

	updated, _ := m.Update(moveResult{at: at})
	got := updated.(model)

	if len(got.logs) != 1 || !got.logs[0].success || got.logs[0].detail != "moved + restored" {
		t.Fatalf("activity = %#v, want successful movement", got.logs)
	}
}

func TestMoveResultRecordsFailure(t *testing.T) {
	at := time.Date(2026, time.August, 25, 9, 30, 45, 0, time.Local)
	m := newModel(time.Minute, false, func() error { return nil }, func() time.Time { return at })

	updated, _ := m.Update(moveResult{at: at, err: errors.New("event unavailable")})
	got := updated.(model)

	if len(got.logs) != 1 || got.logs[0].success || got.logs[0].detail != "movement failed: event unavailable" {
		t.Fatalf("activity = %#v, want failed movement", got.logs)
	}
}

func TestMoveTickUpdatesNextMovement(t *testing.T) {
	at := time.Date(2026, time.August, 25, 9, 30, 45, 0, time.Local)
	m := newModel(90*time.Second, false, func() error { return nil }, func() time.Time { return at })

	updated, command := m.Update(moveTick(at))
	got := updated.(model)

	if command == nil {
		t.Fatal("move tick did not schedule movement")
	}
	want := at.Add(90 * time.Second)
	if !got.nextMove.Equal(want) {
		t.Fatalf("next move = %s, want %s", got.nextMove, want)
	}
}

func TestActivityLogIsBounded(t *testing.T) {
	m := newModel(time.Minute, false, func() error { return nil }, time.Now)
	for range maxLogEntries + 3 {
		m.addLog(activity{detail: "movement"})
	}

	if len(m.logs) != maxLogEntries {
		t.Fatalf("log length = %d, want %d", len(m.logs), maxLogEntries)
	}
}

func TestIntervalPromptStartsWithConfirmedSeconds(t *testing.T) {
	m := newModel(time.Minute, true, func() error { return nil }, time.Now)
	m.input = "2"

	updated, command := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got := updated.(model)

	if got.askingInterval || got.interval != 2*time.Second || !got.waitingFirst || command == nil {
		t.Fatalf("model = %#v, want running two-second model", got)
	}
}

func TestIntervalPromptRejectsZero(t *testing.T) {
	m := newModel(time.Minute, true, func() error { return nil }, time.Now)
	m.input = "0"

	updated, command := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got := updated.(model)

	if !got.askingInterval || got.inputError != "Enter a positive whole number of seconds." || command != nil {
		t.Fatalf("model = %#v, want validation error", got)
	}
}

func TestFormatRemainingUsesStableClockFormat(t *testing.T) {
	if got := formatRemaining(42*time.Second + time.Nanosecond); got != "00:43" {
		t.Fatalf("countdown = %q, want 00:43", got)
	}
	if got := formatRemaining(61 * time.Minute); got != "1:01:00" {
		t.Fatalf("countdown = %q, want 1:01:00", got)
	}
}
