package repo

import (
	"fmt"

	"github.com/gitsight/go-vcsurl"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
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

func DetermineChanges(dir string, oldCommit string, newCommit string) ([]string, error) {
	// Load the repo.
	repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't instantiate repo: %w", err)
	}

	// Do a bunch of plumbing to get these commits usable for go-git
	oldHash, err := repo.ResolveRevision(plumbing.Revision(oldCommit))
	if err != nil {
		return nil, fmt.Errorf("couldn't parse old commit: %w", err)
	}
	newHash, err := repo.ResolveRevision(plumbing.Revision(newCommit))
	if err != nil {
		return nil, fmt.Errorf("couldn't parse new commit: %w", err)
	}
	old, err := repo.CommitObject(*oldHash)
	if err != nil {
		return nil, fmt.Errorf("couldn't find old commit: %w", err)
	}
	new, err := repo.CommitObject(*newHash)
	if err != nil {
		return nil, fmt.Errorf("couldn't find new commit: %w", err)
	}

	// Validate that old is before the new commit.
	// TODO

	// Get commits to iterate over.
	commits, err := repo.Log(&git.LogOptions{
		From: new.Hash,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't get git log: %w", err)
	}
	defer commits.Close()

	// Discover the set of changed files between the two commits.
	changed := map[string]bool{}
	err = commits.ForEach(func(c *object.Commit) error {
		if c.Hash == old.Hash {
			return storer.ErrStop
		}

		stats, err := c.Stats()
		if err != nil {
			return err
		}

		for _, stat := range stats {
			changed[stat.Name] = true
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't iterate over commits: %w", err)
	}

	changedFiles := []string{}
	for item := range changed {
		changedFiles = append(changedFiles, item)
	}

	return changedFiles, nil
}
