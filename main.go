package main

import (
	"fmt"
	"os"

	"nfdownloader/internal/fonts"
	"nfdownloader/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	platform := fonts.DetectPlatform()

	if platform.OS != "linux" {
		fmt.Println("Error: NFDownloader currently only supports Linux (Ubuntu/Debian-based distros). :-(")
		os.Exit(1)
	}

	m := ui.NewModel(platform)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
