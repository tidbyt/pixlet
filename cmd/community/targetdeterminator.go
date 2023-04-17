package community

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/tools/repo"
)

var (
	oldCommit string
	newCommit string
)

func init() {
	TargetDeterminatorCmd.Flags().StringVarP(&oldCommit, "old", "o", "", "The old commit to compare against")
	TargetDeterminatorCmd.MarkFlagRequired("old")
	TargetDeterminatorCmd.Flags().StringVarP(&newCommit, "new", "n", "", "The new commit to compare against")
	TargetDeterminatorCmd.MarkFlagRequired("new")
}

var TargetDeterminatorCmd = &cobra.Command{
	Use:   "target-determinator",
	Short: "Determines what files have changed between old and new commit",
	Example: `  pixlet community target-determinator \
    --old 4d69e9bbf181434229a98e87909a619634072930 \
    --new 2fc2a1fcfa48bbb0b836084c1b1e259322c4e133`,
	Long: `This command determines what files have changed between two commits
so that we can limit the build in the community repo to only the files that
have changed.`,
	RunE: determineTargets,
}

func determineTargets(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not determine targets, something went wrong with the local filesystem: %w", err)
	}

	changedFiles, err := repo.DetermineChanges(cwd, oldCommit, newCommit)
	if err != nil {
		return fmt.Errorf("could not determine targets: %w", err)
	}

	changedApps := map[string]bool{}
	for _, f := range changedFiles {
		dir := filepath.Dir(f) + string(os.PathSeparator)

		// We only care about things in apps/{{ app package }}/ changing and
		// nothing else. Skip any file changes that are not in that structure.
		parts := strings.Split(f, string(os.PathSeparator))
		if len(parts) < 3 {
			continue
		}

		// If the filepath does not start with apps, we also don't care about
		// it.
		if !strings.HasPrefix(dir, "apps") {
			continue
		}

		// If the directory no longer exists, we don't care about it. This would
		// happen if someone deleted an app in a PR.
		if !dirExists(dir) {
			continue
		}

		changedApps[dir] = true
	}

	for app := range changedApps {
		fmt.Println(app)
	}

	return nil
}

func dirExists(dir string) bool {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}

	if err != nil {
		return false
	}

	return true
}
