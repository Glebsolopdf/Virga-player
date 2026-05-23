package io

import (
	"errors"
	"os/exec"
	"runtime"
)

func CopyToClipboard(text string) error {
	if text == "" {
		return errors.New("no log data to copy")
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		if path, _ := exec.LookPath("wl-copy"); path != "" {
			cmd = exec.Command(path)
		} else if path, _ := exec.LookPath("xclip"); path != "" {
			cmd = exec.Command(path, "-selection", "clipboard")
		} else if path, _ := exec.LookPath("xsel"); path != "" {
			cmd = exec.Command(path, "--clipboard", "--input")
		}
	case "darwin":
		if path, _ := exec.LookPath("pbcopy"); path != "" {
			cmd = exec.Command(path)
		}
	case "windows":
		if path, _ := exec.LookPath("clip"); path != "" {
			cmd = exec.Command(path)
		}
	}

	if cmd == nil {
		return errors.New("clipboard tool not found")
	}

	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		_ = in.Close()
		return err
	}
	_, _ = in.Write([]byte(text))
	_ = in.Close()
	return cmd.Wait()
}
