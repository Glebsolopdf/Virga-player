package music

import (
	"context"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const playerctlTimeout = 750 * time.Millisecond

var (
	playerctlCheck sync.Once
	playerctlAvail bool
)

func playerctlAvailable() bool {
	playerctlCheck.Do(func() {
		_, err := exec.LookPath("playerctl")
		playerctlAvail = err == nil
	})
	return playerctlAvail
}

func playerctlRun(args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), playerctlTimeout)
	defer cancel()

	out, err := exec.CommandContext(ctx, "playerctl", args...).Output()
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func playerctlClean(val string) string {
	val = strings.TrimSpace(val)
	if val == "" || strings.Contains(strings.ToLower(val), "no player could handle this command") {
		return ""
	}
	return val
}
