package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dnjdsxor21/mouse-jiggler-go/internal/platform"
	"github.com/dnjdsxor21/mouse-jiggler-go/internal/tui"
)

var version = "dev"

type mouse interface {
	Trusted() bool
	Jiggle() error
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, platform.Mouse{}, tui.Run))
}

type launcher func(context.Context, time.Duration, bool, func() error) error

func run(args []string, stdout, stderr io.Writer, pointer mouse, start launcher) int {
	flags := flag.NewFlagSet("mouse-jiggler", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() {
		fmt.Fprintln(stderr, "Usage: mouse-jiggler [--interval DURATION]")
		fmt.Fprintln(stderr, "\nWithout --interval, choose the interval in the interactive TUI.")
		flags.PrintDefaults()
	}

	interval := flags.Duration("interval", time.Minute, "time between pointer movements")
	showVersion := flags.Bool("version", false, "print version and exit")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	intervalProvided := false
	flags.Visit(func(flag *flag.Flag) {
		intervalProvided = intervalProvided || flag.Name == "interval"
	})
	if *showVersion {
		fmt.Fprintln(stdout, version)
		return 0
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "mouse-jiggler: unexpected positional arguments")
		return 2
	}
	if *interval <= 0 {
		fmt.Fprintln(stderr, "mouse-jiggler: --interval must be positive")
		return 2
	}
	if !pointer.Trusted() {
		printAccessibilityInstructions(stderr)
		return 1
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := start(ctx, *interval, !intervalProvided, pointer.Jiggle); err != nil {
		fmt.Fprintf(stderr, "mouse-jiggler: %v\n", err)
		return 1
	}
	return 0
}

func printAccessibilityInstructions(stderr io.Writer) {
	executable, err := os.Executable()
	if err != nil {
		executable = "the mouse-jiggler executable"
	}

	fmt.Fprintln(stderr, "mouse-jiggler requires macOS Accessibility access.")
	fmt.Fprintln(stderr, "1. Open System Settings > Privacy & Security > Accessibility.")
	fmt.Fprintf(stderr, "2. Add and enable: %s\n", executable)
	fmt.Fprintln(stderr, "3. Run mouse-jiggler again.")
	fmt.Fprintln(stderr, "Open the settings pane with:")
	fmt.Fprintln(stderr, "  open 'x-apple.systempreferences:com.apple.preference.security?Privacy_Accessibility'")
}
