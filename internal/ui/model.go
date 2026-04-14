package ui

import (
	"fmt"
	"strings"

	"nfdownloader/internal/fonts"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	stateLoading state = iota
	stateList
	stateDownloading
	stateConfirmUninstall
	stateSearch
)

// Messages
type fontsLoadedMsg struct {
	fonts []fonts.Font
	err   error
}

type downloadProgressMsg struct {
	downloaded int64
	total      int64
}

type downloadDoneMsg struct{ err error }
type uninstallDoneMsg struct{ err error }

// displayItem represents a row in the font list (font or separator).
type displayItem struct {
	font      *fonts.Font
	separator string // non-empty means this is a section header
}

type Model struct {
	platform  fonts.Platform
	state     state
	spinner   spinner.Model
	fonts     []fonts.Font
	installed map[string]bool
	cursor     int
	offset     int
	height     int
	width      int
	termHeight int
	errMsg    string
	statusMsg string

	// display list (installed first, then separator, then available)
	display []displayItem

	// search
	searchQuery  string
	searchActive bool

	// download state
	downloading  string
	dlDownloaded int64
	dlTotal      int64
	progressCh   chan downloadProgressMsg
}

func NewModel(platform fonts.Platform) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return Model{
		platform:  platform,
		state:     stateLoading,
		spinner:   s,
		installed: make(map[string]bool),
		height:    20,
		width:     80,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchFontsCmd)
}

func fetchFontsCmd() tea.Msg {
	fl, err := fonts.FetchFontList()
	return fontsLoadedMsg{fonts: fl, err: err}
}

func waitForProgress(ch chan downloadProgressMsg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return nil
		}
		return msg
	}
}

// buildDisplay creates the display list: installed fonts first, separator, then available.
// If a search query is active, it filters by substring match.
func (m *Model) buildDisplay() {
	m.display = nil
	query := strings.ToLower(m.searchQuery)

	var installedFonts, availableFonts []fonts.Font
	for i := range m.fonts {
		f := &m.fonts[i]
		if query != "" && !strings.Contains(strings.ToLower(f.Name), query) {
			continue
		}
		if m.installed[f.Name] {
			installedFonts = append(installedFonts, *f)
		} else {
			availableFonts = append(availableFonts, *f)
		}
	}

	if len(installedFonts) > 0 {
		m.display = append(m.display, displayItem{separator: "installed"})
		for i := range installedFonts {
			m.display = append(m.display, displayItem{font: &installedFonts[i]})
		}
	}

	if len(availableFonts) > 0 {
		m.display = append(m.display, displayItem{separator: "available"})
		for i := range availableFonts {
			m.display = append(m.display, displayItem{font: &availableFonts[i]})
		}
	}

	// Ensure cursor is on a font row, not a separator
	m.clampCursor()
}

func (m *Model) clampCursor() {
	if len(m.display) == 0 {
		m.cursor = 0
		m.offset = 0
		return
	}
	if m.cursor >= len(m.display) {
		m.cursor = len(m.display) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	// Skip separator
	if m.display[m.cursor].separator != "" {
		m.moveCursorDown()
	}
}

func (m *Model) moveCursorDown() {
	for m.cursor < len(m.display)-1 {
		m.cursor++
		if m.display[m.cursor].font != nil {
			break
		}
	}
	// If we landed on separator at end, try going up
	if m.display[m.cursor].separator != "" {
		for m.cursor > 0 {
			m.cursor--
			if m.display[m.cursor].font != nil {
				break
			}
		}
	}
	m.fixOffset()
}

func (m *Model) moveCursorUp() {
	for m.cursor > 0 {
		m.cursor--
		if m.display[m.cursor].font != nil {
			break
		}
	}
	// If we landed on separator at top, go down
	if m.display[m.cursor].separator != "" {
		m.moveCursorDown()
	}
	m.fixOffset()
}

func (m *Model) fixOffset() {
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+m.height {
		m.offset = m.cursor - m.height + 1
	}
	// Don't start on a separator if possible
	if m.offset > 0 && m.offset < len(m.display) && m.display[m.offset].separator != "" {
		if m.offset > 0 {
			m.offset--
		}
	}
}

func (m *Model) fontCount() (total, installed int) {
	for _, f := range m.fonts {
		if m.installed[f.Name] {
			installed++
		}
	}
	return len(m.fonts), installed
}

func (m *Model) refreshInstalled() {
	inst, err := fonts.InstalledFonts(m.platform.FontDir)
	if err == nil {
		m.installed = inst
	}
	m.buildDisplay()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.termHeight = msg.Height
		m.height = msg.Height - 10
		if m.height < 5 {
			m.height = 5
		}

	case tea.KeyMsg:
		return m.handleKey(msg)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case fontsLoadedMsg:
		if msg.err != nil {
			m.errMsg = fmt.Sprintf("failed to fetch fonts: %v", msg.err)
			m.state = stateList
			return m, nil
		}
		m.fonts = msg.fonts
		m.state = stateList
		m.refreshInstalled()
		return m, nil

	case downloadProgressMsg:
		m.dlDownloaded = msg.downloaded
		m.dlTotal = msg.total
		return m, waitForProgress(m.progressCh)

	case downloadDoneMsg:
		if msg.err != nil {
			m.errMsg = fmt.Sprintf("download failed: %v", msg.err)
		} else {
			m.statusMsg = fmt.Sprintf("installed %s", m.downloading)
		}
		m.state = stateList
		m.downloading = ""
		m.progressCh = nil
		m.refreshInstalled()
		return m, nil

	case uninstallDoneMsg:
		if msg.err != nil {
			m.errMsg = fmt.Sprintf("uninstall failed: %v", msg.err)
		} else {
			m.statusMsg = "font uninstalled"
		}
		m.state = stateList
		m.refreshInstalled()
		return m, nil
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Search mode input handling
	if m.state == stateSearch {
		return m.handleSearchKey(msg)
	}

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "q":
		if m.state == stateConfirmUninstall {
			m.state = stateList
			return m, nil
		}
		if m.state == stateDownloading {
			return m, nil
		}
		return m, tea.Quit
	}

	switch m.state {
	case stateList:
		return m.handleListKey(msg)
	case stateConfirmUninstall:
		return m.handleConfirmKey(msg)
	}
	return m, nil
}

