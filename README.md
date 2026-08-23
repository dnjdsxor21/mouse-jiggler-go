# mouse-jiggler

A minimal macOS command-line mouse jiggler. It moves the pointer by one point and restores it at a fixed interval. It never clicks or sends keyboard input.

## Requirements

- macOS on Apple Silicon
- Accessibility permission for the installed executable

## Install

```sh
brew tap wontaek/tap
brew install --cask mouse-jiggler
```

## Use

```sh
# Move once immediately, then every 60 seconds.
mouse-jiggler

# Use a different positive Go duration.
mouse-jiggler --interval 2m30s

# Stop with Ctrl-C.
mouse-jiggler --interval 10s
```

`mouse-jiggler --version` prints the installed version. `mouse-jiggler --help` shows the command options.

## Grant Accessibility access

Before the first run, open **System Settings → Privacy & Security → Accessibility**, add the executable path printed by `mouse-jiggler`, and enable it. Then run the command again.

Open the settings pane directly:

```sh
open 'x-apple.systempreferences:com.apple.preference.security?Privacy_Accessibility'
```

## Build locally

```sh
go test ./...
go build -o mouse-jiggler ./cmd/mouse-jiggler
./mouse-jiggler --interval 10s
```

## Release

1. Create public repositories `wontaek/mouse-jiggler-go` and `wontaek/homebrew-tap`.
2. Add the repository secret `RELEASE_GITHUB_TOKEN`: a fine-grained GitHub token with Contents read/write access limited to both `wontaek/mouse-jiggler-go` and `wontaek/homebrew-tap`.
3. Push `main`, then push a tag such as `v0.1.0`.
4. GitHub Actions tests the source, creates the GitHub release, and updates the Homebrew cask in the tap.

The tag workflow uses GoReleaser's supported `homebrew_casks` publisher rather than its deprecated formula publisher.

## Scope

macOS Apple Silicon only. No GUI, daemon mode, clicks, keyboard events, Windows/Linux support, Intel artifact, or code signing/notarization.

## License

MIT. See [LICENSE](LICENSE).
