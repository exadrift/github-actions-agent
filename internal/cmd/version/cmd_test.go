package version

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

type PhonyGit struct {
	ShortHash string
	Branch    string
}

func (g *PhonyGit) GetShortCommitHash() (string, error) {
	return g.ShortHash, nil
}

// GetCurrentBranch returns the current branch
func (g *PhonyGit) GetCurrentBranch() (string, error) {
	return g.Branch, nil
}

func TestInferFileType(t *testing.T) {
	assert.Equal(t, FileTypeYaml, inferFileType("my/path.file/test.yaml"))
}

func TestEnsureVersionProperlyInferred(t *testing.T) {
	v, err := GetVersion(&VersionGetCmd{
		Location: "./version.yaml",
		Key:      "application.version",
	}, &PhonyGit{
		ShortHash: "abcdef",
		Branch:    "test",
	},
		[]byte(`application:
  version: 1.2.4
  `))
	assert.NoError(t, err)
	assert.Equal(t, "1.2.4", v)
}

func TestEnsureVersionProperlyInferredOnTextFile(t *testing.T) {
	v, err := GetVersion(&VersionGetCmd{
		Location: "./VERSION",
		Key:      "application.version",
	}, &PhonyGit{
		ShortHash: "abcdef",
		Branch:    "test",
	},
		[]byte(`1.2.4
  `))
	assert.NoError(t, err)
	assert.Equal(t, "1.2.4", v)
}

func TextExtractMajorVersion(t *testing.T) {
	assert.Equal(t, "v1", extractMajorVersion("v1.0.2"))
	assert.Equal(t, "v1", extractMajorVersion("v1"))
	assert.Equal(t, "v1", extractMajorVersion("v1.0"))
}

func TextExtractMajorMinorVersion(t *testing.T) {
	assert.Equal(t, "v1.0", extractMajorMinorVersion("v1.0.2"))
	assert.Equal(t, "v1.0", extractMajorMinorVersion("v1"))
	assert.Equal(t, "v1.1", extractMajorMinorVersion("v1.1"))
	assert.Equal(t, "1.2", extractMajorMinorVersion("1.2.3"))
}

func TestFurnishShortCommitHash(t *testing.T) {
	assert.Equal(t, "v1.0.0-abcdef", furnishShortCommitHash("v1.0.0", "abcdef"))
}

func TestFurnishNpmShortCommitHash(t *testing.T) {
	assert.Equal(t, "v1.0.0-dev.abcdef", furnishNpmShortCommitHash("v1.0.0", "abcdef"))
}

func TestFurnishPythonDevTimestamp(t *testing.T) {
	assert.Equal(t, "v1.0.0.dev123456", furnishPoetryDevTimestamp("v1.0.0", int64(123456)))
}

func TestGetVersionShortHash(t *testing.T) {
	version, err := GetVersion(&VersionGetCmd{
		Location:              "test.yaml",
		Key:                   "application.version",
		AppendShortCommitHash: true,
	}, &PhonyGit{
		ShortHash: "abcdef",
		Branch:    "test",
	}, []byte(`application:
  version: 1.2.4
  `))
	assert.NoError(t, err)
	assert.Equal(t, "1.2.4-abcdef", version)
}

func TestGetVersionBranch(t *testing.T) {
	version, err := GetVersion(&VersionGetCmd{
		Location:     "test.yaml",
		Key:          "application.version",
		AppendBranch: true,
	}, &PhonyGit{
		ShortHash: "abcdef",
		Branch:    "test",
	}, []byte(`application:
  version: 1.2.4
  `))
	assert.NoError(t, err)
	assert.Equal(t, "1.2.4-test", version)
}
