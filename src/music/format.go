package music

import "fmt"

func formatDuration(seconds int) string {
	mins := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%d:%02d", mins, secs)
}

func (t *TrackInfo) DurationFormatted() string {
	if t.Duration <= 0 {
		return "--:--"
	}
	return formatDuration(t.Duration)
}

func (t *TrackInfo) ElapsedFormatted() string {
	return formatDuration(t.Elapsed)
}
