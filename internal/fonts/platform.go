package fonts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Platform represents the current OS platform and its support status.
type Platform struct {
	OS        string
	Distro    string
	Supported bool
	FontDir   string
}

// DetectPlatform checks the current OS and returns platform info.
func DetectPlatform() Platform {
	p := Platform{OS: runtime.GOOS}

	if p.OS != "linux" {
		return p
	}

	p.Distro = detectDistro()
	supported := []string{"ubuntu", "debian", "pop", "linuxmint", "elementary", "zorin", "neon", "fedora"}
	distroLower := strings.ToLower(p.Distro)
	for _, s := range supported {
		if strings.Contains(distroLower, s) {
			p.Supported = true
			break
		}
	}

	home, err := os.UserHomeDir()
	if err == nil {
		p.FontDir = filepath.Join(home, ".local", "share", "fonts")
	}

	return p
}

func detectDistro() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "unknown"
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "ID_LIKE=") {
			return strings.Trim(strings.TrimPrefix(line, "ID_LIKE="), "\"")
		}
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "ID=") {
			return strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
		}
	}
	return "unknown"
}

// InstalledFonts returns a set of nerd font family names installed in the font directory.
func InstalledFonts(fontDir string) (map[string]bool, error) {
	installed := make(map[string]bool)
	if fontDir == "" {
		return installed, nil
	}

	entries, err := os.ReadDir(fontDir)
	if err != nil {
		if os.IsNotExist(err) {
			return installed, nil
		}
		return nil, err
	}

	for _, e := range entries {
		if e.IsDir() {
			installed[e.Name()] = true
		}
	}
	return installed, nil
}

// InstallFont downloads and installs a font to fontDir/<fontName>/.
func InstallFont(fontDir, fontName, zipURL string, progress func(downloaded, total int64)) error {
	destDir := filepath.Join(fontDir, fontName)
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("create font dir: %w", err)
	}

	zipPath := filepath.Join(os.TempDir(), fontName+".zip")
	defer os.Remove(zipPath)

	if err := downloadFile(zipPath, zipURL, progress); err != nil {
		os.RemoveAll(destDir)
		return fmt.Errorf("download: %w", err)
	}

	if err := unzip(zipPath, destDir); err != nil {
		os.RemoveAll(destDir)
		return fmt.Errorf("extract: %w", err)
	}

	// Refresh font cache
	if cmd, err := exec.LookPath("fc-cache"); err == nil {
		exec.Command(cmd, "-fv", destDir).Run()
	}

	return nil
}

// UninstallFont removes an installed font directory and refreshes the font cache.
func UninstallFont(fontDir, fontName string) error {
	destDir := filepath.Join(fontDir, fontName)
	if err := os.RemoveAll(destDir); err != nil {
		return fmt.Errorf("remove font dir: %w", err)
	}

	if cmd, err := exec.LookPath("fc-cache"); err == nil {
		exec.Command(cmd, "-fv", fontDir).Run()
	}

	return nil
}
