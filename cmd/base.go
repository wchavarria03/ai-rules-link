package cmd

import (
	"ai-rules-link/internal/service"
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var baseCmd = &cobra.Command{
	Use:   "base",
	Short: "Generate only the base rules for the project",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := service.NewContextService(rulesFS)
		if err := service.InitializeFlexible(ctx, "go", true, false); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating base rules: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Successfully generated base rules in .context/prompt.mdc")
	},
}

func init() {
	rootCmd.AddCommand(baseCmd)
}
