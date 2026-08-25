# Project overview

## Objective

Ship a minimal, Homebrew-installable macOS TUI that periodically moves the pointer by one point and restores it without clicking or sending keystrokes.

## Scope

- macOS Apple Silicon only.
- Foreground full-screen `mouse-jiggler` TUI; it displays the interval, a live countdown, and eight recent movement results. `q`, Esc, Ctrl-C, or SIGTERM stops it.
- Default interval: 60 seconds; configurable with `--interval`.
- Accessibility permission preflight with actionable guidance.
- Public GitHub releases from `dnjdsxor21/mouse-jiggler-go` and a Homebrew formula in `dnjdsxor21/homebrew-tap`.
- MIT license.

## Decisions

- Native CoreGraphics and ApplicationServices through a small cgo bridge; no large input-automation dependency.
- Bubble Tea v2 owns the display timer and schedules each movement, so movement success or failure is visible in the terminal.
- Each jiggle sends a one-point move and then restores the exact original location. No click or keyboard event is emitted.
- GitHub Actions releases Apple Silicon archives and updates the separate tap through GoReleaser.

## Current state

Version `v0.2.0` is released. GitHub Actions verification and release jobs passed. The Homebrew formula upgrade from `0.1.1` to `0.2.0` and its formula test passed on Apple Silicon. The installed TUI was run at a two-second interval; it rendered its countdown, recorded initial and periodic `moved + restored` entries, and exited with `q` and status 0.

## Next step

Users install or upgrade with Homebrew, grant Accessibility access if needed, then run `mouse-jiggler`.

See [plan/PLAN.md](plan/PLAN.md) for the executable plan.
