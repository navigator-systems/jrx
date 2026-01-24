package scm

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/go-github/v58/github"
	"github.com/navigator-systems/jrx/internal/config"
	"golang.org/x/oauth2"
)

// GitHubClient wraps the GitHub client with configuration
type GitHubClient struct {
	client *github.Client
	org    string
}

// NewGitHubClient creates a new GitHub client from JRX config
func NewGitHubClient(cfg config.JRXConfig, githubOrg string) (*GitHubClient, error) {
	if cfg.GitProvider.GithubToken == "" {
		return nil, fmt.Errorf("github_token not found in config")
	}

	if slices.Contains(cfg.GitProvider.GithubOrganization, githubOrg) {
		// Ok
	} else {
		return nil, fmt.Errorf("github_organization %s not found in config", githubOrg)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitProvider.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	var client *github.Client
	var err error
	domain := cfg.GitProvider.GithubURL
	// If it's not github.com, use Enterprise client
	if domain != "" && domain != "github.com" {
		baseURL := fmt.Sprintf("https://%s/api/v3/", domain)
		client, err = github.NewClient(tc).WithEnterpriseURLs(baseURL, baseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create enterprise client: %w", err)
		}
	} else {
		client = github.NewClient(tc)
	}

	return &GitHubClient{
		client: client,
		org:    githubOrg,
	}, nil
}

// CreateRepository creates a new repository in the organization
func (gc *GitHubClient) CreateRepository(ctx context.Context, repoName string, description string, private bool) (*github.Repository, error) {
	repo := &github.Repository{
		Name:        github.String(repoName),
		Description: github.String(description),
		Private:     github.Bool(private),
		AutoInit:    github.Bool(false), // Don't initialize with README yet
	}

	createdRepo, _, err := gc.client.Repositories.Create(ctx, gc.org, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	return createdRepo, nil
}
