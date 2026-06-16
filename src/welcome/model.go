package welcome

import (
	tea "github.com/charmbracelet/bubbletea"
	"virga-player/settings"
)

type model struct {
	width      int
	height     int
	frame      int
	selected   int
	stage      welcomeStage
	items      []menuItem
	chosen     *settings.Config
	result     *ConfigResult
	bannerAnim blockAnimation
	menuAnim   blockAnimation
}

type ConfigResult struct {
	Config   *settings.Config
	Canceled bool
}

func newModel(opts Options) model {
	return model{
		width:  80,
		height: 24,
		bannerAnim: blockAnimation{
			borderFrames: 28,
			settleFrames: 10,
			textFrames:   12,
		},
		menuAnim: blockAnimation{
			borderFrames: 28,
			settleFrames: 10,
			textFrames:   12,
		},
		items: []menuItem{
			{
				Title:       "I just want to watch rain in the terminal",
				Description: "Start Virga with the default atmospheric rain configuration.",
				Config:      opts.DefaultConfig.Clone(),
			},
			{
				Title:       "I want the player visualization",
				Description: "Enable the player-focused layout and the audio integration profile.",
				Config:      opts.PlayerConfig.Clone(),
			},
		},
	}
}

func (m model) Init() tea.Cmd {
	if m.animationsRunning() {
		return nextFrame()
	}
	return nil
}

func (m model) animationsRunning() bool {
	return !m.bannerAnim.state(m.frame).done || !m.menuAnim.state(m.frame).done
}

func (m model) menuReady() bool {
	return m.menuAnim.state(m.frame).done
}
