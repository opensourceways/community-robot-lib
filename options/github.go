package options

import (
	"flag"
)

// GithubOptions holds options for interacting with Github.
type GithubOptions struct {
	TokenPath     string
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
}

// Validate validates Gitee options.
func (o GithubOptions) Validate() error {
	return nil
}
