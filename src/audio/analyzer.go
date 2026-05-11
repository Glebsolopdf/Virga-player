package audio

import (
	"bufio"
	"errors"
	"io"
	"math"
	"os/exec"
	"sync"
)

type Analyzer struct {
	mu        sync.RWMutex
	bands     Bands
	envelope  float64
	cmd       *exec.Cmd
	stopOnce  sync.Once
	stopCh    chan struct{}
	doneCh    chan struct{}
	available bool
}

func NewAnalyzer() (*Analyzer, error) {
	source, err := detectMonitorSource()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(
		"parec",
		"-d", source,
		"--format=s16le",
		"--rate", "11025",
		"--channels", "1",
		"--latency-msec", "10",
		"--raw",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	a := &Analyzer{
		cmd:       cmd,
		stopCh:    make(chan struct{}),
		doneCh:    make(chan struct{}),
		available: true,
	}

	go a.readLoop(stdout)
	return a, nil
}

func (a *Analyzer) readLoop(r io.Reader) {
	defer close(a.doneCh)
	reader := bufio.NewReader(r)
	readBuf := make([]byte, readChunkBytes)

	win := make([]float64, 0, analysisWinSize)

	for {
		select {
		case <-a.stopCh:
			return
		default:
		}

		n, err := io.ReadFull(reader, readBuf)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				return
			}
			return
		}
		if n < 2 {
			continue
		}

		chunk := decodePCM16(readBuf[:n])

		win = append(win, chunk...)
		if len(win) > analysisWinSize {
			win = win[len(win)-analysisWinSize:]
		}
		if len(win) < 64 {
			continue
		}

		low := bandEnergy(win, []float64{60, 80, 120, 180})
		mid := bandEnergy(win, []float64{500, 800, 1300, 2000})
		high := bandEnergy(win, []float64{2800, 3600, 4400, 5000})
		env := rms(chunk)

		normEnv := clamp((env-minNormEnvelope)*envSensitivity, 0, 1)
		nLow := clamp(math.Sqrt(low)*envSensitivity*0.6, 0, 1)
		nMid := clamp(math.Sqrt(mid)*envSensitivity*0.6, 0, 1)
		nHigh := clamp(math.Sqrt(high)*envSensitivity*0.6, 0, 1)

		a.mu.Lock()
		step := func(prev, next float64) float64 {
			if next > prev {
				return next
			}
			return smooth(prev, next, 0.40)
		}
		a.envelope = step(a.envelope, normEnv)
		a.bands.Low = step(a.bands.Low, nLow)
		a.bands.Mid = step(a.bands.Mid, nMid)
		a.bands.High = step(a.bands.High, nHigh)
		a.mu.Unlock()
	}
}

func (a *Analyzer) Bands() (Bands, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if !a.available {
		return Bands{}, false
	}
	b := a.bands
	b.Envelope = a.envelope
	return b, true
}

func (a *Analyzer) Stop() {
	a.stopOnce.Do(func() {
		close(a.stopCh)
		if a.cmd != nil && a.cmd.Process != nil {
			_ = a.cmd.Process.Kill()
		}
		if a.cmd != nil {
			_ = a.cmd.Wait()
		}
		<-a.doneCh
		a.mu.Lock()
		a.available = false
		a.mu.Unlock()
	})
}
