# Project overview

## Objective

Ship a minimal, Homebrew-installable macOS CLI that periodically moves the pointer by one point and restores it without clicking or sending keystrokes.

## Scope

- macOS Apple Silicon only.
- Foreground `mouse-jiggler` process; Ctrl-C or SIGTERM stops it.
- Default interval: 60 seconds; configurable with `--interval`.
- Accessibility permission preflight with actionable guidance.
- Public GitHub releases from `wontaek/mouse-jiggler-go` and a Homebrew cask in `wontaek/homebrew-tap`.
- MIT license.

## Decisions

- Native CoreGraphics and ApplicationServices through a small cgo bridge; no large input-automation dependency.
- Each jiggle sends a one-point move and then restores the exact original location. No click or keyboard event is emitted.
- GitHub Actions releases Apple Silicon archives and updates the separate tap through GoReleaser.

## Current state

The CLI, tests, native bridge, GitHub Actions workflow, GoReleaser cask configuration, README, and MIT license are present. `go test ./...`, `go vet ./...`, and the local arm64 build pass. GitHub CLI is not authenticated, and GoReleaser is not installed locally, so repository creation, release, tap update, and live Accessibility verification remain outstanding.

## Next step

Authenticate GitHub CLI, create the public source repository and tap, add the tap token secret, push `main`, then test a tagged Homebrew release with Accessibility access.

See [plan/PLAN.md](plan/PLAN.md) for the executable plan.
