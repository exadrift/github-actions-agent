package branch

import (
	"fmt"
	"os"

	"github.com/exadrift/github-actions-agent/internal/git"
)

type BranchCmd struct {
	Default bool `short:"d" help:"print the default branch instead of the active branch"`
}

func (c *BranchCmd) Run() error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	gitIfx := git.NewGit(workingDirectory)
	var branch string
	if c.Default {
		branch, err = gitIfx.GetDefaultBranch()
	} else {
		branch, err = gitIfx.GetCurrentBranch()
	}
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", branch)
	return nil
}
