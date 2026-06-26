package version

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestTomlVersion(t *testing.T) {
	tomlBytes := []byte(`version = '0.1.2'`)
	versioner, err := NewTomlVersion(tomlBytes)
	assert.NoError(t, err)
	v, err := versioner.Get("version")
	assert.NoError(t, err)
	assert.Equal(t, "0.1.2", v)
}

func TestDeeperTomlVersion(t *testing.T) {
	tomlBytes := []byte(`[application]
version = '0.1.2'`)
	versioner, err := NewTomlVersion(tomlBytes)
	assert.NoError(t, err)
	v, err := versioner.Get("application.version")
	assert.NoError(t, err)
	assert.Equal(t, "0.1.2", v)
}

func TestEvenDeeperTomlVersion(t *testing.T) {
	tomlBytes := []byte(`[application.package]
version = '0.1.2'`)
	versioner, err := NewTomlVersion(tomlBytes)
	assert.NoError(t, err)
	v, err := versioner.Get("application.package.version")
	assert.NoError(t, err)
	assert.Equal(t, "0.1.2", v)
}

func TestTomlVersionNotFound(t *testing.T) {
	tomlBytes := []byte(`[application]
version = '0.1.2'`)
	versioner, err := NewTomlVersion(tomlBytes)
	assert.NoError(t, err)
	_, err = versioner.Get("application.test.code")
	assert.Error(t, err)
}
