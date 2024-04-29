package repo

import (
	"fmt"

	"github.com/gitsight/go-vcsurl"
	"github.com/go-git/go-git/v5"
)

// IsInRepo determines if the provided directory is in the provided git
// repository. Git repositories can be named differently on a local clone then
// the remote. In addition, a git repo can have multiple remotes. In practice
// though, the business logic question is something like:
// "Am I in the community repo?". To answer that, this function iterates over
// the remotes and if any of them have the same name as the one requested, it
// returns true. Any other case returns false.
func IsInRepo(dir string, name string) bool {
	repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return false
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return false
	}

	for _, remote := range remotes {
		for _, url := range remote.Config().URLs {
			info, err := vcsurl.Parse(url)
			if err != nil {
				return false
			}

			if info.Name == name {
				return true
			}
		}
	}

	return false
}

func RepoRoot(dir string) (string, error) {
	repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return "", fmt.Errorf("couldn't instantiate repo: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("couldn't get worktree: %w", err)
	}

	return worktree.Filesystem.Root(), nil
}
