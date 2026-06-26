package changes

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

var files = []string{
	"database/README.md",
	"database/VERSION",
	"database/CHANGELOG.md",
	"redmist/tests/new.go",
	"redmist/VERSION",
	"blue/abcdef.thing",
	"blue/yefdbz.sum",
	"test.yaml",
	".gitignore",
	"/.somedir/test",
}

type PhonyGit struct {
	changes []string
}

func NewPhonyGit(changes []string) *PhonyGit {
	return &PhonyGit{
		changes,
	}
}

func (g *PhonyGit) GetRepoRootDir() (string, error) {
	return "/home/user/repo", nil
}

func (g *PhonyGit) GetDefaultBranch() (string, error) {
	return "main", nil
}

func (g *PhonyGit) GetChangedFiles(parentBranch string, ignoreDeleted bool) ([]string, error) {
	return g.changes, nil
}

func TestGetChanges(t *testing.T) {
	changedDirs, err := getChangedDirs(
		&ChangesCmd{},
		NewPhonyGit(files),
		"main",
	)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"/home/user/repo/blue",
		"/home/user/repo/database",
		"/home/user/repo/redmist",
	}, changedDirs)
}

func TestGetChangesExcludeMatch(t *testing.T) {
	changedDirs, err := getChangedDirs(
		&ChangesCmd{
			ExcludePrefix: "redmist/",
		},
		NewPhonyGit(files),
		"main",
	)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"/home/user/repo/blue",
		"/home/user/repo/database",
	}, changedDirs)
}

func TestGetChangesIncludeMatch(t *testing.T) {
	changedDirs, err := getChangedDirs(
		&ChangesCmd{
			IncludePrefix: "redmist/",
		},
		NewPhonyGit(files),
		"main",
	)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"/home/user/repo/redmist",
	}, changedDirs)
}

func TestGetChangesIncludeHidden(t *testing.T) {
	changedDirs, err := getChangedDirs(
		&ChangesCmd{
			IncludeHiddenDirectories: "true",
		},
		NewPhonyGit(files),
		"main",
	)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"/home/user/repo/.somedir",
		"/home/user/repo/blue",
		"/home/user/repo/database",
		"/home/user/repo/redmist",
	}, changedDirs)
}
