package install

import (
	"os"
	"os/exec"
	"path/filepath"
)

func ensureSystemCommand(executablePath string) (string, error) {
	if executablePath == systemCommandPath {
		return systemCommandPath, nil
	}

	if err := os.MkdirAll(filepath.Dir(systemCommandPath), 0o755); err != nil {
		return "", err
	}

	data, err := os.ReadFile(executablePath)
	if err != nil {
		return "", err
	}

	tmpPath := systemCommandPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o755); err != nil {
		return "", err
	}
	if err := os.Rename(tmpPath, systemCommandPath); err != nil {
		_ = os.Remove(tmpPath)
		if os.IsPermission(err) {
			return ensureSystemCommandWithSudo(executablePath)
		}
		return "", err
	}

	return systemCommandPath, nil
}

func ensureSystemCommandWithSudo(executablePath string) (string, error) {
	sudoPath, err := exec.LookPath("sudo")
	if err != nil {
		return "", err
	}

	cmd := exec.Command(sudoPath, "install", "-Dm755", executablePath, systemCommandPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return systemCommandPath, nil
}

func removeSystemCommand() error {
	var err error
	if err = os.Remove(systemCommandPath); err == nil {
		return nil
	} else if os.IsNotExist(err) {
		return nil
	} else if os.IsPermission(err) {
		sudoPath, lookErr := exec.LookPath("sudo")
		if lookErr != nil {
			return err
		}

		cmd := exec.Command(sudoPath, "rm", "-f", systemCommandPath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		runErr := cmd.Run()
		if runErr == nil {
			return nil
		}
		return runErr
	}
	return err
}
