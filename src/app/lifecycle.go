package app

import (
	"time"

	"virga-player/app/bootstrap"
	debugmgr "virga-player/debug/manager"
)

func New(opts Options, dbg *debugmgr.Manager) *App {
	if dbg == nil {
		dbg = debugmgr.NewManager(opts.Debug, opts.Debug)
	}
	return &App{
		debug:       dbg,
		debugForced: opts.Debug,
	}
}

func (a *App) Run() error {
	var err error
	a.screen, err = bootstrap.NewScreen()
	if err != nil {
		return err
	}
	defer a.screen.Fini()
	defer func() {
		if a.audioAnalyzer != nil {
			a.audioAnalyzer.Stop()
		}
	}()

	a.initComponents()
	a.initEvents()
	a.lastTick = time.Now()

	for {
		select {
		case now := <-a.animEngine.Tick():
			dt := now.Sub(a.lastTick).Seconds()
			if dt <= 0 {
				dt = a.animEngine.FrameDuration().Seconds()
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
