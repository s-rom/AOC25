package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var countStyle lipgloss.Style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#BB8ED0")).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#3291B6")).
	Align(lipgloss.Left)

var splitterStyle lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#D1D3D4"))

var activatedSplitterStyle lipgloss.Style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#3291B6"))

var startStyle lipgloss.Style = lipgloss.NewStyle().
	Bold(true).
	Background(lipgloss.Color("#BB8ED0"))

var rayStyle lipgloss.Style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#BB8ED0"))

var emptyStyle lipgloss.Style = lipgloss.NewStyle().
	Bold(true)

type model struct {
	Grid        *Grid[string]
	Queue       []Position
	TotalSplits uint64
	CurrentRow  int
}

type tickMsg time.Time

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.CurrentRow < m.Grid.Rows {
			m.CurrentRow++
			return m, tickCmd()
		} else {
			return m, nil
		}
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	default:
		return m, nil
	}

	return m, nil
}

func (m model) View() string {

	var strBuilder strings.Builder
	strBuilder.WriteString(countStyle.Render(fmt.Sprintf("Splits: %d", m.TotalSplits)))
	strBuilder.WriteString("\n")

	start := max(0, m.CurrentRow-20)
	end := min(start+30, m.Grid.Rows)

	for r := start; r < end; r++ {
		for c := range m.Grid.Columns {
			ptr := m.Grid.GetPtrAt(Position{Row: r, Column: c})

			switch *ptr {
			case "S":
				strBuilder.WriteString(startStyle.Render("*"))
			case ".":
				strBuilder.WriteString(emptyStyle.Render(" "))
			case "|":
				if r <= m.CurrentRow {
					strBuilder.WriteString(rayStyle.Render("|"))
				} else {
					strBuilder.WriteString(emptyStyle.Render(" "))
				}
			case "^":
				if r <= m.CurrentRow {
					strBuilder.WriteString(activatedSplitterStyle.Render(*ptr))
				} else {
					strBuilder.WriteString(splitterStyle.Render(*ptr))
				}
			}
		}
		strBuilder.WriteRune('\n')
	}

	return strBuilder.String()

}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/15.0, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func StartViz(grid *Grid[string]) {
	m := model{
		Grid:        grid,
		Queue:       make([]Position, 0, 200),
		CurrentRow:  0,
		TotalSplits: 0,
	}
	m.Queue = append(m.Queue, MoveDown(FindStartPosition(grid)))

	var splits, total uint64
	splits, m.Queue = PropagateBreadth(grid, m.Queue)
	for len(m.Queue) != 0 {
		total += splits
		splits, m.Queue = PropagateBreadth(grid, m.Queue)
	}

	m.TotalSplits = total

	if _, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
