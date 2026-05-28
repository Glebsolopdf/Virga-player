package artwork

import (
	"os"
	"os/exec"
	"strings"
)

// DetectSixelSupport
func DetectSixelSupport() bool {
	term := strings.ToLower(os.Getenv("TERM"))
	termProgram := strings.ToLower(os.Getenv("TERM_PROGRAM"))

	// kitty does not support sixel output; it uses its own image protocol.
	if strings.Contains(term, "kitty") || strings.Contains(termProgram, "kitty") {
		return false
	}

	if _, err := exec.LookPath("convert"); err != nil {
		return false
	}

	if strings.Contains(term, "sixel") {
		return true
	}

	if strings.Contains(termProgram, "iterm") || strings.Contains(termProgram, "wezterm") {
		return true
	}

	if strings.Contains(os.Getenv("XTERM_VERSION"), "sixel") {
		return true
	}

	return false
}
