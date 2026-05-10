package audio

import (
	"bufio"
	"errors"
	"io"
	"math"
	"os/exec"
	"strings"
	"sync"
)

const (
	sampleRate      = 11025
	readChunkBytes  = 256
	analysisWinSize = 512
	minNormEnvelope = 0.003
	envSensitivity  = 28.0
)

type Bands struct {
	Low      float64
	Mid      float64
	High     float64
	Envelope float64
}

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

	// Rolling window of float64 samples for Goertzel analysis.
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

		// Append new samples; keep only the last analysisWinSize.
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

func detectMonitorSource() (string, error) {
	if _, err := exec.LookPath("parec"); err != nil {
		return "", errors.New("parec is not installed")
	}
	if _, err := exec.LookPath("pactl"); err != nil {
		return "", errors.New("pactl is not installed")
	}

	infoOut, err := exec.Command("pactl", "info").Output()
	if err == nil {
		for _, line := range strings.Split(string(infoOut), "\n") {
			if !strings.HasPrefix(line, "Default Sink:") {
				continue
			}
			sink := strings.TrimSpace(strings.TrimPrefix(line, "Default Sink:"))
			if sink != "" && sink != "(null)" {
				return sink + ".monitor", nil
			}
		}
	}

	sourcesOut, err := exec.Command("pactl", "list", "short", "sources").Output()
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(sourcesOut), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		name := fields[1]
		if strings.HasSuffix(name, ".monitor") {
			return name, nil
		}
	}

	return "", errors.New("no PulseAudio monitor source found")
}

func decodePCM16(buf []byte) []float64 {
	n := len(buf) / 2
	samples := make([]float64, 0, n)
	for i := 0; i+1 < len(buf); i += 2 {
		raw := int16(buf[i]) | int16(buf[i+1])<<8
		samples = append(samples, float64(raw)/32768.0)
	}
	return samples
}

func bandEnergy(samples []float64, freqs []float64) float64 {
	if len(samples) == 0 {
		return 0
	}
	sum := 0.0
	for _, f := range freqs {
		sum += goertzelPower(samples, sampleRate, f)
	}
	return sum / float64(len(freqs))
}

func goertzelPower(samples []float64, rate int, freq float64) float64 {
	n := len(samples)
	if n == 0 {
		return 0
	}
	k := int(0.5 + (float64(n)*freq)/float64(rate))
	omega := 2.0 * math.Pi * float64(k) / float64(n)
	coeff := 2.0 * math.Cos(omega)
	s0, s1, s2 := 0.0, 0.0, 0.0
	for _, x := range samples {
		s0 = x + coeff*s1 - s2
		s2 = s1
		s1 = s0
	}
	power := s1*s1 + s2*s2 - coeff*s1*s2
	if power < 0 {
		return 0
	}
	return power / float64(n)
}

func rms(samples []float64) float64 {
	if len(samples) == 0 {
		return 0
	}
	sum := 0.0
	for _, s := range samples {
		sum += s * s
	}
	return math.Sqrt(sum / float64(len(samples)))
}

func smooth(prev, next, alpha float64) float64 {
	return prev + alpha*(next-prev)
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
