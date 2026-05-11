package install

import (
	"os"
	"path/filepath"
	"strings"
)

const virgaPathExport = `export PATH="$HOME/.local/bin:$PATH"`
const systemCommandPath = "/usr/bin/virga"

func EnsureCommandAliases() bool {
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

	aliasTarget := executablePath
	systemInstalled := false
	if installPath, err := ensureSystemCommand(executablePath); err == nil {
		aliasTarget = installPath
		systemInstalled = true
		_ = removeUserCommand(filepath.Join(binDir, "virga"))
	} else if userPath, err := ensureUserCommand(executablePath, filepath.Join(binDir, "virga")); err == nil {
		aliasTarget = userPath
	}

	if err := ensureShellPathEntry(filepath.Join(homeDir, ".profile")); err != nil {
		return false
	}
	if err := ensureShellPathEntry(filepath.Join(homeDir, ".bash_profile")); err != nil {
		return false
	}
	if err := ensureShellPathEntry(filepath.Join(homeDir, ".bash_login")); err != nil {
		return false
	}
	if err := ensureShellPathEntry(filepath.Join(homeDir, ".bashrc")); err != nil {
		return false
	}
	if err := ensureShellPathEntry(filepath.Join(homeDir, ".zshrc")); err != nil {
		return false
	}

	aliases := []string{"virgaplayer"}
	if !systemInstalled {
		aliases = append([]string{"virga"}, aliases...)
	}
	for _, alias := range aliases {
		aliasPath := filepath.Join(binDir, alias)
		if err := ensureAlias(aliasPath, aliasTarget); err != nil {
			return false
		}
	}

	return true
}

func RemoveVirgaInstallation() error {
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
	_ = removeUserCommand(filepath.Join(binDir, "virga"))
	for _, alias := range []string{"virga", "virgaplayer"} {
		aliasPath := filepath.Join(binDir, alias)
		if info, statErr := os.Lstat(aliasPath); statErr == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				_ = os.Remove(aliasPath)
			}
		}
	}

	_ = removeSystemCommand()
	_ = removeShellPathEntry(filepath.Join(homeDir, ".profile"))
	_ = removeShellPathEntry(filepath.Join(homeDir, ".bash_profile"))
	_ = removeShellPathEntry(filepath.Join(homeDir, ".bash_login"))
	_ = removeShellPathEntry(filepath.Join(homeDir, ".bashrc"))
	_ = removeShellPathEntry(filepath.Join(homeDir, ".zshrc"))
	_ = removePathFromEnv(binDir)

	return nil
}

func removePathFromEnv(pathToRemove string) error {
	currentPath := os.Getenv("PATH")
	if currentPath == "" {
		return nil
	}

	parts := strings.Split(currentPath, string(os.PathListSeparator))
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == pathToRemove {
			continue
		}
		filtered = append(filtered, part)
	}

	return os.Setenv("PATH", strings.Join(filtered, string(os.PathListSeparator)))
}
