package options

import (
	"flag"
)

// GitLabOptions holds options for interacting with GitLab.
type GitLabOptions struct {
	TokenPath     string
	RepoCacheDir  string
	CacheRepoOnPV bool
}

// NewGitLabOptions creates a GiteeOptions with default values.
func NewGitLabOptions() *GiteeOptions {
	return &GiteeOptions{}
}

// AddFlags injects Gitee options into the given FlagSet.
func (o *GitLabOptions) AddFlags(fs *flag.FlagSet) {
	o.addFlags("/etc/gitee/oauth", fs)
}

// AddFlagsWithoutDefaultGitLabTokenPath injects Gitlab options into the given
// FlagSet without setting a default for the gitlabTokenPath, allowing to
// use an anonymous Gitlab client
func (o *GitLabOptions) AddFlagsWithoutDefaultGitLabTokenPath(fs *flag.FlagSet) {
	o.addFlags("", fs)
}

func (o *GitLabOptions) addFlags(defaultGitlabTokenPath string, fs *flag.FlagSet) {
	fs.StringVar(
		&o.TokenPath,
		"gitlab-token-path",
		defaultGitlabTokenPath,
		"Path to the file containing the GitLab OAuth secret.",
	)
	fs.StringVar(&o.RepoCacheDir, "repo-cache-dir", "", "Path to which clone repo.")
	fs.BoolVar(&o.CacheRepoOnPV, "cache-repo-on-pv", false, "Specify whether to cache repo on persistent volume.")
}

// Validate validates Gitlab options.
func (o GitLabOptions) Validate() error {
	return nil
}
