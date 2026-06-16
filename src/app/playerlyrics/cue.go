package playerlyrics

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var cuePattern = regexp.MustCompile(`\[(\d{1,2}):(\d{2})(?:\.(\d{1,3}))?\](.*)`)

type Cue struct {
	AtMillis int
	Text     string
}

func Parse(raw string) []Cue {
	lines := strings.Split(raw, "\n")
	cues := make([]Cue, 0, len(lines))
	for _, line := range lines {
		match := cuePattern.FindStringSubmatch(strings.TrimSpace(line))
		if len(match) == 0 {
			continue
		}
		text := strings.TrimSpace(match[4])
		if text == "" {
			continue
		}
		cues = append(cues, Cue{AtMillis: cueMillis(match), Text: text})
	}
	sort.SliceStable(cues, func(i, j int) bool {
		return cues[i].AtMillis < cues[j].AtMillis
	})
	return cues
}

func ActiveCue(cues []Cue, elapsedSec int) (Cue, int, bool) {
	elapsed := elapsedSec * 1000
	index := sort.Search(len(cues), func(i int) bool {
		return cues[i].AtMillis > elapsed
	}) - 1
	if index < 0 {
		return Cue{}, -1, false
	}
	return cues[index], index, true
}

func cueMillis(match []string) int {
	minutes, _ := strconv.Atoi(match[1])
	seconds, _ := strconv.Atoi(match[2])
	millis := 0
	if match[3] != "" {
		part := match[3]
		for len(part) < 3 {
			part += "0"
		}
		millis, _ = strconv.Atoi(part[:3])
	}
	return minutes*60000 + seconds*1000 + millis
}
