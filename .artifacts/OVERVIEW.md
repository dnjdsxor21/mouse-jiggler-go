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

Version `v0.1.0` was released, and its GitHub Actions verification and release jobs passed. The generated cask installed but macOS Gatekeeper rejected its unsigned archive binary. The release pipeline is being migrated to a Homebrew formula, which installs the same CLI outside cask quarantine; `v0.1.1` will replace the obsolete cask with the formula.

## Next step

Release `v0.1.1` with the Homebrew formula, remove the obsolete cask, install the formula, and test the live Accessibility path.

See [plan/PLAN.md](plan/PLAN.md) for the executable plan.
