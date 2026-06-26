package version

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestYamlVersion(t *testing.T) {
	yamlBytes := []byte(`version: 0.1.2`)
	versioner, err := NewYamlVersion(yamlBytes)
	assert.NoError(t, err)
	v, err := versioner.Get("version")
	assert.NoError(t, err)
	assert.Equal(t, "0.1.2", v)
}

func TestDeeperYamlVersion(t *testing.T) {
	yamlBytes := []byte(`application:
  version: 0.1.2`)
	versioner, err := NewYamlVersion(yamlBytes)
	assert.NoError(t, err)
	v, err := versioner.Get("application.version")
	assert.NoError(t, err)
	assert.Equal(t, "0.1.2", v)
}

func TestYamlVersionNotFound(t *testing.T) {
	yamlBytes := []byte(`application:
  version: 0.1.2`)
	versioner, err := NewYamlVersion(yamlBytes)
	assert.NoError(t, err)
	_, err = versioner.Get("application.test.code")
	assert.Error(t, err)
}
