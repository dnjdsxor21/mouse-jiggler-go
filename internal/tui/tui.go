// Package tui renders the foreground mouse-jiggler status display.
package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
)

const (
	defaultPromptSeconds = 60
	maxLogEntries        = 8
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true)
	accentStyle  = lipgloss.NewStyle().Foreground(adaptiveColor("130", "214"))
	mutedStyle   = lipgloss.NewStyle().Foreground(adaptiveColor("242", "245"))
	successStyle = lipgloss.NewStyle().Foreground(adaptiveColor("28", "42"))
	dangerStyle  = lipgloss.NewStyle().Foreground(adaptiveColor("124", "203"))
	errorStyle   = dangerStyle.Bold(true)
)

type secondTick time.Time
type moveTick time.Time
type moveResult struct {
	at  time.Time
	err error
}
type quitMsg struct{}

type activity struct {
	at      time.Time
	success bool
	detail  string
}

type model struct {
	interval       time.Duration
	mover          func() error
	now            time.Time
	nextMove       time.Time
	logs           []activity
	askingInterval bool
	input          string
	inputError     string
	width          int
	height         int
	waitingFirst   bool
}

// Run starts the full-screen UI and returns when the user quits.
func Run(ctx context.Context, interval time.Duration, promptForInterval bool, mover func() error) error {
	program := tea.NewProgram(newModel(interval, promptForInterval, mover, time.Now))

	go func() {
		<-ctx.Done()
		program.Send(quitMsg{})
	}()

	_, err := program.Run()
	return err
}

func newModel(interval time.Duration, promptForInterval bool, mover func() error, now func() time.Time) model {
	current := now()
	return model{
		interval:       interval,
		mover:          mover,
		now:            current,
		nextMove:       current.Add(interval),
		logs:           make([]activity, 0, maxLogEntries),
		askingInterval: promptForInterval,
		input:          strconv.Itoa(defaultPromptSeconds),
	}
}

func (m model) Init() tea.Cmd {
	if m.askingInterval {
		return nil
	}
	return m.startRunning()
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	if m.askingInterval {
		return m.updateIntervalPrompt(message)
	}

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		m.height = message.Height
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
		m.waitingFirst = false
		if message.err != nil {
			m.addLog(activity{at: message.at, detail: fmt.Sprintf("movement failed: %v", message.err)})
		} else {
			m.addLog(activity{at: message.at, success: true, detail: "moved + restored"})
		}
	}

	return m, nil
}

func (m model) updateIntervalPrompt(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		m.height = message.Height
	case quitMsg:
		return m, tea.Quit
	case tea.KeyPressMsg:
		key := message.String()
		switch key {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			interval, err := parseSeconds(m.input)
			if err != nil {
				m.inputError = err.Error()
				return m, nil
			}
			m.interval = interval
			m.askingInterval = false
			m.inputError = ""
			command := m.startRunning()
			return m, command
		case "backspace", "ctrl+h":
			if len(m.input) > 0 {
				_, size := utf8.DecodeLastRuneInString(m.input)
				m.input = m.input[:len(m.input)-size]
			}
			m.inputError = ""
		default:
			if len(key) == 1 && key[0] >= '0' && key[0] <= '9' {
				m.input += key
				m.inputError = ""
			}
		}
	}

	return m, nil
}

func (m *model) startRunning() tea.Cmd {
	m.now = time.Now()
	m.nextMove = m.now.Add(m.interval)
	m.logs = m.logs[:0]
	m.waitingFirst = true
	return tea.Batch(secondTickCmd(), moveCmd(m.mover), moveTickCmd(m.interval))
}

func (m model) View() tea.View {
	if m.askingInterval {
		return m.intervalPromptView()
	}

	var output strings.Builder
	m.writeLine(&output, titleStyle.Render("mouse-jiggler")+"  "+accentStyle.Render("[running]"))
	output.WriteString("\n")
	m.writeLine(&output, fmt.Sprintf("%-14s%s", "interval", m.interval))
	m.writeLine(&output, fmt.Sprintf("%-14s%s", "next move", formatRemaining(m.nextMove.Sub(m.now))))
	output.WriteString("\n")
	m.writeLine(&output, mutedStyle.Render("recent activity"))
	if m.waitingFirst {
		m.writeLine(&output, mutedStyle.Render("--:--:--  [wait]   first movement in progress"))
	} else {
		for _, entry := range m.visibleLogs() {
			m.writeLine(&output, m.renderActivity(entry))
		}
	}
	output.WriteString("\n")
	m.writeLine(&output, mutedStyle.Render(m.quitHelp()))

	view := tea.NewView(output.String())
	view.AltScreen = true
	return view
}

