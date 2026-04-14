# TODO

## Phase 1: Baseline Quality

- [ ] Add unit tests for `internal/fonts/fetch.go`
  - Test `parseFonts()` with sample HTML
  - Verify duplicate font links are deduplicated
  - Verify parsed font names are sorted

- [ ] Add unit tests for `internal/fonts/download.go`
  - Test that `unzip()` extracts only `.ttf` and `.otf`
  - Verify non-font files are skipped
  - Verify nested zip paths are flattened safely

- [ ] Add unit tests for `internal/fonts/platform.go`
  - Test distro detection from sample `/etc/os-release` content
  - Test installed font detection from a temp directory
  - Avoid depending on the real machine state in tests

- [ ] Add a basic CI workflow
  - Run `go test ./...`
  - Run `go build ./...`

## Phase 2: Reliability

- [ ] Replace raw `http.Get` calls with an `http.Client`
  - Add request timeout
  - Reuse the same client for page fetches and downloads
  - Consider setting a simple `User-Agent`

- [ ] Handle `fc-cache` failures properly
  - Do not silently ignore errors
  - Surface warnings or errors in the UI after install/uninstall

- [ ] Improve install failure cleanup
  - Verify partial installs are removed cleanly on download/extract failure
  - Check for leftover temp files and directories

- [ ] Add local caching for the fetched font list
  - Cache the last successful result
  - Use cache as fallback if the Nerd Fonts page is unreachable

## Phase 3: UX Improvements

- [ ] Clean up search state handling
  - Remove or properly use `searchActive`
  - Preserve the last search query when reopening search
  - Show match count in the UI

- [ ] Improve download UX
  - Allow canceling an in-progress download
  - Show clearer status when total size is unknown

- [ ] Improve uninstall UX
  - Include the font name in the success message
  - Keep cursor position stable after uninstall

- [ ] Improve empty/error states
  - Better messaging when font fetch fails
  - Better messaging when no fonts match search

## Phase 4: Better Font Detection

- [ ] Make installed-font detection more robust
  - Current approach only checks `~/.local/share/fonts/<FontName>/`
  - Investigate using `fc-list` or a local install manifest
  - Make sure fonts installed outside this tool are detected correctly

## Phase 5: Scraper Hardening

- [ ] Harden HTML scraping
  - Add test fixtures from the Nerd Fonts downloads page
  - Validate parsing against realistic page samples
  - Make regex failures easier to diagnose

- [ ] Reduce dependence on page structure
  - Investigate whether a more stable source exists
  - If not, keep scraping isolated and well-tested

## Phase 6: Product/Docs Consistency

- [ ] Reconcile supported distro messaging
  - Code currently includes Fedora support
  - Docs and project description should match actual behavior

- [ ] Expand README
  - Add screenshots or terminal GIF
  - Add install instructions for built releases
  - Document supported/unsupported platform behavior clearly

## Phase 7: One Bigger Feature

- [ ] Add bulk install/remove
  - Multi-select fonts with keyboard shortcuts
  - Batch install selected fonts
  - Show combined progress/status

## Nice-to-Have

- [ ] Add a details pane for the selected font
- [ ] Add favorites or recently installed fonts
- [ ] Add release builds for Linux binaries
- [ ] Add config options for font directory or cache location
