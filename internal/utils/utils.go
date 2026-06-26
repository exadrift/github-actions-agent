package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func RunIt(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	oBytes, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if len(exitError.Stderr) > 0 {
				return "", fmt.Errorf("error %d - %s", exitError.ExitCode(), exitError.Stderr)
			}

			return "", fmt.Errorf("error %d", exitError.ExitCode())
		}

		return "", err
	}

	return strings.TrimSpace(string(oBytes)), nil
}
