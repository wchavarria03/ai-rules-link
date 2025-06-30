package cmd

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the base command for ai-rules-link.
var rootCmd = &cobra.Command{
	Use:   "ai-rules-link",
	Short: "A CLI to manage AI context for different tools and technologies",
	Long:  `ai-rules-link is a tool to standardize AI-assisted development by generating context-aware prompts.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var rulesFS fs.FS

// Execute runs the root command with the provided embedded rules filesystem.
func Execute(embeddedFS fs.FS) {
	rulesFS = embeddedFS
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI: %v\n", err)
		os.Exit(1)
	}
}
