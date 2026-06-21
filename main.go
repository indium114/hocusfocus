package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/huh"
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

func stopSession(sessions []Session) {
	now := time.Now()

	for idx := range sessions {
		if sessions[idx].End == nil {
			sessions[idx].End = &now
			break
		}
	}

	saveSessions(sessions)
}

func startSession(kind string, sessions []Session) {
	now := time.Now()

	// End active session if one exists
	for idx := range sessions {
		if sessions[idx].End == nil {
			sessions[idx].End = &now
			break
		}
	}
	// Start new session
	sessions = append(sessions, Session{
		Type:  string(kind),
		Start: now,
		End:   nil,
	})
	saveSessions(sessions)
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

	var kind string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("HocusFocus").
				Options(
					huh.NewOption("Work", "work"),
					huh.NewOption("Study", "study"),
					huh.NewOption("Waste", "waste"),
					huh.NewOption("Stop Current Session", "stop"),
				).
				Value(&kind),
		),
	)

	err := form.Run()
	if err != nil {
		panic(err)
	}

	if kind == "stop" {
		if currentSession(sessions) == nil {
			fmt.Println("No current session")
		} else {
			stopSession(sessions)
		}
	} else {
		startSession(kind, sessions)
	}
}
