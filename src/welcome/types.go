package welcome

import (
	"errors"

	"virga-player/settings"
)

var ErrCanceled = errors.New("welcome canceled")

type Options struct {
	DefaultConfig *settings.Config
	PlayerConfig  *settings.Config
}

type welcomeStage int

const (
	stageSelect welcomeStage = iota
	stageControls
)

type menuItem struct {
	Title       string
	Description string
	Config      *settings.Config
}
