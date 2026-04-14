# NFDownloader

A TUI tool to browse, install, and uninstall Nerd Fonts on Ubuntu/Debian-based Linux distros.

## Build & Run

```bash
go build -o nfdownloader .
./nfdownloader
```

## Project Structure

```
main.go                      — Entry point, platform check, launches TUI
internal/
  fonts/
    fetch.go                 — Scrapes nerdfonts.com for font names and download URLs
    download.go              — HTTP download with progress callback, zip extraction
    platform.go              — OS/distro detection, font install/uninstall, installed font detection
  ui/
    model.go                 — Bubble Tea model: states, key handling, view rendering
    styles.go                — Lipgloss style definitions (dark cyan/red terminal theme)
```

## Architecture

- **State machine TUI**: The Bubble Tea model has 5 states: loading, list, downloading, confirmUninstall, search
- **Search**: Vim-style `/` search filters fonts by substring match in real-time
- **Grouped display**: Installed fonts shown at top with section separators, then available fonts
- **Real-time download progress**: Uses a channel to stream progress from the HTTP download goroutine back to the TUI
- **Font detection**: Checks `~/.local/share/fonts/` for subdirectories matching font names
- **Platform guard**: Only runs on Linux; shows a warning banner for non-Ubuntu/Debian distros

## Key Decisions

- Fonts are installed to `~/.local/share/fonts/<FontName>/` (user-local, no sudo needed)
- Only .ttf and .otf files are extracted from the zip
- `fc-cache` is called after install/uninstall to refresh the system font cache
- Font list is scraped from the HTML page using regex on GitHub release URLs (no API key needed)

## Supported Distros

Ubuntu, Debian, Pop!_OS, Linux Mint, elementary OS, Zorin OS, KDE Neon. Adding more requires updating the `supported` slice in `platform.go`.
