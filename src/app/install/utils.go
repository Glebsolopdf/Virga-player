package install

import (
	"os"
	"path/filepath"
	"strings"
)

func isEphemeralExecutable(executablePath string) bool {
	tempDir := os.TempDir()
	if tempDir != "" && strings.HasPrefix(executablePath, tempDir) {
		return true
	}
	return strings.Contains(executablePath, string(filepath.Separator)+"go-build")
}
