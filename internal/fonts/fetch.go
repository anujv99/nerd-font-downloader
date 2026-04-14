package fonts

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

const NerdFontsURL = "https://www.nerdfonts.com/font-downloads"

// Font represents a downloadable nerd font.
type Font struct {
	Name       string
	DownloadURL string
}

// FetchFontList scrapes the nerd fonts website and returns available fonts.
func FetchFontList() ([]Font, error) {
	resp, err := http.Get(NerdFontsURL)
	if err != nil {
		return nil, fmt.Errorf("fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return parseFonts(string(body))
}

func parseFonts(html string) ([]Font, error) {
	// The nerd fonts page has download links like:
	// https://github.com/ryanoasis/nerd-fonts/releases/download/v3.x.x/FontName.zip
	re := regexp.MustCompile(`https://github\.com/ryanoasis/nerd-fonts/releases/download/[^"'\s]+\.zip`)
	matches := re.FindAllString(html, -1)

	if len(matches) == 0 {
		return nil, fmt.Errorf("no font download links found on page")
	}

	seen := make(map[string]bool)
	var fonts []Font
	for _, url := range matches {
		// Extract font name from URL: .../FontName.zip
		parts := strings.Split(url, "/")
		filename := parts[len(parts)-1]
		name := strings.TrimSuffix(filename, ".zip")

		if seen[name] {
			continue
		}
		seen[name] = true

		fonts = append(fonts, Font{
			Name:       name,
			DownloadURL: url,
		})
	}

	sort.Slice(fonts, func(i, j int) bool {
		return fonts[i].Name < fonts[j].Name
	})

	return fonts, nil
}
