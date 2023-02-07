package community

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/icons"
)

var ListIconsCmd = &cobra.Command{
	Use:     "list-icons",
	Short:   "List icons that are available in our mobile app.",
	Example: `  pixlet community list-icons`,
	Long:    `This command lists all in your icons that are supported by our mobile app.`,
	RunE:    listIcons,
}

func listIcons(cmd *cobra.Command, args []string) error {
	iconSet := []string{}
	for icon := range icons.IconsMap {
		iconSet = append(iconSet, icon)
	}

	sort.Strings(iconSet)
	for _, icon := range iconSet {
		fmt.Println(icon)
	}

	return nil
}
