package community

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/tools/repo"
)

var (
	oldCommit string
	newCommit string
)

func init() {
	TargetDeterminatorCmd.Flags().StringVarP(&oldCommit, "old", "o", "", "The old commit to compare against")
	TargetDeterminatorCmd.Flags().StringVarP(&newCommit, "new", "n", "", "The new commit to compare against")
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

	files, err := repo.DetermineChanges(cwd, oldCommit, newCommit)
	if err != nil {
		return fmt.Errorf("could not determine targets: %w", err)
	}

	for _, f := range files {
		fmt.Println(f)
	}

	return nil
}
