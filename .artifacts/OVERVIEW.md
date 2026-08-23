# Project overview

## Objective

Ship a minimal, Homebrew-installable macOS CLI that periodically moves the pointer by one point and restores it without clicking or sending keystrokes.

## Scope

- macOS Apple Silicon only.
- Foreground `mouse-jiggler` process; Ctrl-C or SIGTERM stops it.
- Default interval: 60 seconds; configurable with `--interval`.
- Accessibility permission preflight with actionable guidance.
- Public GitHub releases from `dnjdsxor21/mouse-jiggler-go` and a Homebrew formula in `dnjdsxor21/homebrew-tap`.
- MIT license.

## Decisions

- Native CoreGraphics and ApplicationServices through a small cgo bridge; no large input-automation dependency.
- Each jiggle sends a one-point move and then restores the exact original location. No click or keyboard event is emitted.
- GitHub Actions releases Apple Silicon archives and updates the separate tap through GoReleaser.

## Current state

Version `v0.1.1` is released. Its GitHub Actions verification and release jobs passed; the formula is published at the root of `dnjdsxor21/homebrew-tap`; the obsolete `v0.1.0` cask was deleted. `brew install dnjdsxor21/tap/mouse-jiggler` and the formula's `brew test` both passed on Apple Silicon.

## Next step

Users grant Accessibility access for the installed executable and run `mouse-jiggler`; a maintainer rotates `RELEASE_GITHUB_TOKEN` before its expiry.

See [plan/PLAN.md](plan/PLAN.md) for the executable plan.
