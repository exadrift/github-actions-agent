package changelog

import (
	"fmt"
	"os"
	"strings"
)

type ChangeLogCmd struct {
	Location string `arg:"" help:"path to the changelog file"`
	Version  string `arg:"" help:"version to verify"`
}

func (c *ChangeLogCmd) Run() error {
	_, err := os.Stat(c.Location)
	if err != nil {
		return fmt.Errorf("change log file \"%s\" could not be located", c.Location)
	}

	logBytes, err := os.ReadFile(c.Location)
	if err != nil {
		return fmt.Errorf("unable to read changelog %s: %w", c.Location, err)
	}

	return verifyChangeLog(logBytes, c.Version)
}

func verifyChangeLog(logData []byte, version string) error {
	logLines := strings.Split(string(logData), "\n")
	expectVersion := true
	hadMessage := false
	for i, logLine := range logLines {
		if len(strings.Trim(logLine, " \t")) == 0 {
			if i == len(logLines)-1 {
				// skip empty trailing line
				continue
			}
			return fmt.Errorf("changelog should not include blank lines or other superfluous formatting")
		}

		if !expectVersion && hadMessage && strings.HasPrefix(logLine, "#") {
			// prepare for the next series
			expectVersion = true
			hadMessage = false
		}

		if expectVersion {
			if !strings.HasPrefix(logLine, "# v") {
				return fmt.Errorf("changelog version headings must be in the form of \"# v<version>\"")
			}

			if i == 0 {
				// first item must match the expected version
				item := strings.TrimPrefix(logLine, "# ")
				if !strings.HasPrefix(version, "v") {
					item = strings.TrimPrefix(item, "v")
				}
				if item != version {
					return fmt.Errorf("expected version is \"%s\", found version \"%s\"", version, item)
				}
			}

			expectVersion = false
		} else {
			if !strings.HasPrefix(logLine, "- ") {
				return fmt.Errorf("changelog notes must be in the form of \"- <message>\"")
			}

			hadMessage = true
		}
	}

	return nil
}
