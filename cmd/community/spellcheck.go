package community

import (
	"fmt"
	"io"
	"os"

	"github.com/client9/misspell"
	"github.com/spf13/cobra"
)

var (
	FixSpelling    bool
	SilentSpelling bool
)

func init() {
	SpellCheckCmd.Flags().BoolVarP(&FixSpelling, "fix", "f", false, "fixes spelling mistakes automatically")
	SpellCheckCmd.Flags().BoolVarP(&SilentSpelling, "silent", "s", false, "silences spelling mistakes")
}

var SpellCheckCmd = &cobra.Command{
	Use:   "spell-check <filespec>",
	Short: "Spell check for a file",
	Example: `  pixlet community spell-check manifest.yaml
  pixlet community spell-check app.star`,
	Long: `This command checks the spelling of strings located in a file. This can be used
both for a manifest and Tidbyt app.`,
	Args: cobra.ExactArgs(1),
	RunE: SpellCheck,
}

func SpellCheck(cmd *cobra.Command, args []string) error {
	// Load file for checking.
	f, err := os.OpenFile(args[0], os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	// Create replacer.
	r := misspell.Replacer{
		Replacements: misspell.DictMain,
	}

	// Tidbyt is primarily in US markets. We only ship a US power plug, and all
	// materials are in the US locale. In the future, we will need to consider
	// how we manage spell check as we look to support more markets.
	r.AddRuleList(misspell.DictAmerican)
	r.Compile()

	// Run replacer.
	updated, diffs := r.Replace(string(b))

	// If FixSpelling is true, we only want to fix spelling and return
	if FixSpelling {
		// Updating a file in line gets a bit tricky. The file would first have
		// to be cleared of the file contents, which feels dangerous. So
		// instead, create a temp file, write the contents, and then replace
		// the original file with the new file.
		temp := args[0] + ".temp"
		t, err := os.Create(temp)
		if err != nil {
			return fmt.Errorf("could not create file: %w", err)
		}
		defer t.Close()

		_, err = t.WriteString(updated)
		if err != nil {
			return fmt.Errorf("could not update file: %w", err)
		}

		err = os.Rename(temp, args[0])
		if err != nil {
			return fmt.Errorf("could not replace file: %w", err)
		}

		return nil
	}

	if !SilentSpelling {
		for _, diff := range diffs {
			fmt.Printf("`%s` is a misspelling of `%s` at line: %d\n", diff.Original, diff.Corrected, diff.Line)
		}
	}

	// Return error if there are any diffs.
	if len(diffs) > 0 {
		return fmt.Errorf("%s contains spelling errors", args[0])
	}

	return nil
}
