package playerlyrics

import (
	"time"

	"virga-player/audio"
	"virga-player/renderer"
	"virga-player/settings"

	"github.com/gdamore/tcell/v2"
)

type RenderOptions struct {
	Mode           settings.LyricsDisplayMode
	StaticPosition settings.PositionPreset
	PlayerPosition settings.PositionPreset
	Elapsed        int
	Bands          audio.Bands
}

func (s *State) Render(screen tcell.Screen, r *renderer.Renderer, opts RenderOptions) {
	if s == nil {
		return
	}
	width, height := screen.Size()
	player := PlayerBounds(width, height, opts.PlayerPosition)
	if opts.Mode == settings.LyricsDisplayDynamic {
		s.UpdateDynamic(opts.Elapsed, width, height, player, opts.Bands)
		s.renderDynamic(screen, r)
		return
	}
	s.renderStatic(screen, r, width, height, player, opts)
}

func (s *State) renderStatic(screen tcell.Screen, r *renderer.Renderer, width, height int, player Rect, opts RenderOptions) {
	cue, ok := s.StaticLine(opts.Elapsed)
	if !ok {
		return
	}
	theme := settings.CurrentTheme()
	w, h := lyricsCellSize(cue.Text)
	rect := StaticRect(width, height, w, h, player, opts.StaticPosition)
	drawLyricsText(screen, r, rect.X, rect.Y, cue.Text, theme.LyricsCurrent)
}

func (s *State) renderDynamic(screen tcell.Screen, r *renderer.Renderer) {
	theme := settings.CurrentTheme()
	now := time.Now()
	for _, line := range s.Finished() {
		elapsedFrames := int(now.Sub(line.EndedAt) / (time.Second / 30))
		if elapsedFrames < line.Hold {
			drawLyricsText(screen, r, line.Bounds.X, line.Bounds.Y, line.Cue.Text, theme.LyricsInactive)
			continue
		}
		progress := frameProgress(now.Sub(line.EndedAt)-time.Second/30*time.Duration(line.Hold), line.Wash)
		if text, dx, dy := exitFrame(line.Cue.Text, progress, line.Exit); text != "" {
			drawLyricsText(screen, r, line.Bounds.X+dx, line.Bounds.Y+dy, text, theme.LyricsInactive)
		}
	}
	if line := s.Active(); line != nil {
		progress := frameProgress(now.Sub(line.StartedAt), line.Appear)
		text := appearText(line.Cue.Text, progress, line.Variant)
		drawLyricsText(screen, r, line.Bounds.X, line.Bounds.Y, text, theme.LyricsCurrent)
	}
}

func frameProgress(elapsed time.Duration, frames int) float64 {
	if frames <= 1 {
		return 1
	}
	progress := float64(elapsed) / float64(time.Second/30*time.Duration(frames))
	if progress < 0 {
		return 0
	}
	if progress > 1 {
		return 1
	}
	return progress
}
