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
			currentThemeFileContent = defaultThemeCSS
			SetCurrentTheme(theme)
			return theme, true, nil
		}
		return DefaultTheme(), false, err
	}

	theme := parseThemeCSS(string(data), DefaultTheme())
	currentThemeFileContent = string(data)
	SetCurrentTheme(theme)
	return theme, false, nil
}

func ReloadThemeIfChanged(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	if currentThemeFileContent == string(data) {
		return false, nil
	}

	theme := parseThemeCSS(string(data), DefaultTheme())
	SetCurrentTheme(theme)
	currentThemeFileContent = string(data)
	return true, nil
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
	currentThemeFileContent = defaultThemeCSS
	SetCurrentTheme(DefaultTheme())
	return nil
}
