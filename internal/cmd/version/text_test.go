package version

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestTextVersion(t *testing.T) {
	textBytes := []byte("v0.1.2\n")
	versioner, err := NewTextVersion(textBytes)
	assert.NoError(t, err)
	v, err := versioner.Get("")
	assert.NoError(t, err)
	assert.Equal(t, "v0.1.2", v)
}