func (m Model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.state = stateList
		return m, nil
	case "escape":
		m.searchQuery = ""
		m.state = stateList
		m.buildDisplay()
		return m, nil
	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.cursor = 0
			m.offset = 0
			m.buildDisplay()
		}
		return m, nil
	default:
		// Only accept printable characters
		r := msg.String()
		if len(r) == 1 && r[0] >= 32 && r[0] <= 126 {
			m.searchQuery += r
			m.cursor = 0
			m.offset = 0
			m.buildDisplay()
		}
		return m, nil
	}
}

func (m Model) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.errMsg = ""
	m.statusMsg = ""

	switch msg.String() {
	case "up", "k":
		m.moveCursorUp()
	case "down", "j":
		m.moveCursorDown()
	case "home", "g":
		m.cursor = 0
		m.offset = 0
		m.clampCursor()
	case "end", "G":
		m.cursor = len(m.display) - 1
		m.clampCursor()
	case "pgup":
		m.cursor -= m.height
		if m.cursor < 0 {
			m.cursor = 0
		}
		m.offset = m.cursor
		m.clampCursor()
	case "pgdown":
		m.cursor += m.height
		if m.cursor >= len(m.display) {
			m.cursor = len(m.display) - 1
		}
		m.clampCursor()
	case "/":
		m.state = stateSearch
		m.searchQuery = ""
		return m, nil
	case "escape":
		if m.searchActive || m.searchQuery != "" {
			m.searchQuery = ""
			m.searchActive = false
			m.cursor = 0
			m.offset = 0
			m.buildDisplay()
		}
	case "enter":
		if len(m.display) == 0 || m.display[m.cursor].font == nil {
			return m, nil
		}
		font := m.display[m.cursor].font
		if m.installed[font.Name] {
			m.state = stateConfirmUninstall
			return m, nil
		}
		return m.startDownload(*font)
	case "d":
		if len(m.display) > 0 && m.cursor < len(m.display) && m.display[m.cursor].font != nil {
			if m.installed[m.display[m.cursor].font.Name] {
				m.state = stateConfirmUninstall
			}
		}
	}

	return m, nil
}

func (m Model) startDownload(font fonts.Font) (tea.Model, tea.Cmd) {
	m.state = stateDownloading
	m.downloading = font.Name
	m.dlDownloaded = 0
	m.dlTotal = 0
	m.progressCh = make(chan downloadProgressMsg, 100)

	ch := m.progressCh
	fontDir := m.platform.FontDir

	downloadCmd := func() tea.Msg {
		err := fonts.InstallFont(fontDir, font.Name, font.DownloadURL, func(downloaded, total int64) {
			select {
			case ch <- downloadProgressMsg{downloaded: downloaded, total: total}:
			default:
			}
		})
		close(ch)
		return downloadDoneMsg{err: err}
	}

	return m, tea.Batch(m.spinner.Tick, downloadCmd, waitForProgress(m.progressCh))
}

func (m Model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		if m.cursor >= len(m.display) || m.display[m.cursor].font == nil {
			m.state = stateList
			return m, nil
		}
		font := m.display[m.cursor].font
		fontDir := m.platform.FontDir
		m.state = stateLoading
		return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
			err := fonts.UninstallFont(fontDir, font.Name)
			return uninstallDoneMsg{err: err}
		})
	case "n", "N", "escape":
		m.state = stateList
	}
	return m, nil
}

// ─── View ───────────────────────────────────────────────────────────