func (m model) intervalPromptView() tea.View {
	var output strings.Builder
	m.writeLine(&output, titleStyle.Render("mouse-jiggler")+"  "+accentStyle.Render("[setup]"))
	output.WriteString("\n")
	m.writeLine(&output, "Movement interval")
	m.writeLine(&output, mutedStyle.Render("How often should the pointer move one point and return?"))
	output.WriteString("\n")
	m.writeLine(&output, mutedStyle.Render("Seconds"))
	m.writeLine(&output, accentStyle.Render("> ")+m.input+accentStyle.Render("▏"))
	if m.inputError != "" {
		m.writeLine(&output, "  "+errorStyle.Render(m.inputError))
	}
	output.WriteString("\n")
	m.writeLine(&output, mutedStyle.Render(m.quitHelpWithStart()))

	view := tea.NewView(output.String())
	view.AltScreen = true
	return view
}

func (m model) writeLine(output *strings.Builder, text string) {
	output.WriteString(m.gutter())
	output.WriteString(text)
	output.WriteByte('\n')
}

func (m model) gutter() string {
	if m.width > 0 && m.width < 40 {
		return " "
	}
	return "  "
}

func (m model) quitHelp() string {
	if m.width > 0 && m.width < 40 {
		return "q/esc/ctrl-c: quit"
	}
	return "q / esc / ctrl-c: quit"
}

func (m model) quitHelpWithStart() string {
	if m.width > 0 && m.width < 40 {
		return "enter: start  q/esc/ctrl-c: quit"
	}
	return "enter: start    q / esc / ctrl-c: quit"
}

func (m model) visibleLogs() []activity {
	availableRows := maxLogEntries
	if m.height > 0 {
		availableRows = max(1, m.height-9)
	}
	if len(m.logs) <= availableRows {
		return m.logs
	}
	return m.logs[:availableRows]
}

func (m model) renderActivity(entry activity) string {
	prefix := entry.at.Format("15:04:05") + "  "
	status := "[error]"
	statusStyle := dangerStyle
	if entry.success {
		status = "[ok]"
		statusStyle = successStyle
	}
	available := m.contentWidth() - lipgloss.Width(prefix) - len(status) - 4
	detail := truncate(entry.detail, available)
	return prefix + statusStyle.Render(status) + "     " + detail
}

func (m model) contentWidth() int {
	if m.width <= 0 {
		return 76
	}
	return max(1, m.width-lipgloss.Width(m.gutter())-2)
}

func (m *model) addLog(entry activity) {
	if len(m.logs) < maxLogEntries {
		m.logs = append(m.logs, activity{})
	}
	copy(m.logs[1:], m.logs[:len(m.logs)-1])
	m.logs[0] = entry
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

func parseSeconds(input string) (time.Duration, error) {
	if input == "" {
		return 0, fmt.Errorf("Enter an interval in seconds.")
	}

	seconds, err := strconv.ParseInt(input, 10, 64)
	if err != nil || seconds <= 0 {
		return 0, fmt.Errorf("Enter a positive whole number of seconds.")
	}
	if seconds > int64(time.Duration(1<<63-1)/time.Second) {
		return 0, fmt.Errorf("That interval is too large.")
	}
	return time.Duration(seconds) * time.Second, nil
}

func formatRemaining(remaining time.Duration) string {
	if remaining <= 0 {
		return "00:00"
	}
	seconds := int64((remaining + time.Second - 1) / time.Second)
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds%60)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds%60)
}

func truncate(text string, width int) string {
	if width <= 0 {
		return ""
	}
	if lipgloss.Width(text) <= width {
		return text
	}
	if width <= 3 {
		return strings.Repeat(".", width)
	}

	var output strings.Builder
	for _, runeValue := range text {
		candidate := output.String() + string(runeValue)
		if lipgloss.Width(candidate) > width-3 {
			break
		}
		output.WriteRune(runeValue)
	}
	return output.String() + "..."
}

func adaptiveColor(light, dark string) compat.AdaptiveColor {
	return compat.AdaptiveColor{
		Light: lipgloss.Color(light),
		Dark:  lipgloss.Color(dark),
	}
}
