package changelog

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestChangeLogVersionMatch(t *testing.T) {
	changelogBytes := []byte(`# v1.0.0
- major update
# v0.5.0
- some minor change
# v0.0.1
- initial release
`)
	err := verifyChangeLog(changelogBytes, "v1.0.0")
	assert.NoError(t, err)
}

func TestChangeLogVersionIncorrect(t *testing.T) {
	changelogBytes := []byte(`# v1.0.0
- major update
# v0.5.0
- some minor change
# v0.0.1
- initial release
`)
	err := verifyChangeLog(changelogBytes, "v2.0.0")
	assert.Error(t, err)
}

func TestChangeLogFormattingIncorrect(t *testing.T) {
	changelogBytes := []byte(`#v1.0.0
- major update
# v0.5.0
- some minor change
# v0.0.1
- initial release
`)
	err := verifyChangeLog(changelogBytes, "v1.0.0")
	assert.Error(t, err)
}

func TestChangeLogSuperfluousSpacing(t *testing.T) {
	changelogBytes := []byte(`# v1.0.0
- major update
# v0.5.0

- some minor change
# v0.0.1
- initial release
`)
	err := verifyChangeLog(changelogBytes, "v1.0.0")
	assert.Error(t, err)
}
