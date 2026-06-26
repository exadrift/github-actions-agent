package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/exadrift/github-actions-agent/internal/cmd/branch"
	"github.com/exadrift/github-actions-agent/internal/cmd/changelog"
	"github.com/exadrift/github-actions-agent/internal/cmd/changes"
	"github.com/exadrift/github-actions-agent/internal/cmd/release"
	"github.com/exadrift/github-actions-agent/internal/cmd/tag"
	"github.com/exadrift/github-actions-agent/internal/cmd/version"
)

// VERSION gets set at build time
var Version = "unknown"

type VersionCmd struct {
}

func (c *VersionCmd) Run() error {
	fmt.Printf("%s\n", Version)
	return nil
}

var CLI struct {
	VersionGet         version.VersionGetCmd  `cmd:"" help:"establish version from version file"`
	VersionSet         version.VersionSetCmd  `cmd:"" help:"update version file with new version"`
	ChangedDirectories changes.ChangesCmd     `cmd:"" help:"return a list of directories containing changes"`
	VerifyChangeLog    changelog.ChangeLogCmd `cmd:"" help:"verify the changelog format and version"`
	VerifyTag          tag.TagCmd             `cmd:"" help:"verify a given tag exists and print the commit hash or an empty string"`
	VerifyBranch       tag.VerifyBranchCmd    `cmd:"" help:"verify a given tag exists and print the commit hash or an empty string"`
	Branch             branch.BranchCmd       `cmd:"" help:"print branch name"`
	Release            release.ReleaseCmd     `cmd:"" help:"create a release (including tags)"`
	WorkingDirectory   string                 `short:"w" help:"change working directory for the command"`
	TagCommit          tag.TagCommitCmd       `cmd:"" help:"tag the current commit point"`
	Version            VersionCmd             `cmd:"" help:"return the application version"`
}

func main() {
	ctx := kong.Parse(&CLI)

	if CLI.WorkingDirectory != "" {
		// change working directory
		err := os.Chdir(CLI.WorkingDirectory)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
	}

	if err := ctx.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
