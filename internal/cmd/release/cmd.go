package release

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/exadrift/github-actions-agent/internal/git"
	"github.com/exadrift/github-actions-agent/internal/github"
)

type ReleaseCmd struct {
	Tags                string `arg:"" help:"comma separated list of tags, only the first tag will receive the binary attachment if specified"`
	BinaryLocations     string `short:"b" optional:"" help:"binary locations for upload, comma separated"`
	BinaryNameOverrides string `optional:"" help:"override name of binary in release, each corresponds to one binary location in order"`
	Description         string `short:"d" help:"release description"`
	Owner               string `short:"o" help:"repository owner"`
	Repo                string `short:"r" help:"repository name"`
}

func (c *ReleaseCmd) Run() error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	gh, err := github.NewGithubApi(c.Repo, c.Owner)
	if err != nil {
		return err
	}

	gitIfx := git.NewGit(workingDirectory)

	var binLocations []string
	if c.BinaryLocations != "" {
		binLocations = strings.Split(c.BinaryLocations, ",")
	}

	var binNameOverrides []string
	if c.BinaryNameOverrides != "" {
		binNameOverrides = strings.Split(c.BinaryNameOverrides, ",")
	} else {
		for _, binLoc := range binLocations {
			_, name := filepath.Split(binLoc)
			binNameOverrides = append(binNameOverrides, name)
		}
	}

	return createReleases(strings.Split(c.Tags, ","), c.Description, binLocations, binNameOverrides, gh, gitIfx)
}

func createReleases(tagNames []string, description string, binaryLocations []string, binaryNames []string, apiIfx github.GithubApiInterface, gitIfx git.GitReleaseInterface) error {
	if len(binaryLocations) != len(binaryNames) {
		return fmt.Errorf("number of binary locations must match number of binary name overrides")
	}
	tags := map[string]struct{}{}
	for _, tag := range tagNames {
		if tag == "" {
			return fmt.Errorf("unable to process empty tag")
		}

		if _, ok := tags[tag]; ok {
			return fmt.Errorf("tag \"%s\" is duplicated in this request", tag)
		}

		tags[tag] = struct{}{}
	}

	for i, tag := range tagNames {
		// check first if the tag exists
		tagCommit, err := gitIfx.GetCommitish(fmt.Sprintf("refs/tags/%s", tag))
		if err != nil {
			return err
		}

		sha, err := gitIfx.GetCommitHash()
		if err != nil {
			return err
		}

		if tagCommit == "" {
			if _, err := apiIfx.CreateTag(tag, sha); err != nil {
				return err
			}

			releaseResp, err := apiIfx.CreateRelease(tag, description)
			if err != nil {
				return err
			}

			if i == 0 {
				for bInd := range binaryLocations {
					binaryName := binaryNames[bInd]
					binaryLocation := binaryLocations[bInd]

					name := binaryName
					if name == "" {
						_, name = filepath.Split(binaryLocation)
					}

					// we must read the entire file in order to use the gh-client REST api.  this may use a lot of RAM but that's fine
					// given the 2GB limit
					rBytes, err := os.ReadFile(binaryLocation)
					if err != nil {
						return err
					}
					if _, err = apiIfx.UploadReleaseAsset(releaseResp.Id, name, bytes.NewReader(rBytes)); err != nil {
						return fmt.Errorf("release id %d upload of %s as %s errored: %w", releaseResp.Id, binaryLocation, name, err)
					}
				}
			}

			continue
		}

		// if this tag already exists, we merely update the existing tag to point to the new commitish
		if _, err = apiIfx.UpdateTag(tag, sha); err != nil {
			return err
		}
	}

	return nil
}
