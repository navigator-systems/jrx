package templates

import (
	"path/filepath"
	"sort"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/navigator-systems/jrx/internal/errors"
)

// GetVersionsTags fetches tags from remote repository, filters by pattern, sorts by most recent, and limits results
func (tm *TemplateManager) GetVersionsTags() ([]string, error) {
	publicKeys, err := ssh.NewPublicKeysFromFile("git", tm.config.SshKeyPath, tm.config.SshKeyPassphrase)
	if err != nil {
		return nil, errors.NewError("create SSH keys", err)
	}

	rem := git.NewRemote(nil, &config.RemoteConfig{
		Name: "origin",
		URLs: []string{tm.config.TemplatesRepo},
	})

	refs, err := rem.List(&git.ListOptions{Auth: publicKeys})
	if err != nil {
		return nil, errors.NewError("list remote references", err)
	}

	var tagVersions []string
	pattern := tm.config.TemplatesPatternGlob

	// Extract and filter tags
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tagName := ref.Name().Short()

			// Filter by pattern if specified
			if pattern != "" {
				matched, err := filepath.Match(pattern, tagName)
				if err != nil {
					return nil, errors.NewError("match pattern", err)
				}
				if matched {
					tagVersions = append(tagVersions, tagName)
				}
			} else {
				tagVersions = append(tagVersions, tagName)
			}
		}
	}

	// Sort tags in descending order (most recent first)
	// This is a simple lexicographic sort - for semantic versioning, consider using a semver library
	sort.Sort(sort.Reverse(sort.StringSlice(tagVersions)))

	// Limit to max versions if specified
	if tm.config.TemplatesMaxVersions > 0 && len(tagVersions) > tm.config.TemplatesMaxVersions {
		tagVersions = tagVersions[:tm.config.TemplatesMaxVersions]
	}

	return tagVersions, nil
}

// GetAvailableVersions returns a list of all available template versions
// by combining branches and tags from config
func (tm *TemplateManager) GetAvailableVersions() []string {
	versions := make([]string, 0)

	// Add all configured branches
	versions = append(versions, tm.config.TemplatesBranch...)

	// Add all configured tags
	tags, err := tm.GetVersionsTags()
	if err == nil {
		versions = append(versions, tags...)
	}

	return versions
}

// ValidateVersion checks if a given version exists in the available versions
func (tm *TemplateManager) ValidateVersion(version string) bool {
	if version == "" {
		return true // Empty version is valid (uses default)
	}

	availableVersions := tm.GetAvailableVersions()
	for _, v := range availableVersions {
		if v == version {
			return true
		}
	}
	return false
}

// GetCurrentVersion returns the currently loaded version
func (tm *TemplateManager) GetCurrentVersion() string {
	return tm.currentVersion
}

// ListAll returns all available templates
func (tm *TemplateManager) ListAll() ([]RootTemplate, error) {
	if !tm.loaded {
		return nil, errors.NewError("list templates", errors.ErrLoadTemplates)
	}

	templates := make([]RootTemplate, 0, len(tm.templateFile.Templates))
	for _, tpl := range tm.templateFile.Templates {
		templates = append(templates, tpl)
	}

	return templates, nil
}
