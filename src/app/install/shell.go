package install

import (
	"os"
	"strings"
)

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
		if trimmed == "# Added by Virga Player" || strings.Contains(trimmed, virgaPathExport) {
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
