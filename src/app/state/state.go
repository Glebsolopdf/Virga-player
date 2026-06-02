package state

import (
	"time"

	"virga-player/app/message"
	"virga-player/app/player"
	"virga-player/settings"
)

type AppState struct {
	Width          int
	Height         int
	Message        *message.Message
	NextMessage    string
	IntroStep      int
	IntroStepStart time.Time
	Player         *player.Player
	PlayerEnabled  bool
	WashTriggered  bool
	WashStart      time.Time
	StartTime      time.Time
}

func NewAppState(width, height int, text, nextMessage string, cfg *settings.Config) *AppState {
	p := &player.Player{}
	if cfg.Player {
		p = player.New(width/2, height/2)
		p.UpdatePosition(width, height)
		p.SetTrackInfo("Loading...", "Please wait", "...", 0, 0)
	}
	return &AppState{
		Width:          width,
		Height:         height,
		Message:        message.New(text, width, height),
		NextMessage:    nextMessage,
		IntroStep:      introStep(nextMessage),
		IntroStepStart: time.Now(),
		Player:         p,
		PlayerEnabled:  cfg.Player,
		StartTime:      time.Now(),
	}
}

func introStep(nextMessage string) int {
	if nextMessage != "" {
		return 1
	}
	return 0
}

func (s *AppState) Resize(width, height int) {
	s.Width = width
	s.Height = height
	s.Message.UpdatePosition(width, height)
	if s.PlayerEnabled && s.Player != nil {
		s.Player.UpdatePosition(width, height)
	}
}

func (s *AppState) TriggerWash() {
	s.WashTriggered = true
	s.WashStart = time.Now()
}

func (s *AppState) TriggerWashNow() {
	s.WashTriggered = true
	s.WashStart = time.Now().Add(-5 * time.Second)
}

func (s *AppState) AdvanceMessage() bool {
	if s.NextMessage == "" {
		return false
	}

	s.Message.SetText(s.NextMessage, s.Width, s.Height)
	s.NextMessage = ""
	s.WashTriggered = false
	s.WashStart = time.Time{}
	s.StartTime = time.Now()
	return true
}

func (s *AppState) SinceStart() time.Duration {
	return time.Since(s.StartTime)
}

func (s *AppState) ShouldSpawnMessageDrops() bool {
	return s.WashTriggered && !s.Message.Converted && !s.Message.Persistent && time.Since(s.WashStart) >= 5*time.Second
}

func (s *AppState) IsIntroActive() bool {
	return s.IntroStep != 0
}

func (s *AppState) IsMessageProtected() bool {
	return s.IntroStep == 1 || s.IntroStep == 2
}

func (s *AppState) ShouldPauseRain() bool {
	return s.IsMessageProtected()
}

func (s *AppState) AdvanceIntroIfNeeded() (startedWash bool) {
	if s.IntroStep == 0 {
		return false
	}

	if time.Since(s.IntroStepStart) < 5*time.Second {
		return false
	}

	if s.IntroStep == 1 {
		if s.AdvanceMessage() {
			s.IntroStep = 2
			s.IntroStepStart = time.Now()
			return false
		}
		s.IntroStep = 0
		return false
	}

	if s.IntroStep == 2 {
		s.IntroStep = 0
		return true
	}

	return false
}
