package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "List all symlinks in .cursor/rules/ and their targets",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		// Determine rules source
		rulesSource := "embedded (copied/generated)"
		xdg := os.Getenv("XDG_CONFIG_HOME")
		if xdg != "" {
			candidate := filepath.Join(xdg, "ai-rules")
			if stat, err := os.Stat(candidate); err == nil && stat.IsDir() {
				rulesSource = "$XDG_CONFIG_HOME/ai-rules: " + candidate
			}
		}
		if rulesSource == "embedded (copied/generated)" {
			home, _ := os.UserHomeDir()
			candidate := filepath.Join(home, "ai-rules")
			if stat, err := os.Stat(candidate); err == nil && stat.IsDir() {
				rulesSource = "~/ai-rules: " + candidate
			}
		}
		fmt.Printf("Rules source: %s\n", rulesSource)

		rulesDir := filepath.Join(cwd, ".cursor", "rules")
		entries, err := os.ReadDir(rulesDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read .cursor/rules/: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Symlinks in .cursor/rules/:")
		for _, entry := range entries {
			if entry.Type()&os.ModeSymlink != 0 {
				linkPath := filepath.Join(rulesDir, entry.Name())
				target, err := os.Readlink(linkPath)
				if err != nil {
					fmt.Printf("  %s -> [broken symlink]\n", entry.Name())
				} else {
					fmt.Printf("  %s -> %s\n", entry.Name(), target)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
