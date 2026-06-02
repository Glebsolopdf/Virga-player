package artwork

import (
	"sort"
	"strconv"
	"strings"
)

func parseSyncedLyrics(raw string) []lyricCue {
	lines := strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n")
	cues := make([]lyricCue, 0, len(lines))
	for _, line := range lines {
		rest := strings.TrimSpace(line)
		if rest == "" {
			continue
		}

		timestamps := make([]int, 0, 2)
		for strings.HasPrefix(rest, "[") {
			end := strings.IndexByte(rest, ']')
			if end <= 1 {
				break
			}

			if atMillis, ok := parseLyricTimestamp(rest[1:end]); ok {
				timestamps = append(timestamps, atMillis)
			}
			rest = strings.TrimLeft(rest[end+1:], " \t")
		}

		text := strings.TrimSpace(rest)
		if len(timestamps) == 0 || text == "" {
			continue
		}

		for _, atMillis := range timestamps {
			cues = append(cues, lyricCue{atMillis: atMillis, text: text})
		}
	}

	sort.SliceStable(cues, func(i, j int) bool {
		return cues[i].atMillis < cues[j].atMillis
	})
	return cues
}

func parseLyricTimestamp(tag string) (int, bool) {
	parts := strings.SplitN(strings.TrimSpace(tag), ":", 2)
	if len(parts) != 2 {
		return 0, false
	}

	minutes, err := strconv.Atoi(parts[0])
	if err != nil || minutes < 0 {
		return 0, false
	}

	secondsPart := strings.ReplaceAll(strings.TrimSpace(parts[1]), ",", ".")
	seconds, err := strconv.ParseFloat(secondsPart, 64)
	if err != nil || seconds < 0 {
		return 0, false
	}

	return minutes*60000 + int(seconds*1000+0.5), true
}
