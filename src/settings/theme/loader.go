package theme

import (
	"os"
	"path/filepath"
)

func LoadOrCreateTheme(path string) (Theme, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			theme := DefaultTheme()
			if writeErr := writeDefaultTheme(path); writeErr != nil {
				return theme, true, writeErr
			}
			SetCurrentTheme(theme)
			return theme, true, nil
		}
		return DefaultTheme(), false, err
	}

	theme := parseThemeCSS(string(data), DefaultTheme())
	SetCurrentTheme(theme)
	return theme, false, nil
}

func writeDefaultTheme(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(defaultThemeCSS), 0o644)
}

func ResetStyleToDefault(path string) error {
	if err := writeDefaultTheme(path); err != nil {
		return err
	}
	SetCurrentTheme(DefaultTheme())
	return nil
}
