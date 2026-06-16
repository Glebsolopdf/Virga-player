package welcome

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"virga-player/settings"
)

func Run(opts Options) (*settings.Config, error) {
	if opts.DefaultConfig == nil || opts.PlayerConfig == nil {
		return nil, fmt.Errorf("welcome options must include both configs")
	}

	program := tea.NewProgram(newModel(opts), tea.WithAltScreen())
	finalModel, err := program.Run()
	if err != nil {
		return nil, err
	}

	result := finalModel.(model).result
	if result == nil || result.Canceled {
		return nil, ErrCanceled
	}
	return result.Config, nil
}
