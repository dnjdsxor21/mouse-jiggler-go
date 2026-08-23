package jiggler

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestRunRejectsNonPositiveInterval(t *testing.T) {
	err := Run(context.Background(), 0, func() error { return nil })
	if err == nil || !strings.Contains(err.Error(), "positive") {
		t.Fatalf("Run() error = %v, want positive-interval error", err)
	}
}

func TestRunJigglesImmediately(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	calls := 0
	err := Run(ctx, time.Hour, func() error {
		calls++
		cancel()
		return nil
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if calls != 1 {
		t.Fatalf("jiggle calls = %d, want 1", calls)
	}
}
