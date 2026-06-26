package version

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/exadrift/github-actions-agent/internal/git"
)

type VersionSetCmd struct {
	Location string `arg:"" help:"location of the file containing a version string"`
	Version  string `arg:"" help:"value to use as the new version"`
	Key      string `short:"k" help:"key used to locate version within hierachical structure"`
	Type     string `short:"t" help:"file type hint if it cannot be inferred from the extension"`
}

func (c *VersionSetCmd) Run(workDir string) error {
	fileType := fileTypeFromExtension(c.Type)
	if fileType == FileTypeUnknown {
		fileType = inferFileType(c.Location)
	}

	vBytes, err := os.ReadFile(c.Location)
	if err != nil {
		return err
	}

	versioner, err := getVersioner(vBytes, fileType)
	if err != nil {
		return err
	}

	err = versioner.Set(c.Key, c.Version)
	if err != nil {
		return fmt.Errorf("unable to set new version: %w", err)
	}

	vBytes, err = versioner.Marshal()
	if err != nil {
		return fmt.Errorf("unable to marshal version file: %w", err)
	}

	return os.WriteFile(c.Location, vBytes, 0644)
}

type VersionGetCmd struct {
	Location                 string `arg:"" help:"location of the file containing a version string"`
	Key                      string `short:"k" help:"key used to locate version within hierachical structure"`
	Type                     string `short:"t" help:"file type hint if it cannot be inferred from the extension"`
	AppendShortCommitHash    bool   `help:"append the short commit hash to the version"`
	AppendPoetryDevTimestamp bool   `help:"append the timestamp to the version as a dev marker as poetry expects"`
	AppendNpmShortCommitHash bool   `help:"append the short commit hash to the version as npm expects"`
	AppendBranch             bool   `help:"append branch name to the version"`
	ExtractMajorVersion      bool   `help:"extract only the major version component of the verison string"`
	ExtractToMinorVersion    bool   `help:"extract up to and including the minor version component of the verison string"`
}

func (c *VersionGetCmd) Run() error {
	vBytes, err := os.ReadFile(c.Location)
	if err != nil {
		return err
	}

	dir, _ := filepath.Split(c.Location)
	version, err := GetVersion(c, git.NewGit(dir), vBytes)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", version)
	return nil
}

func GetVersion(c *VersionGetCmd, gitIfx git.GitVersionInterface, b []byte) (string, error) {
	fileType := FileTypeUnknown
	if c.Type != "" {
		fileType = fileTypeFromExtension(c.Type)
	}
	if fileType == FileTypeUnknown {
		fileType = inferFileType(c.Location)
	}

	versioner, err := getVersioner(b, fileType)
	if err != nil {
		return "", err
	}

	version, err := versioner.Get(c.Key)
	if err != nil {
		return "", err
	}

	var shortCommitHash string
	if c.AppendShortCommitHash || c.AppendNpmShortCommitHash {
		shortCommitHash, err = gitIfx.GetShortCommitHash()
		if err != nil {
			return "", err
		}
	}

	if c.ExtractMajorVersion {
		version = extractMajorVersion(version)
	} else if c.ExtractToMinorVersion {
		version = extractMajorMinorVersion(version)
	}

	if c.AppendShortCommitHash {
		version = furnishShortCommitHash(version, shortCommitHash)
	} else if c.AppendNpmShortCommitHash {
		version = furnishNpmShortCommitHash(version, shortCommitHash)
	} else if c.AppendPoetryDevTimestamp {
		version = furnishPoetryDevTimestamp(version, time.Now().UTC().Unix())
	} else if c.AppendBranch {
		branch, err := gitIfx.GetCurrentBranch()
		if err != nil {
			return "", err
		}
		branch = strings.TrimPrefix(branch, "remotes/origin/")
		version = fmt.Sprintf("%s-%s", version, branch)
	}

	return version, nil
}

type Versioner interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Marshal() ([]byte, error)
}

type FileType string

func extractMajorVersion(version string) string {
	return strings.Split(version, ".")[0]
}

func extractMajorMinorVersion(version string) string {
	parts := strings.Split(version, ".")
	if len(parts) == 1 {
		return fmt.Sprintf("%s.0", parts[0])
	}

	return fmt.Sprintf("%s.%s", parts[0], parts[1])
}

func furnishNpmShortCommitHash(version string, shortCommitHash string) string {
	return fmt.Sprintf("%s-dev.%s", version, shortCommitHash)
}

func furnishShortCommitHash(version string, shortCommitHash string) string {
	return fmt.Sprintf("%s-%s", version, shortCommitHash)
}

func furnishPoetryDevTimestamp(version string, timestamp int64) string {
	return fmt.Sprintf("%s.dev%d", version, timestamp)
}

func fileTypeFromExtension(ext string) FileType {
	ext = strings.ToLower(ext)
	switch ext {
	case "yaml":
		return FileTypeYaml
	case "yml":
		return FileTypeYaml
	case "json":
		return FileTypeJson
	case "toml":
		return FileTypeToml
	case "txt":
		return FileTypeText
	default:
		return FileTypeUnknown
	}
}

func inferFileType(location string) FileType {
	// first strip anything before a final slash
	parts := strings.Split(location, "/")
	location = parts[len(parts)-1]
	parts = strings.Split(location, ".")
	var ext string
	if len(parts) < 2 {
		ext = ""
	} else {
		ext = parts[len(parts)-1]
	}
	if ext == "" {
		return FileTypeText
	}
	return fileTypeFromExtension(ext)
}

const (
	FileTypeUnknown FileType = ""
	FileTypeYaml    FileType = "YAML"
	FileTypeJson    FileType = "JSON"
	FileTypeToml    FileType = "TOML"
	FileTypeText    FileType = "TEXT"
)

// getVersioner returns a versionion for the provided reader
func getVersioner(b []byte, fileType FileType) (Versioner, error) {
	var versioner Versioner
	var err error
	switch fileType {
	case FileTypeYaml:
		versioner, err = NewYamlVersion(b)
	case FileTypeJson:
		versioner, err = NewJsonVersion(b)
	case FileTypeToml:
		versioner, err = NewTomlVersion(b)
	case FileTypeText:
		versioner, err = NewTextVersion(b)
	default:
		return nil, fmt.Errorf("unknown file type")
	}

	if err != nil {
		return nil, fmt.Errorf("unable to create versioner: %w", err)
	}

	return versioner, nil
}
