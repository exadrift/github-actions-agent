package changes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/exadrift/github-actions-agent/internal/git"
)

type ChangesCmd struct {
	ParentCommitish          string `short:"p" help:"the parent commitish against which to compare the current branch"`
	ExcludePrefix            string `short:"e" help:"exclude any change that matches this prefix list from the repository root (comma separated)"`
	IncludePrefix            string `short:"i" help:"include any change that matches this prefix list from the repository root (comma separated)"`
	IncludeDeleted           string `short:"d" help:"ignore changes which are deletions (true or false) defaults to false"`
	Json                     bool   `short:"j" help:"output changes as json list"`
	IncludeHiddenDirectories string `short:"z" help:"include hidden directories (true or false) defaults to false"`
}

type ChangedDir struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func (c *ChangesCmd) Run() error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	gitIfx := git.NewGit(workingDirectory)
	if c.ParentCommitish == "" {
		defBranch, err := gitIfx.GetDefaultBranch()
		if err != nil {
			return err
		}
		c.ParentCommitish = defBranch
	}
	changedPaths, err := getChangedDirs(c, gitIfx, c.ParentCommitish)
	if err != nil {
		return err
	}

	if c.Json {
		jsonPaths := []ChangedDir{}
		for _, path := range changedPaths {
			_, dirName := filepath.Split(path)
			jsonPaths = append(jsonPaths, ChangedDir{
				Path: path,
				Name: dirName,
			})
		}
		pathBytes, err := json.Marshal(jsonPaths)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", string(pathBytes))
	} else {
		for _, path := range changedPaths {
			fmt.Printf("%s\n", path)
		}
	}

	return nil
}

func getChangedDirs(c *ChangesCmd, gitChangesIfx git.GitChangesInterface, parentCommitish string) ([]string, error) {
	changedFiles, err := gitChangesIfx.GetChangedFiles(parentCommitish, strings.ToLower(c.IncludeDeleted) == "true")
	if err != nil {
		return nil, err
	}

	repoRoot, err := gitChangesIfx.GetRepoRootDir()
	if err != nil {
		return nil, err
	}

	includeHiddenDirs := strings.ToLower(c.IncludeHiddenDirectories) == "true"
	dirs := map[string]struct{}{}
	skippedDirs := map[string]struct{}{}
	for _, changedFile := range changedFiles {
		// change to a full path first
		changedFile = filepath.Join(repoRoot, changedFile)
		changedFile = strings.TrimPrefix(changedFile, repoRoot)
		changedFile = strings.TrimPrefix(changedFile, "/")
		parts := strings.Split(changedFile, "/")
		if len(parts) < 2 {
			continue
		}
		firstDir := parts[0]
		if _, ok := dirs[firstDir]; ok {
			continue
		}
		if _, ok := skippedDirs[firstDir]; ok {
			continue
		}

		if len(c.IncludePrefix) > 0 {
			contains := false
			for _, match := range strings.Split(c.IncludePrefix, ",") {
				if strings.HasPrefix(changedFile, match) {
					contains = true
					break
				}
			}

			if !contains {
				skippedDirs[firstDir] = struct{}{}
				continue
			}
		}

		if len(c.ExcludePrefix) > 0 {
			contains := false
			for _, match := range strings.Split(c.ExcludePrefix, ",") {
				if strings.HasPrefix(changedFile, match) {
					contains = true
					break
				}
			}
			if contains {
				skippedDirs[firstDir] = struct{}{}
				continue
			}
		}

		if !includeHiddenDirs && strings.HasPrefix(firstDir, ".") {
			continue
		}
		dirs[firstDir] = struct{}{}
	}

	var dirList []string
	for dir := range dirs {
		dirList = append(dirList, filepath.Join(repoRoot, dir))
	}

	sort.Slice(dirList, func(i, j int) bool {
		return dirList[i] < dirList[j]
	})
	return dirList, nil
}
