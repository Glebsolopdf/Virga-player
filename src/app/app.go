package app

import (
	"sync"
	"sync/atomic"
	"time"

	"virga-player/animation"
	"virga-player/app/state"
	"virga-player/audio"
	debugmgr "virga-player/debug/manager"
	"virga-player/lyricsearch"
	"virga-player/notification"
	"virga-player/rain"
	"virga-player/renderer"
	"virga-player/settings"
	"virga-player/settings/page"

	"github.com/gdamore/tcell/v2"
)

type App struct {
	screen         tcell.Screen
	particleSystem *rain.ParticleSystem
	animEngine     *animation.Engine
	renderEngine   *renderer.Renderer
	state          *state.AppState
	settingsPage   *page.Page
	cfg            *settings.Config
	eventChan      <-chan tcell.Event

	width        int
	height       int
	lastTick     time.Time
	settingsOpen bool
	exitAt       time.Time

	audioAnalyzer     *audio.Analyzer
	debug             *debugmgr.Manager
	debugForced       bool
	notifications     *notification.Service
	notificationToast *notification.Toast

	lyricsResults    chan lyricFetchResult
	lyricsRequestKey string
	lyricsResultKey  string
	currentLyrics    string
	lyricsManager    *lyricsearch.LyricsManager

	lyricsPromptMu sync.Mutex
	lyricsPrompt   *lyricsPromptState

	lyricsDoubleConfirm atomic.Bool

	notificationsSupported    bool
	notifyRemoteInFlight      atomic.Bool
	remoteNotificationResults chan remoteNotificationSyncResult
	uninstallInProgress       atomic.Bool
}

type lyricsPromptState struct {
	trackKey       string
	message        string
	showUntil      time.Time
	firstConfirmAt time.Time
	resultCh       chan bool
}

type Options struct {
	Debug bool
}
