package playerctlcmd

import (
	"context"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const commandTimeout = 750 * time.Millisecond

var (
	checkOnce sync.Once
	available bool
)

func Available() bool {
	checkOnce.Do(func() {
		_, err := exec.LookPath("playerctl")
		available = err == nil
	})
	return available
}

func run(args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
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

func cleanValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || strings.Contains(strings.ToLower(value), "no player could handle this command") {
		return ""
	}
	return value
}
