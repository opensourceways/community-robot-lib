package options

import (
	"flag"
	"fmt"
)

// GithubOptions holds options for interacting with Gitee.
type GithubOptions struct {
	TokenPath     string
	RepoCacheDir  string
	CacheRepoOnPV bool
}

// NewGithubOptions creates a GiteeOptions with default values.
func NewGithubOptions() *GiteeOptions {
	return &GiteeOptions{}
}

// AddFlags injects Gitee options into the given FlagSet.
func (o *GithubOptions) AddFlags(fs *flag.FlagSet) {
	o.addFlags("/etc/gitee/oauth", fs)
}

// AddFlagsWithoutDefaultGithubTokenPath injects Gitee options into the given
// Flagset without setting a default for for the giteeTokenPath, allowing to
// use an anonymous Gitee client
func (o *GithubOptions) AddFlagsWithoutDefaultGithubTokenPath(fs *flag.FlagSet) {
	o.addFlags("", fs)
}

func (o *GithubOptions) addFlags(defaultGithubTokenPath string, fs *flag.FlagSet) {
	fs.StringVar(
		&o.TokenPath,
		"github-token-path",
		defaultGithubTokenPath,
		"Path to the file containing the Gitee OAuth secret.",
	)
	fs.StringVar(&o.RepoCacheDir, "repo-cache-dir", "", "Path to which clone repo.")
	fs.BoolVar(&o.CacheRepoOnPV, "cache-repo-on-pv", false, "Specify whether to cache repo on persistent volume.")
}

// Validate validates Gitee options.
func (o GithubOptions) Validate() error {
	if o.CacheRepoOnPV && o.RepoCacheDir == "" {
		return fmt.Errorf("must set repo-cache-dir if caching repo on persistent volume")
	}

	return nil
}
