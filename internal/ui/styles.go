package ui

import "github.com/charmbracelet/lipgloss"

// Dark terminal palette — cyan/teal primary, red accent
var (
	cyan    = lipgloss.Color("#5AC8AF")
	cyanDim = lipgloss.Color("#3A8A78")
	red     = lipgloss.Color("#E06C75")
	yellow  = lipgloss.Color("#D19A66")
	green   = lipgloss.Color("#98C379")
	white   = lipgloss.Color("#ABB2BF")
	dim     = lipgloss.Color("#5C6370")
	faint   = lipgloss.Color("#3E4451")
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	accentStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true)

	countStyle = lipgloss.NewStyle().
			Foreground(dim)

	cursorStyle = lipgloss.NewStyle().
			Foreground(red)

	selectedStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(white)

	installedBadge = lipgloss.NewStyle().
			Foreground(green)

	dimStyle = lipgloss.NewStyle().
			Foreground(dim)

	separatorStyle = lipgloss.NewStyle().
			Foreground(faint)

	errorStyle = lipgloss.NewStyle().
			Foreground(red)

	successStyle = lipgloss.NewStyle().
			Foreground(green)

	warningStyle = lipgloss.NewStyle().
			Foreground(yellow)

	progressFilled = lipgloss.NewStyle().
			Foreground(cyan)

	progressEmpty = lipgloss.NewStyle().
			Foreground(faint)

	pctStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	helpKey = lipgloss.NewStyle().
		Foreground(cyan)

	helpDesc = lipgloss.NewStyle().
			Foreground(dim)

	searchStyle = lipgloss.NewStyle().
			Foreground(cyan)

	searchInputStyle = lipgloss.NewStyle().
				Foreground(white)

	sectionLabel = lipgloss.NewStyle().
			Foreground(dim).
			Bold(true)
)
