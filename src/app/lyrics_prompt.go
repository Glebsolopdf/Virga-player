package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"virga-player/lyricsearch"
)

const (
	lyricsPromptVisibleDuration = 6 * time.Second
	lyricsPromptSecondWindow    = 10 * time.Second
)

func trackPromptKey(track lyricsearch.Track) string {
	return strings.TrimSpace(track.Artist) + "\x00" + strings.TrimSpace(track.Title)
}

func (a *App) onLyricsSavePrompt(ctx context.Context, request lyricsearch.PromptRequest) bool {
	if a.debug != nil {
		a.debug.Infof("lyrics prompt: %s", request.Message)
	}

	now := time.Now()
	doubleConfirm := a.lyricsDoubleConfirm.Load()
	message := request.Message
	if doubleConfirm {
		message = "Save lyrics? Press Y, then Y or Enter within 10s"
	} else {
		message = "Save lyrics? Press Y or Enter to confirm"
	}
	state := &lyricsPromptState{
		trackKey:  trackPromptKey(request.Track),
		message:   message,
		showUntil: now.Add(lyricsPromptVisibleDuration),
		resultCh:  make(chan bool, 1),
	}

	a.lyricsPromptMu.Lock()
	a.lyricsPrompt = state
	a.lyricsPromptMu.Unlock()

	defer a.clearLyricsPromptState(state)

	select {
	case approved := <-state.resultCh:
		return approved
	case <-ctx.Done():
		return false
	}
}

func (a *App) clearLyricsPromptState(target *lyricsPromptState) {
	a.lyricsPromptMu.Lock()
	defer a.lyricsPromptMu.Unlock()
	if a.lyricsPrompt == target {
		a.lyricsPrompt = nil
	}
}

func (a *App) tryHandleLyricsPromptConfirm(r rune, enter bool) bool {
	now := time.Now()
	a.lyricsPromptMu.Lock()
	defer a.lyricsPromptMu.Unlock()

	state := a.lyricsPrompt
	if state == nil {
		return false
	}

	if r == 'n' || r == 'N' {
		a.resolveLyricsPromptLocked(state, false)
		if a.debug != nil {
			a.debug.Debugf("lyrics prompt declined")
		}
		return true
	}

	if r == 'y' || r == 'Y' {
		if !a.lyricsDoubleConfirm.Load() {
			a.resolveLyricsPromptLocked(state, true)
			if a.debug != nil {
				a.debug.Infof("lyrics prompt accepted")
			}
			return true
		}

		if state.firstConfirmAt.IsZero() || now.Sub(state.firstConfirmAt) > lyricsPromptSecondWindow {
			state.firstConfirmAt = now
			state.message = "Press Y or Enter again within 10s to save lyrics"
			state.showUntil = now.Add(lyricsPromptVisibleDuration)
			return true
		}

		a.resolveLyricsPromptLocked(state, true)
		if a.debug != nil {
			a.debug.Infof("lyrics prompt accepted")
		}
		return true
	}

	if enter {
		if !a.lyricsDoubleConfirm.Load() {
			a.resolveLyricsPromptLocked(state, true)
			if a.debug != nil {
				a.debug.Infof("lyrics prompt accepted")
			}
			return true
		}

		if state.firstConfirmAt.IsZero() {
			return true
		}
		if now.Sub(state.firstConfirmAt) <= lyricsPromptSecondWindow {
			a.resolveLyricsPromptLocked(state, true)
			if a.debug != nil {
				a.debug.Infof("lyrics prompt accepted")
			}
			return true
		}
		state.firstConfirmAt = time.Time{}
		state.message = "First press Y, then Y or Enter within 10s"
		state.showUntil = now.Add(lyricsPromptVisibleDuration)
		return true
	}

	return false
}

func (a *App) resolveLyricsPromptLocked(state *lyricsPromptState, approved bool) {
	select {
	case state.resultCh <- approved:
	default:
	}
	if a.lyricsPrompt == state {
		a.lyricsPrompt = nil
	}
}

func (a *App) lyricsPromptBanner(now time.Time) (text string, show bool) {
	a.lyricsPromptMu.Lock()
	defer a.lyricsPromptMu.Unlock()

	state := a.lyricsPrompt
	if state == nil {
		return "", false
	}
	if now.After(state.showUntil) {
		return "", false
	}

	text = state.message
	if a.lyricsDoubleConfirm.Load() && !state.firstConfirmAt.IsZero() {
		remaining := lyricsPromptSecondWindow - now.Sub(state.firstConfirmAt)
		if remaining < 0 {
			remaining = 0
		}
		text = fmt.Sprintf("%s (%ds left)", text, int(remaining.Seconds())+1)
	}
	return text, true
}
