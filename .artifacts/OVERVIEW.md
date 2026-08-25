# Project overview

## Objective

Ship a minimal, Homebrew-installable macOS TUI that periodically moves the pointer by one point and restores it without clicking or sending keystrokes.

## Scope

- macOS Apple Silicon only.
- Foreground full-screen `mouse-jiggler` TUI. Without `--interval`, it requires a positive whole-number seconds selection before starting; with the flag, it starts immediately. It displays the interval, a live countdown, and eight recent movement results. `q`, Esc, Ctrl-C, or SIGTERM stops it.
- Default interval: 60 seconds; configurable interactively or with `--interval`.
- Accessibility permission preflight with actionable guidance.
- Public GitHub releases from `dnjdsxor21/mouse-jiggler-go` and a Homebrew formula in `dnjdsxor21/homebrew-tap`.
- MIT license.

## Decisions

- Native CoreGraphics and ApplicationServices through a small cgo bridge; no large input-automation dependency.
- Bubble Tea v2 owns setup, the display timer, and each scheduled movement. The sparse, adaptive-color TUI uses semantic states, readable status alignment, and bounded movement history.
- Each jiggle sends a one-point move and then restores the exact original location. No click or keyboard event is emitted.
- GitHub Actions releases Apple Silicon archives and updates the separate tap through GoReleaser.

## Current state

Version `v0.2.0` is released. Interval setup and visual hierarchy improvements are implemented and locally tested; they await a tagged release.

## Next step

Release the onboarding TUI as `v0.3.0`, upgrade the Homebrew formula, and verify the full interactive flow.

See [plan/PLAN.md](plan/PLAN.md) for the executable plan.
