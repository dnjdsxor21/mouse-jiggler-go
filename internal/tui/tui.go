// Package tui renders the foreground mouse-jiggler status display.
package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

const maxLogEntries = 8

type secondTick time.Time
type moveTick time.Time

type moveResult struct {
	at  time.Time
	err error
}
type quitMsg struct{}

type model struct {
	interval time.Duration
	mover    func() error
	now      time.Time
	nextMove time.Time
	logs     []string
}

// Run starts the full-screen status UI and returns when the user quits.
func Run(ctx context.Context, interval time.Duration, mover func() error) error {
	program := tea.NewProgram(newModel(interval, mover, time.Now))

	go func() {
		<-ctx.Done()
		program.Send(quitMsg{})
	}()

	_, err := program.Run()
	return err
}

func newModel(interval time.Duration, mover func() error, now func() time.Time) model {
	current := now()
	return model{
		interval: interval,
		mover:    mover,
		now:      current,
		nextMove: current.Add(interval),
		logs:     append(make([]string, 0, maxLogEntries), "waiting for first movement"),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(secondTickCmd(), moveCmd(m.mover), moveTickCmd(m.interval))
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyPressMsg:
		switch message.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case quitMsg:
		return m, tea.Quit
	case secondTick:
		m.now = time.Time(message)
		return m, secondTickCmd()
	case moveTick:
		m.nextMove = time.Time(message).Add(m.interval)
		return m, tea.Batch(moveCmd(m.mover), moveTickCmd(m.interval))
	case moveResult:
		if message.err != nil {
			m.addLog(fmt.Sprintf("%s  movement failed: %v", message.at.Format("15:04:05"), message.err))
		} else {
			m.addLog(fmt.Sprintf("%s  moved + restored", message.at.Format("15:04:05")))
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	remaining := m.nextMove.Sub(m.now)
	if remaining < 0 {
		remaining = 0
	}

	var output strings.Builder
	fmt.Fprintln(&output, "mouse-jiggler")
	fmt.Fprintf(&output, "interval: %s    next movement: %s\n\n", m.interval, formatRemaining(remaining))
	fmt.Fprintln(&output, "recent activity")
	for index := len(m.logs) - 1; index >= 0; index-- {
		fmt.Fprintf(&output, "  %s\n", m.logs[index])
	}
	fmt.Fprintln(&output, "\nq / esc / ctrl-c: quit")

	view := tea.NewView(output.String())
	view.AltScreen = true
	return view
}

func (m *model) addLog(entry string) {
	if len(m.logs) == maxLogEntries {
		copy(m.logs, m.logs[1:])
		m.logs = m.logs[:maxLogEntries-1]
	}
	m.logs = append(m.logs, entry)
}

func secondTickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(now time.Time) tea.Msg {
		return secondTick(now)
	})
}

func moveTickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(now time.Time) tea.Msg {
		return moveTick(now)
	})
}

func moveCmd(mover func() error) tea.Cmd {
	return func() tea.Msg {
		return moveResult{at: time.Now(), err: mover()}
	}
}

func formatRemaining(remaining time.Duration) string {
	return remaining.Round(time.Second).String()
}
