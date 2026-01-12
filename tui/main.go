package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"tact-tui/api"
	"tact-tui/ui"
)

const defaultAPIURL = "http://localhost:2100"

func main() {
	apiURL := flag.String("api", "", "Backend API URL (default: http://localhost:2100)")
	flag.Parse()

	// Priority: flag > env var > default
	url := *apiURL
	if url == "" {
		url = os.Getenv("TACT_API_URL")
	}
	if url == "" {
		url = defaultAPIURL
	}

	client := api.NewClient(url)
	app := ui.NewApp(client)

	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
