package version

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestJsonVersion(t *testing.T) {
	jsonBytes := []byte(`{"version":"0.1.2"}`)
	versioner, err := NewJsonVersion(jsonBytes)
	assert.NoError(t, err)
	v, err := versioner.Get("version")
	assert.NoError(t, err)
	assert.Equal(t, "0.1.2", v)
}

func TestDeeperJsonVersion(t *testing.T) {
	jsonBytes := []byte(`{"application": {"version":"0.1.2"}}`)
	versioner, err := NewJsonVersion(jsonBytes)
	assert.NoError(t, err)
	v, err := versioner.Get("application.version")
	assert.NoError(t, err)
	assert.Equal(t, "0.1.2", v)
}

func TestJsonVersionNotFound(t *testing.T) {
	jsonBytes := []byte(`{"application": {"version":"0.1.2"}}`)
	versioner, err := NewJsonVersion(jsonBytes)
	assert.NoError(t, err)
	_, err = versioner.Get("application.test.code")
	assert.Error(t, err)
}