func (m Model) View() string {
	var content strings.Builder

	content.WriteString("\n")

	// Header
	content.WriteString("  " + accentStyle.Render("nfdownloader") + dimStyle.Render(":") + titleStyle.Render("~") + " " + dimStyle.Render("$"))
	content.WriteString("\n\n")

	if !m.platform.Supported {
		content.WriteString("  " + warningStyle.Render("unsupported distro ("+m.platform.Distro+")"))
		content.WriteString("\n\n")
	}

	switch m.state {
	case stateLoading:
		content.WriteString("  " + m.spinner.View() + dimStyle.Render("  fetching fonts..."))
		content.WriteString("\n")

	case stateDownloading:
		content.WriteString("  " + m.spinner.View() + "  downloading " + selectedStyle.Render(m.downloading))
		content.WriteString("\n\n")
		if m.dlTotal > 0 {
			pct := float64(m.dlDownloaded) / float64(m.dlTotal) * 100
			bar := renderProgressBar(pct, 30)
			dlMB := float64(m.dlDownloaded) / 1024 / 1024
			totalMB := float64(m.dlTotal) / 1024 / 1024
			content.WriteString("  " + bar + "  " + pctStyle.Render(fmt.Sprintf("%3.0f%%", pct)))
			content.WriteString(dimStyle.Render(fmt.Sprintf("  %.1f/%.1f MB", dlMB, totalMB)))
		} else {
			content.WriteString("  " + m.spinner.View() + dimStyle.Render("  connecting..."))
		}
		content.WriteString("\n")

	case stateConfirmUninstall:
		if m.cursor < len(m.display) && m.display[m.cursor].font != nil {
			font := m.display[m.cursor].font
			content.WriteString("  " + warningStyle.Render("uninstall "+font.Name+"?"))
			content.WriteString("\n\n")
			content.WriteString("  " + helpKey.Render("y") + helpDesc.Render(" confirm") + "    " + helpKey.Render("n") + helpDesc.Render(" cancel"))
			content.WriteString("\n")
		}

	case stateSearch:
		m.viewList(&content)
		content.WriteString("\n")
		content.WriteString("  " + searchStyle.Render("/") + searchInputStyle.Render(m.searchQuery) + dimStyle.Render("_"))
		content.WriteString("\n")

	case stateList:
		m.viewList(&content)
	}

	// Build footer
	footer := m.viewFooter()

	// Pad content so footer is always at the bottom
	contentStr := content.String()
	contentLines := strings.Count(contentStr, "\n")
	footerLines := strings.Count(footer, "\n")
	totalUsed := contentLines + footerLines
	pad := m.termHeight - totalUsed
	if pad < 0 {
		pad = 0
	}

	return contentStr + strings.Repeat("\n", pad) + footer
}

func (m Model) viewList(b *strings.Builder) {
	if m.errMsg != "" {
		b.WriteString("  " + errorStyle.Render(m.errMsg))
		b.WriteString("\n\n")
	}
	if m.statusMsg != "" {
		b.WriteString("  " + successStyle.Render(m.statusMsg))
		b.WriteString("\n\n")
	}

	if len(m.display) == 0 {
		if m.searchQuery != "" {
			b.WriteString(dimStyle.Render("  no matches for \"" + m.searchQuery + "\""))
		} else {
			b.WriteString(dimStyle.Render("  no fonts available"))
		}
		b.WriteString("\n")
		return
	}

	end := m.offset + m.height
	if end > len(m.display) {
		end = len(m.display)
	}

	for i := m.offset; i < end; i++ {
		item := m.display[i]

		if item.separator != "" {
			label := sectionLabel.Render("  " + item.separator)
			b.WriteString(label)
			b.WriteString("\n")
			continue
		}

		font := item.font
		isInstalled := m.installed[font.Name]

		if i == m.cursor {
			b.WriteString("  " + cursorStyle.Render(">") + " " + selectedStyle.Render(font.Name))
			if isInstalled {
				b.WriteString("  " + installedBadge.Render("installed"))
			}
		} else {
			b.WriteString("    " + normalStyle.Render(font.Name))
			if isInstalled {
				b.WriteString("  " + installedBadge.Render("installed"))
			}
		}
		b.WriteString("\n")
	}

	// Position + scroll hint
	fontIdx, fontTotal := m.cursorFontIndex()
	if fontTotal > 0 {
		b.WriteString("\n")
		pos := fmt.Sprintf("  %d/%d", fontIdx, fontTotal)
		if m.searchQuery != "" && m.state != stateSearch {
			pos += dimStyle.Render("  /" + m.searchQuery)
		}
		b.WriteString(dimStyle.Render(pos))
		b.WriteString("\n")
	}
}

func (m Model) cursorFontIndex() (current, total int) {
	idx := 0
	cur := 0
	for i, item := range m.display {
		if item.font != nil {
			idx++
			if i == m.cursor {
				cur = idx
			}
		}
	}
	return cur, idx
}

func (m Model) viewFooter() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString("  " +
		helpKey.Render("j/k") + helpDesc.Render(" navigate") + "  " +
		helpKey.Render("enter") + helpDesc.Render(" install") + "  " +
		helpKey.Render("d") + helpDesc.Render(" remove") + "  " +
		helpKey.Render("/") + helpDesc.Render(" search") + "  " +
		helpKey.Render("q") + helpDesc.Render(" quit"))
	b.WriteString("\n")
	return b.String()
}

func renderProgressBar(pct float64, width int) string {
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled
	return progressFilled.Render(strings.Repeat("━", filled)) +
		progressEmpty.Render(strings.Repeat("─", empty))
}
