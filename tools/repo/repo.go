package repo

import (
	"fmt"

	"github.com/gitsight/go-vcsurl"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
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
	oldTree, err := old.Tree()
	if err != nil {
		return nil, fmt.Errorf("couldn't generate tree for old commit: %w", err)
	}
	newTree, err := new.Tree()
	if err != nil {
		return nil, fmt.Errorf("couldn't generate tree for new commit: %w", err)
	}

	// Diff the two trees to determine what changed.
	changes, err := oldTree.Diff(newTree)
	if err != nil {
		return nil, fmt.Errorf("couldn't get changes between commits: %w", err)
	}

	// Create a unique list of changed files.
	changedFiles := []string{}
	for _, change := range changes {
		action, err := change.Action()
		if err != nil {
			return nil, fmt.Errorf("couldn't determine action for commit: %w", err)
		}

		// Skip deleted files.
		if action == merkletrie.Delete {
			continue
		}

		changedFiles = append(changedFiles, getChangeName(change))
	}

	return changedFiles, nil
}

func getChangeName(change *object.Change) string {
	var empty = object.ChangeEntry{}
	if change.From != empty {
		return change.From.Name
	}

	return change.To.Name
}
