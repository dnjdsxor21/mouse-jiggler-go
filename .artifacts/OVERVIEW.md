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

Version `v0.1.1` is released. A Bubble Tea TUI with countdown and bounded movement log is implemented and locally verified; it awaits a tagged release.

## Next step

Release the TUI as `v0.2.0`, install it through Homebrew, and repeat the terminal smoke test.

See [plan/PLAN.md](plan/PLAN.md) for the executable plan.
