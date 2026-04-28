package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

// Session storage
type Session struct {
	Type  string     `json:"type"`
	Start time.Time  `json:"start"`
	End   *time.Time `json:"end"`
}

func storePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".hocusfocus.json")
}

func loadSessions() []Session {
	data, err := os.ReadFile(storePath())
	if err != nil {
		return []Session{}
	}

	var sessions []Session
	_ = json.Unmarshal(data, &sessions)
	return sessions
}

func saveSessions(sessions []Session) {
	data, _ := json.MarshalIndent(sessions, "", "  ")
	_ = os.WriteFile(storePath(), data, 0o644)
}

func currentSession(sessions []Session) *Session {
	for i := range sessions {
		if sessions[i].End == nil {
			return &sessions[i]
		}
	}
	return nil
}

// CLI commands
func printStats() {
	sessions := loadSessions()
	totals := map[string]time.Duration{}

	for _, s := range sessions {
		if s.End == nil {
			continue
		}
		totals[s.Type] += s.End.Sub(s.Start)
	}

	if len(totals) == 0 {
		fmt.Println(" No sessions have been completed.")
		return
	}

	for typ, d := range totals {
		fmt.Printf("%s: %s\n", typ, d.Round(time.Second))
	}
}

func printCurrentSession() {
	sessions := loadSessions()
	s := currentSession(sessions)

	if s == nil {
		fmt.Println(" No current session")
		return
	}

	fmt.Printf(
		" Current session: %s (%s)\n",
		s.Type,
		time.Since(s.Start).Round(time.Second),
	)
}

func printHelp() {
	fmt.Println("hocusfocus help:")
	fmt.Println("<no args>      : choose/stop session")
	fmt.Println("help           : print this message")
	fmt.Println("currentsession : print current session")
	fmt.Println("stats          : print statistics")
}

// Bubbletea TUI
type item string

func (i item) FilterValue() string { return "" }

type stopItem struct{}

func (s stopItem) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var str string
	switch i := listItem.(type) {
	case item:
		str = fmt.Sprintf("%d. %s", index+1, string(i))
	case stopItem:
		str = fmt.Sprintf("%d. %s", index+1, "Stop current session")
	default:
		return
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	sessions []Session
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			i := m.list.SelectedItem()
			if i == nil {
				return m, nil
			}

			now := time.Now()

			switch v := i.(type) {
			case stopItem:
				// Stop session
				for idx := range m.sessions {
					if m.sessions[idx].End == nil {
						m.sessions[idx].End = &now
						break
					}
				}
				saveSessions(m.sessions)
				return m, tea.Quit

			case item:
				// End active session if one exists
				for idx := range m.sessions {
					if m.sessions[idx].End == nil {
						m.sessions[idx].End = &now
						break
					}
				}
				// Start new session
				m.sessions = append(m.sessions, Session{
					Type:  string(v),
					Start: now,
					End:   nil,
				})
				saveSessions(m.sessions)
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return "\n" + m.list.View()
}

// Build list
func buildItems(sessions []Session) []list.Item {
	items := []list.Item{
		item("Work"),
		item("Study"),
		item("Waste"),
	}

	if currentSession(sessions) != nil {
		// Prepend Stop item if a session is active
		items = append([]list.Item{stopItem{}}, items...)
	}

	return items
}

// Main function
func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "stats":
			printStats()
			return
		case "currentsession":
			printCurrentSession()
			return
		case "help":
			printHelp()
			return
		}
	}

	sessions := loadSessions()
	items := buildItems(sessions)

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select a Session"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{
		list:     l,
		sessions: sessions,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
