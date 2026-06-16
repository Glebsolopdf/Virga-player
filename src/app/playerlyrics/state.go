package playerlyrics

import (
	"hash/fnv"
	"math/rand"
	"time"

	"virga-player/audio"
)

type State struct {
	raw       string
	cues      []Cue
	active    *dynamicLine
	finished  []dynamicLine
	lastIndex int
	rng       *rand.Rand
}

type dynamicLine struct {
	Cue       Cue
	Bounds    Rect
	StartedAt time.Time
	EndedAt   time.Time
	Appear    int
	Hold      int
	Wash      int
	Variant   int
	Exit      int
}

func NewState() *State {
	return &State{lastIndex: -1, rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (s *State) SetRaw(raw string) {
	if s.raw == raw {
		return
	}
	s.raw = raw
	s.cues = Parse(raw)
	s.active = nil
	s.finished = nil
	s.lastIndex = -1
}

func (s *State) StaticLine(elapsed int) (Cue, bool) {
	cue, _, ok := ActiveCue(s.cues, elapsed)
	return cue, ok
}

func (s *State) UpdateDynamic(elapsed int, width, height int, player Rect, bands audio.Bands) {
	cue, index, ok := ActiveCue(s.cues, elapsed)
	if ok && index != s.lastIndex {
		s.finishActive()
		s.active = s.newLine(cue, width, height, player, Intensity(bands))
		s.lastIndex = index
	}
	s.pruneFinished(time.Now())
}

func (s *State) Active() *dynamicLine {
	return s.active
}

func (s *State) Finished() []dynamicLine {
	return s.finished
}

func (s *State) finishActive() {
	if s.active == nil {
		return
	}
	line := *s.active
	line.EndedAt = time.Now()
	s.finished = append(s.finished, line)
	s.active = nil
}

func (s *State) pruneFinished(now time.Time) {
	kept := s.finished[:0]
	for _, line := range s.finished {
		if int(now.Sub(line.EndedAt)/(time.Second/30)) <= line.Hold+line.Wash {
			kept = append(kept, line)
		}
	}
	s.finished = kept
}

func (s *State) newLine(cue Cue, width, height int, player Rect, pressure float64) *dynamicLine {
	textW, textH := lyricsCellSize(cue.Text)
	rng := seededRand(s.rng, cue)
	rect := nextPlacement(rng, s.active, cue, width, height, player, textW, textH)
	return &dynamicLine{
		Cue:       cue,
		Bounds:    rect,
		StartedAt: time.Now(),
		Appear:    appearFrames(pressure),
		Hold:      6,
		Wash:      washFrames(pressure),
		Variant:   rng.Intn(3),
		Exit:      rng.Intn(3),
	}
}

func seededRand(parent *rand.Rand, cue Cue) *rand.Rand {
	h := fnv.New64a()
	_, _ = h.Write([]byte(cue.Text))
	seed := int64(h.Sum64()) ^ int64(cue.AtMillis) ^ parent.Int63()
	return rand.New(rand.NewSource(seed))
}
