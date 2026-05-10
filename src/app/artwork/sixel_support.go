package artwork

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// DetectSixelSupport проверяет поддержку SIXEL в терминале.
func DetectSixelSupport() bool {
	termProgram := os.Getenv("TERM_PROGRAM")
	if strings.Contains(termProgram, "iTerm") || strings.Contains(termProgram, "WezTerm") {
		return true
	}

	term := os.Getenv("TERM")
	if strings.Contains(term, "xterm") && strings.Contains(term, "256") {
		return checkXtermSixel()
	}

	if strings.Contains(term, "xterm-kitty") {
		return true
	}

	if _, err := exec.LookPath("convert"); err != nil {
		return false
	}

	return true
}

func checkXtermSixel() bool {
	fmt.Print("\x1B[>q")
	cmd := exec.Command("convert", "-size", "10x10", "xc:red", "sixel:-")
	if err := cmd.Run(); err == nil {
		return true
	}

	return false
}
