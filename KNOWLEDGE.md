# NFDownloader Knowledge Base

## What This Program Does

NFDownloader is a terminal UI application that lets users browse, install, and remove [Nerd Fonts](https://www.nerdfonts.com/) — patched developer fonts with icon glyphs. It fetches the available font list from the Nerd Fonts website, shows an interactive list, and handles downloading + extracting + installing fonts to the user's local font directory.

## File Responsibilities

### `main.go`
Entry point. Detects the platform, exits early on non-Linux, and starts the Bubble Tea program with alt-screen mode.

### `internal/fonts/fetch.go`
Responsible for fetching and parsing the font list from `https://www.nerdfonts.com/font-downloads`. Uses a regex to find GitHub release download URLs in the HTML. Returns a sorted list of `Font` structs (name + download URL). No external HTML parsing library needed.

### `internal/fonts/download.go`
Handles HTTP file downloads with a progress callback and zip file extraction. The `progressReader` wraps `io.Reader` to report bytes downloaded. The `unzip` function only extracts `.ttf` and `.otf` files, ignoring licenses and readmes.

### `internal/fonts/platform.go`
- **`DetectPlatform()`**: Reads `/etc/os-release` to identify the distro and checks against a list of supported Ubuntu/Debian-based distros.
- **`InstalledFonts()`**: Scans `~/.local/share/fonts/` for subdirectories (each represents an installed font family).
- **`InstallFont()`**: Downloads the zip, extracts font files to `~/.local/share/fonts/<name>/`, and runs `fc-cache`.
- **`UninstallFont()`**: Removes the font directory and refreshes `fc-cache`.

### `internal/ui/model.go`
The Bubble Tea model implementing the full TUI. Key concepts:
- **States**: `stateLoading` (spinner while fetching), `stateList` (browsable font list), `stateDownloading` (progress bar), `stateConfirmUninstall` (y/n prompt), `stateSearch` (vim-style `/` filter).
- **Display list**: `buildDisplay()` creates a grouped view — installed fonts first with a section header, separator, then available fonts. Search filters this list by substring match.
- **Scrolling**: Virtual scroll with `cursor` and `offset`. Cursor skips separator rows automatically.
- **Progress channel**: Download progress is sent via a buffered channel from the download goroutine, consumed by `waitForProgress` commands that feed back into Bubble Tea's update loop.
- **Keyboard**: Arrow keys/j/k for navigation, Enter to install or trigger uninstall, `d` as uninstall shortcut, `/` to search, `escape` to clear search, `g/G` for home/end, `q`/Ctrl+C to quit.

### `internal/ui/styles.go`
Lipgloss style definitions using a dark terminal palette (cyan/teal primary, red accents). Styles for title, selected items, installed badges, section labels, search input, errors, success messages, progress bar, and help text.

## How to Extend

### Adding a new distro
Add the distro ID to the `supported` slice in `internal/fonts/platform.go`. The ID should match the `ID=` or `ID_LIKE=` field in `/etc/os-release`.

### Adding macOS/Windows support
1. Update `DetectPlatform()` to set appropriate `FontDir` for the OS (e.g., `~/Library/Fonts` on macOS).
2. Replace `fc-cache` calls with OS-appropriate font cache refresh (macOS: `atsutil databases -remove`, Windows: may need registry updates).
3. Remove/adjust the Linux-only guard in `main.go`.

### Adding bulk install
Could add multi-select with space bar, tracking selected items in a `map[int]bool`, then batch downloading on Enter.
