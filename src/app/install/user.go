package install

import "os"

func ensureUserCommand(executablePath, destination string) (string, error) {
	if executablePath == destination {
		return destination, nil
	}

	data, err := os.ReadFile(executablePath)
	if err != nil {
		return "", err
	}

	tmpPath := destination + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o755); err != nil {
		return "", err
	}
	if err := os.Rename(tmpPath, destination); err != nil {
		_ = os.Remove(tmpPath)
		return "", err
	}

	return destination, nil
}

func removeUserCommand(destination string) error {
	if err := os.Remove(destination); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func ensureAlias(aliasPath, executablePath string) error {
	info, err := os.Lstat(aliasPath)
	if err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			target, readErr := os.Readlink(aliasPath)
			if readErr == nil && target == executablePath {
				return nil
			}
			if removeErr := os.Remove(aliasPath); removeErr != nil {
				return removeErr
			}
		} else {
			return nil
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	return os.Symlink(executablePath, aliasPath)
}
