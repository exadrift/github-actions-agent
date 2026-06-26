package git

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/exadrift/github-actions-agent/internal/utils"
)

type GitVersionInterface interface {
	GetShortCommitHash() (string, error)
	GetCurrentBranch() (string, error)
}

type GitChangesInterface interface {
	GetRepoRootDir() (string, error)
	GetDefaultBranch() (string, error)
	GetChangedFiles(parentBranch string, ignoreDeleted bool) ([]string, error)
}

type GitTagsInterface interface {
	GetCommitish(tag string) (string, error)
}

type GitReleaseInterface interface {
	GetCommitHash() (string, error)
	GetCommitish(tag string) (string, error)
}

type Git struct {
	repoPath string
}

func NewGit(repoPath string) *Git {
	return &Git{
		repoPath,
	}
}

// GetShortCommitHash returns the short commit hash of the current HEAD
func (g *Git) GetShortCommitHash() (string, error) {
	result, err := utils.RunIt("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "", fmt.Errorf("unable to obtain short commit hash: %w", err)
	}
	return result, nil
}

// GetCommitHash returns the short commit hash of the current HEAD
func (g *Git) GetCommitHash() (string, error) {
	result, err := utils.RunIt("git", "rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf("unable to obtain commit hash: %w", err)
	}
	return result, nil
}

// GetCurrentBranch returns the current branch
func (g *Git) GetCurrentBranch() (string, error) {
	result, err := utils.RunIt("git", "branch", "--show-current")
	if err != nil {
		return "", fmt.Errorf("unable to obtain current branch: %w", err)
	}
	return fmt.Sprintf("remotes/origin/%s", result), nil
}

// GetDefaultBranch returns the default branch (main/master, etc.)
func (g *Git) GetDefaultBranch() (string, error) {
	result, err := utils.RunIt("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	if err != nil {
		return "", fmt.Errorf("unable to obtain default branch: %w", err)
	}
	parts := strings.Split(result, "/")
	return fmt.Sprintf("remotes/origin/%s", strings.Trim(parts[len(parts)-1], "\r\n")), nil
}

// GetChangedFiles returns a list of files containing changes, relative to the repository location
func (g *Git) GetChangedFiles(parentBranch string, ignoreDeleted bool) ([]string, error) {
	args := []string{
		"diff", "--name-only",
	}

	if ignoreDeleted {
		args = append(args, "--diff-filter=d")
	}

	args = append(args, parentBranch)

	result, err := utils.RunIt("git", args...)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain branch changes: %w", err)
	}

	return strings.Split(result, "\n"), nil
}

// GetRepoRootDir returns the repository root path
func (g *Git) GetRepoRootDir() (string, error) {
	result, err := utils.RunIt("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("unable to obtain repository root path: %w", err)
	}
	return result, nil
}

// GetCommitish returns a lower-cased hex string of the commit hash, or an empty string if the commitish does not exist.
// an error is returned if anything else fails
func (g *Git) GetCommitish(tag string) (string, error) {
	_, err := utils.RunIt("git", "rev-parse", "--git-dir")
	if err != nil {
		return "", fmt.Errorf("unable to determine if this is a git directory: %w", err)
	}

	result, err := utils.RunIt("git", "rev-parse", "--verify", tag)
	if err != nil {
		// we return a nil value here to indicate that the revision does not exist
		return "", nil
	}

	_, err = hex.DecodeString(result)
	if err != nil {
		return "", nil
	}

	return result, nil
}
