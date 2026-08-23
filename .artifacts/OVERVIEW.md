# Project overview

## Objective

Ship a minimal, Homebrew-installable macOS CLI that periodically moves the pointer by one point and restores it without clicking or sending keystrokes.

## Scope

- macOS Apple Silicon only.
- Foreground `mouse-jiggler` process; Ctrl-C or SIGTERM stops it.
- Default interval: 60 seconds; configurable with `--interval`.
- Accessibility permission preflight with actionable guidance.
- Public GitHub releases from `dnjdsxor21/mouse-jiggler-go` and a Homebrew cask in `dnjdsxor21/homebrew-tap`.
- MIT license.

## Decisions

- Native CoreGraphics and ApplicationServices through a small cgo bridge; no large input-automation dependency.
- Each jiggle sends a one-point move and then restores the exact original location. No click or keyboard event is emitted.
- GitHub Actions releases Apple Silicon archives and updates the separate tap through GoReleaser.

## Current state

The CLI, tests, native bridge, GitHub Actions workflow, GoReleaser cask configuration, README, and MIT license are present. `go test ./...`, `go vet ./...`, and the local arm64 build pass. GitHub CLI is authenticated as `dnjdsxor21`; the approved publication owner is now `dnjdsxor21`.

## Next step

Create the public source repository and tap, add the release secret, push `main`, and test a tagged Homebrew release with Accessibility access.

See [plan/PLAN.md](plan/PLAN.md) for the executable plan.
