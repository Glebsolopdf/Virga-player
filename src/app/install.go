package app

import (
	"os"
	"path/filepath"
	"strings"
)

const virgaPathExport = "export PATH=\"$HOME/.local/bin:$PATH\""

func ensureCommandAliases() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	executablePath, err := os.Executable()
	if err != nil {
		return false
	}
	if resolvedPath, resolveErr := filepath.EvalSymlinks(executablePath); resolveErr == nil {
		executablePath = resolvedPath
	}
	if isEphemeralExecutable(executablePath) {
		return false
	}

	binDir := filepath.Join(homeDir, ".local", "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return false
	}

	if currentPath := os.Getenv("PATH"); !strings.Contains(currentPath, binDir) {
		_ = os.Setenv("PATH", binDir+string(os.PathListSeparator)+currentPath)
	}

	if err := ensureShellPathEntry(filepath.Join(homeDir, ".profile")); err != nil {
		return false
	}
	if err := ensureShellPathEntry(filepath.Join(homeDir, ".bashrc")); err != nil {
		return false
	}
	if err := ensureShellPathEntry(filepath.Join(homeDir, ".zshrc")); err != nil {
		return false
	}

	for _, alias := range []string{"virga", "virgaplayer"} {
		aliasPath := filepath.Join(binDir, alias)
		if err := ensureAlias(aliasPath, executablePath); err != nil {
			return false
		}
	}

	return true
}

func removeVirgaInstallation() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = filepath.Join(homeDir, ".config")
	}
	_ = os.RemoveAll(filepath.Join(configHome, "virga-player"))

	binDir := filepath.Join(homeDir, ".local", "bin")
	for _, alias := range []string{"virga", "virgaplayer"} {
		aliasPath := filepath.Join(binDir, alias)
		if info, statErr := os.Lstat(aliasPath); statErr == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				_ = os.Remove(aliasPath)
			}
		}
	}

	_ = removeShellPathEntry(filepath.Join(homeDir, ".profile"))
	_ = removeShellPathEntry(filepath.Join(homeDir, ".bashrc"))
	_ = removeShellPathEntry(filepath.Join(homeDir, ".zshrc"))

	return nil
}

func ensureShellPathEntry(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err == nil {
		if strings.Contains(string(data), virgaPathExport) {
			return nil
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	content := ""
	if err == nil {
		content = string(data)
		if content != "" && !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
	}

	content += "\n# Added by Virga Player\n" + virgaPathExport + "\n"
	return os.WriteFile(filePath, []byte(content), 0o644)
}

func ensureAlias(aliasPath, executablePath string) error {
	info, err := os.Lstat(aliasPath)
	if err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			target, readErr := os.Readlink(aliasPath)
			if readErr == nil && target == executablePath {
				return nil
			}
			if removeErr := os.Remove(aliasPath); removeErr != nil {
				return removeErr
			}
		} else {
			return nil
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	return os.Symlink(executablePath, aliasPath)
}

func removeShellPathEntry(filePath string) error {
	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == virgaPathExport || trimmed == "# Added by Virga Player" {
			continue
		}
		filtered = append(filtered, line)
	}

	content := strings.Join(filtered, "\n")
	for strings.Contains(content, "\n\n\n") {
		content = strings.ReplaceAll(content, "\n\n\n", "\n\n")
	}

	return os.WriteFile(filePath, []byte(content), 0o644)
}

func isEphemeralExecutable(executablePath string) bool {
	tempDir := os.TempDir()
	if tempDir != "" && strings.HasPrefix(executablePath, tempDir) {
		return true
	}
	return strings.Contains(executablePath, string(filepath.Separator)+"go-build")
}
