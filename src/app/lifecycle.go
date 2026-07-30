package app

import (
	"time"

	debugmgr "virga-player/debug/manager"

	"github.com/gdamore/tcell/v2"
)

func newScreen() (tcell.Screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := s.Init(); err != nil {
		return nil, err
	}
	s.EnableMouse()
	s.SetStyle(tcell.StyleDefault.Background(tcell.ColorReset))
	s.Clear()
	return s, nil
}

func eventChan(screen tcell.Screen) <-chan tcell.Event {
	ch := make(chan tcell.Event, 16)
	go func() {
		for {
			e := screen.PollEvent()
			if e == nil {
				close(ch)
				return
			}
			ch <- e
		}
	}()
	return ch
}

func New(opts Options, dbg *debugmgr.Manager) *App {
	if dbg == nil {
		dbg = debugmgr.NewManager(opts.Debug, opts.Debug)
	}
	return &App{
		debug:         dbg,
		debugForced:   opts.Debug,
		lyricsResults: make(chan lyricFetchResult, 4),
	}
}

func (a *App) Run() error {
	var err error
	a.screen, err = newScreen()
	if err != nil {
		return err
	}
	defer a.screen.Fini()
	defer func() {
		if a.audioAnalyzer != nil {
			a.audioAnalyzer.Stop()
		}
		a.closeLyricsManager()
	}()

	a.initComponents()
	a.initEvents()
	a.lastTick = time.Now()

	for {
		select {
		case now := <-a.animTicker.C:
			dt := now.Sub(a.lastTick).Seconds()
			if dt <= 0 {
				dt = (time.Second / time.Duration(a.cfg.FPS)).Seconds()
			}
			a.lastTick = now
			if !a.exitAt.IsZero() && now.After(a.exitAt) {
				return nil
			}
			a.onTick(dt)
		case event := <-a.eventChan:
			if a.handleEvent(event) {
				return nil
			}
		}
	}
}
