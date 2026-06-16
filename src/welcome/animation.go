package welcome

import (
	"math"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lucasb-eyer/go-colorful"
)

const frameStep = time.Second / 30

type tickMsg time.Time

type blockAnimation struct {
	delay        int
	borderFrames int
	settleFrames int
	textFrames   int
}

type blockState struct {
	borderProgress float64
	settleProgress float64
	textProgress   float64
	done           bool
}

func nextFrame() tea.Cmd {
	return tea.Tick(frameStep, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (a blockAnimation) state(frame int) blockState {
	elapsed := frame - a.delay
	if elapsed < 0 {
		return blockState{}
	}

	borderDone := a.borderFrames
	settleDone := borderDone + a.settleFrames
	textDone := settleDone + a.textFrames

	switch {
	case elapsed < borderDone:
		return blockState{borderProgress: progress(elapsed+1, a.borderFrames)}
	case elapsed < settleDone:
		return blockState{
			borderProgress: 1,
			settleProgress: progress(elapsed-borderDone+1, a.settleFrames),
		}
	case elapsed < textDone:
		return blockState{
			borderProgress: 1,
			settleProgress: 1,
			textProgress:   progress(elapsed-settleDone+1, a.textFrames),
		}
	default:
		return blockState{
			borderProgress: 1,
			settleProgress: 1,
			textProgress:   1,
			done:           true,
		}
	}
}

func progress(step, total int) float64 {
	if total <= 0 {
		return 1
	}
	value := float64(step) / float64(total)
	return math.Max(0, math.Min(1, value))
}

func animatedBorderColor(index, total int, settle, phase float64) string {
	rainbow := colorful.Hsv(math.Mod(phase*10+float64(index)*360/float64(total), 360), 0.78, 1)
	white := colorful.LinearRgb(1, 1, 1)
	color := rainbow.BlendLab(white, settle)
	return color.Clamped().Hex()
}

func fadeColor(from, to colorful.Color, amount float64) string {
	blended := from.BlendLab(to, math.Max(0, math.Min(1, amount)))
	return blended.Clamped().Hex()
}
