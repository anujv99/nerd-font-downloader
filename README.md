# NFDownloader

NFDownloader is a terminal UI for browsing, installing, and uninstalling [Nerd Fonts](https://www.nerdfonts.com/) from Linux.

It fetches the available font list from the Nerd Fonts downloads page, shows installed fonts separately from available fonts, and installs selected fonts into your user font directory.

## Features

- Browse Nerd Fonts in a terminal UI built with Bubble Tea
- Detect already installed fonts from `~/.local/share/fonts`
- Install fonts without `sudo`
- Uninstall previously installed fonts from the app
- Vim-style `/` search with live filtering
- Real-time download progress during installs

## Platform Support

- Linux only
- Supported distros currently include: Ubuntu, Debian, Pop!_OS, Linux Mint, elementary OS, Zorin OS, KDE Neon, and Fedora
- On unsupported Linux distros, the app still starts but shows a warning banner
- On non-Linux platforms, the program exits immediately

## Installation Location

Fonts are installed to:

```text
~/.local/share/fonts/<FontName>/
```

Only `.ttf` and `.otf` files are extracted from downloaded archives.

## Requirements

- Go `1.26.1`
- Linux
- `fc-cache` available on the system
- Network access to `https://www.nerdfonts.com/font-downloads` and the Nerd Fonts GitHub releases

## Build And Run

```bash
go build -o nfdownloader .
./nfdownloader
```

## Controls

- `j` / `k` or arrow keys: move through the list
- `Enter`: install the selected font, or prompt to uninstall if it is already installed
- `d`: uninstall the selected installed font
- `/`: start search
- `Esc`: clear search / leave search mode
- `g`: jump to top
- `G`: jump to bottom
- `PgUp` / `PgDn`: page up/down
- `q` or `Ctrl+C`: quit

## How It Works

- Fetches the font list by scraping the Nerd Fonts downloads page for GitHub release zip URLs
- Groups installed fonts above available fonts in the UI
- Downloads the selected zip file to a temporary location
- Extracts font files into the user-local font directory
- Refreshes the font cache with `fc-cache`

## Project Structure

```text
main.go                      Entry point and platform guard
internal/fonts/fetch.go      Fetches and parses Nerd Fonts download links
internal/fonts/download.go   Downloads zip files and extracts font files
internal/fonts/platform.go   Platform detection and install/uninstall logic
internal/ui/model.go         Bubble Tea state machine and key handling
internal/ui/styles.go        Lipgloss styles for the TUI
```

## Notes

- This project currently uses HTML scraping rather than a formal Nerd Fonts API
- The font list depends on the structure of the Nerd Fonts downloads page remaining compatible
