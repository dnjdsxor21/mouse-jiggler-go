package tui

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestMoveResultPrependsSuccessfulActivity(t *testing.T) {
	at := time.Date(2026, time.August, 25, 9, 30, 45, 0, time.Local)
	m := newModel(time.Minute, func() error { return nil }, func() time.Time { return at })

	updated, _ := m.Update(moveResult{at: at})
	got := updated.(model)

	if got.logs[len(got.logs)-1] != "09:30:45  moved + restored" {
		t.Fatalf("latest log = %q", got.logs[len(got.logs)-1])
	}
}

func TestMoveResultRecordsFailure(t *testing.T) {
	at := time.Date(2026, time.August, 25, 9, 30, 45, 0, time.Local)
	m := newModel(time.Minute, func() error { return nil }, func() time.Time { return at })

	updated, _ := m.Update(moveResult{at: at, err: errors.New("event unavailable")})
	got := updated.(model)

	if !strings.Contains(got.logs[len(got.logs)-1], "movement failed: event unavailable") {
		t.Fatalf("latest log = %q", got.logs[len(got.logs)-1])
	}
}

func TestMoveTickUpdatesNextMovement(t *testing.T) {
	at := time.Date(2026, time.August, 25, 9, 30, 45, 0, time.Local)
	m := newModel(90*time.Second, func() error { return nil }, func() time.Time { return at })

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
	m := newModel(time.Minute, func() error { return nil }, time.Now)
	for range maxLogEntries + 3 {
		m.addLog("movement")
	}

	if len(m.logs) != maxLogEntries {
		t.Fatalf("log length = %d, want %d", len(m.logs), maxLogEntries)
	}
}
