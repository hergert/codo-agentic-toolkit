package tui

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type InitResult struct {
	Stacks    []string // exact keys to feed into pack.Files()
	Confirmed bool
}

type option struct {
	Key   string // exact key for pack.Files()
	Label string // pretty label in the TUI
}

var opts = []option{
	{Key: "go", Label: "Go"},
	{Key: "typescript", Label: "TypeScript / Node"},
	{Key: "python", Label: "Python"},
	{Key: "flutter", Label: "Flutter"},
}

type model struct {
	cursor   int
	selected map[int]bool
	quit     bool
	confirm  bool
}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*10, func(time.Time) tea.Msg { return nil })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch key := msg.(type) {
	case tea.KeyMsg:
		switch key.String() {
		case "ctrl+c", "q":
			m.quit, m.confirm = true, false
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(opts)-1 {
				m.cursor++
			}
		case " ":
			if m.selected == nil {
				m.selected = map[int]bool{}
			}
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "a": // select all
			if m.selected == nil {
				m.selected = map[int]bool{}
			}
			for i := range opts {
				m.selected[i] = true
			}
		case "n": // select none
			m.selected = map[int]bool{}
		case "enter":
			m.confirm, m.quit = true, true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder
	fmt.Fprint(&b, "Select stacks to include (space to toggle):\n\n")
	for i, o := range opts {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		box := " "
		if m.selected[i] {
			box = "x"
		}
		fmt.Fprintf(&b, " %s [%s] %s\n", cursor, box, o.Label)
	}
	fmt.Fprintln(&b, "\n[↑/↓/j/k] move   [space] toggle   [a] all   [n] none   [enter] continue   [q] abort")
	return b.String()
}

func RunInitWizard(_ context.Context) (InitResult, error) {
	m := model{selected: map[int]bool{}}
	pm := tea.NewProgram(m)
	res, err := pm.Run()
	if err != nil {
		return InitResult{}, err
	}

	out := res.(model)
	if out.quit && !out.confirm {
		return InitResult{Stacks: nil, Confirmed: false}, nil
	}
	// Build stable list of selected keys
	var keys []string
	for i, on := range out.selected {
		if on {
			keys = append(keys, opts[i].Key)
		}
	}
	sort.Strings(keys)
	return InitResult{Stacks: keys, Confirmed: true}, nil
}
