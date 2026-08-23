// Package jiggler schedules pointer movement at a fixed interval.
package jiggler

import (
	"context"
	"errors"
	"time"
)

// Run invokes jiggle immediately and then once per interval until ctx is done.
func Run(ctx context.Context, interval time.Duration, jiggle func() error) error {
	if interval <= 0 {
		return errors.New("interval must be positive")
	}

	if err := jiggle(); err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := jiggle(); err != nil {
				return err
			}
		}
	}
}
