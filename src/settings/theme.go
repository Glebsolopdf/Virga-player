package settings

import (
	"path/filepath"

	theme "virga-player/settings/theme"
)

type Theme = theme.Theme

func LoadOrCreateTheme() (Theme, bool, error) {
	return theme.LoadOrCreateTheme(filepath.Join(ConfigDirPath(), "style.css"))
}

func ReloadThemeIfChanged() (bool, error) {
	return theme.ReloadThemeIfChanged(filepath.Join(ConfigDirPath(), "style.css"))
}

func DefaultTheme() Theme {
	return theme.DefaultTheme()
}

func SetCurrentTheme(t Theme) {
	theme.SetCurrentTheme(t)
}

func CurrentTheme() Theme {
	return theme.CurrentTheme()
}

func StylePath() string {
	return filepath.Join(ConfigDirPath(), "style.css")
}

func ResetStyleToDefault() error {
	return theme.ResetStyleToDefault(filepath.Join(ConfigDirPath(), "style.css"))
}
