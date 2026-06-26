package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/cli/go-gh/v2/pkg/api"
)

type Ref struct {
	Ref string `json:"ref"`
	Sha string `json:"sha"`
}

type Sha struct {
	Sha string `json:"sha"`
}

type RefObject struct {
	Type string `json:"type"`
	Sha  string `json:"sha"`
	Url  string `json:"url"`
}

type RefResponse struct {
	Ref    string    `json:"ref"`
	NodeId string    `json:"node_id"`
	Url    string    `json:"url"`
	Object RefObject `json:"object"`
}

type ReleaseResponse struct {
	Id      int    `json:"id"`
	TagName string `json:"tag_name"`
}

type ReleasePost struct {
	Tag  string `json:"tag_name"`
	Name string `json:"name"`
	Body string `json:"body"`
}

type ReleaseAssetResponse struct {
	Id                 int    `json:"id"`
	BrowserDownloadUrl string `json:"browser_download_url"`
	Name               string `json:"name"`
	Label              string `json:"label"`
}

func getReader(obj any) (*bytes.Reader, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

type GithubInterface interface {
}

type GithubApi struct {
	client   *api.RESTClient
	repoName string
	owner    string
}

type GithubApiInterface interface {
	CreateTag(ref string, sha string) (*RefResponse, error)
	UpdateTag(ref string, sha string) (*RefResponse, error)
	CreateRelease(tag string, description string) (*ReleaseResponse, error)
	UploadReleaseAsset(releaseId int, name string, reader io.Reader) (*ReleaseAssetResponse, error)
}

func NewGithubApi(repoName string, owner string) (*GithubApi, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &GithubApi{
		client:   client,
		repoName: repoName,
		owner:    owner,
	}, nil
}

func (a *GithubApi) GetEndpoint(suffix string) string {
	return fmt.Sprintf("repos/%s/%s%s", a.owner, a.repoName, suffix)
}

func (a *GithubApi) CreateTag(ref string, sha string) (*RefResponse, error) {
	fullRef := fmt.Sprintf("refs/tags/%s", ref)
	reader, err := getReader(&Ref{Ref: fullRef, Sha: sha})
	if err != nil {
		return nil, err
	}
	resp := &RefResponse{}
	if err = a.client.Post(a.GetEndpoint("/git/refs"), reader, resp); err != nil {
		return nil, fmt.Errorf("unable to post reference \"%s\": %w", fullRef, err)
	}

	return resp, nil
}

func (a *GithubApi) UpdateTag(ref string, sha string) (*RefResponse, error) {
	reader, err := getReader(&Sha{Sha: sha})
	if err != nil {
		return nil, err
	}
	resp := &RefResponse{}
	if err = a.client.Patch(a.GetEndpoint(fmt.Sprintf("/git/refs/tags/%s", ref)), reader, resp); err != nil {
		return nil, fmt.Errorf("unable to patch reference: %w", err)
	}

	return resp, nil
}

func (a *GithubApi) CreateRelease(tag string, description string) (*ReleaseResponse, error) {
	reader, err := getReader(&ReleasePost{Tag: tag, Name: tag, Body: description})
	if err != nil {
		return nil, err
	}

	resp := &ReleaseResponse{}
	if err := a.client.Post(a.GetEndpoint("/releases"), reader, resp); err != nil {
		return nil, fmt.Errorf("unable to post release: %w", err)
	}

	return resp, nil
}

func (a *GithubApi) UploadReleaseAsset(releaseId int, name string, reader io.Reader) (*ReleaseAssetResponse, error) {
	resp := &ReleaseAssetResponse{}
	url := a.GetEndpoint(fmt.Sprintf("/releases/%d/assets?name=%s", releaseId, name))
	// when a full URL is provided to the client as the path, the URL is used verbatim instead of being composed internally
	if err := a.client.Post(fmt.Sprintf("https://uploads.github.com/%s", url), reader, resp); err != nil {
		return nil, fmt.Errorf("unable to upload release asset: %w", err)
	}

	return resp, nil
}
