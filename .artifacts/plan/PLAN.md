# Mouse jiggler release plan

## Product contract

```text
mouse-jiggler [--interval DURATION]
mouse-jiggler --version
```

- `--interval` uses Go duration syntax and defaults to `60s`; values must be positive.
- The foreground full-screen TUI shows the configured interval, a live countdown, and the eight most recent movement results; `q`, Esc, Ctrl-C, or SIGTERM stops it.
- A jiggle reads the current pointer position, posts a one-point mouse-move event, waits briefly so the event is observable, and posts a move back to the original position.
- The program never posts a click or keyboard event.
- Before moving, it verifies macOS Accessibility trust. If missing, it exits non-zero with the executable path, the System Settings path, and the command that opens Accessibility settings.

## Implementation

1. Initialize `github.com/dnjdsxor21/mouse-jiggler-go` as a Go 1.26 module with MIT licensing and a narrow Go `.gitignore`.
2. Add a Bubble Tea v2 full-screen model that ticks every second for the countdown, schedules the initial and periodic jiggler commands, and retains eight movement results.
3. Add the CLI entry point. It parses flags, validates duration, handles `--help` and `--version`, checks Accessibility trust, and uses a signal-aware context.
4. Add a darwin/cgo implementation linked only to `ApplicationServices`. It calls `AXIsProcessTrusted`, reads the location through `CGEventCreate`, and posts the outbound/restoration mouse events. No third-party input library is used.
5. Add a concise README with local build/run, Accessibility setup, Homebrew install, CLI examples, scope limits, and release instructions.

## Release pipeline

1. GitHub Actions runs `go test ./...` and `go vet ./...` for pull requests and pushes.
2. A `v*` tag runs GoReleaser. It builds a `darwin/arm64` archive, injects the tag version, produces checksums, creates the GitHub release, and commits the generated Homebrew formula to `dnjdsxor21/homebrew-tap`.
3. The release workflow requires `RELEASE_GITHUB_TOKEN`: a fine-grained GitHub token with Contents read/write permission limited to both repositories. GoReleaser uses it for the source GitHub release and the separate tap formula commit.
4. The tap formula installs the archived `mouse-jiggler` binary and verifies `mouse-jiggler --version`.
5. A user installs it with `brew tap dnjdsxor21/tap` and `brew install mouse-jiggler`.

## Verification

- Unit test the runner and CLI parsing/validation.
- Build and run the binary on Apple Silicon.
- With Accessibility permission enabled, run at a short interval and interrupt it; verify pointer motion, restoration, and clean exit.
- Tag a prerelease, confirm its GitHub archive/checksum/formula commit, install through Homebrew, and run `--version` plus a short live session.

## Exclusions

No Windows/Linux support, GUI or menu-bar app, daemon/start-stop protocol, clicks, keyboard input, Intel artifact, code signing/notarization, npm/uvx packages, or runtime configuration files.
