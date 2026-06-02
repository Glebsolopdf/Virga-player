package audio

import (
	"errors"
	"os/exec"
	"strings"
)

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
