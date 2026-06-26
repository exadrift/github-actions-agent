package tag

import (
	"fmt"
	"os"

	"github.com/exadrift/github-actions-agent/internal/git"
	"github.com/exadrift/github-actions-agent/internal/github"
)

type VerifyBranchCmd struct {
	Branch string `arg:"" help:"branch for which to return commit hash"`
}

func (c *VerifyBranchCmd) Run() error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	// should contain remotes/origins/<branch>

	gitIfx := git.NewGit(workingDirectory)
	commitHash, err := gitIfx.GetCommitish(c.Branch)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", commitHash)
	return nil
}

type TagCmd struct {
	Tag string `arg:"" help:"tag to check for existance"`
}

func (c *TagCmd) Run() error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	gitIfx := git.NewGit(workingDirectory)
	commitHash, err := gitIfx.GetCommitish(fmt.Sprintf("refs/tags/%s", c.Tag))
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", commitHash)
	return nil
}

type TagCommitCmd struct {
	Tag   string `arg:"" help:"name of the tag to create or update"`
	Sha   string `arg:"" help:"the commit hash to attribute to the tag"`
	Owner string `short:"o" help:"repository owner"`
	Repo  string `short:"r" help:"repository name"`
}

func (c *TagCommitCmd) Run() error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	gitIfx := git.NewGit(workingDirectory)
	commitHash, err := gitIfx.GetCommitish(fmt.Sprintf("refs/tags/%s", c.Tag))
	if err != nil {
		return err
	}
	apiIfx, err := github.NewGithubApi(c.Repo, c.Owner)
	if err != nil {
		return err
	}

	if commitHash == "" {
		// tag does not exist, create it
		_, err := apiIfx.CreateTag(c.Tag, c.Sha)
		if err != nil {
			return fmt.Errorf("unable to create ref: %w", err)
		}

		return nil
	}

	// tag exists, update it
	if _, err = apiIfx.UpdateTag(c.Tag, c.Sha); err != nil {
		return fmt.Errorf("unable to update ref: %w", err)
	}

	return nil
}
