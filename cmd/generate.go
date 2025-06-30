package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"ai-rules-link/internal/service"

	"github.com/spf13/cobra"
)

var ruleFlags []string
var consolidateFlag bool
var globalFlag bool
var embeddedRules fs.FS // will be set from main.go
var forceFlag bool

func SetEmbeddedRules(fs fs.FS) {
	embeddedRules = fs
}

var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "Symlink selected rules into .cursor/rules/ for Cursor IDE integration, or consolidate all into one file if --consolidate is set",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		destRulesPath := os.Getenv("DEST_RULES_PATH")
		if destRulesPath == "" {
			destRulesPath = ".cursor/rules"
		}
		var baseDir string
		if globalFlag {
			baseDir = os.Getenv("HOME")
		} else {
			baseDir = cwd
		}

		// 1. Try XDG_CONFIG_HOME/ai-rules
		rulesDir := ""
		xdg := os.Getenv("XDG_CONFIG_HOME")
		if xdg != "" {
			candidate := filepath.Join(xdg, "ai-rules")
			if stat, err := os.Stat(candidate); err == nil && stat.IsDir() {
				rulesDir = candidate
			} else {

				fmt.Fprintf(os.Stdout, "[ai-rules-link] Warning: No rules found in $XDG_CONFIG_HOME/ai-rules (%s)\n", candidate)
			}
		}
		// 2. Try ~/ai-rules
		if rulesDir == "" {
			home, _ := os.UserHomeDir()
			candidate := filepath.Join(home, "ai-rules")
			if stat, err := os.Stat(candidate); err == nil && stat.IsDir() {
				rulesDir = candidate
			} else {
				fmt.Fprintf(os.Stdout, "[ai-rules-link] Warning: No rules found in ~/ai-rules (%s)\n", candidate)
			}
		}
		usingEmbedded := false
		if rulesDir == "" {
			usingEmbedded = true
			fmt.Fprintf(os.Stdout, "[ai-rules-link] Using embedded rules.\n")
		}

		if consolidateFlag {
			var merged []byte
			for _, flag := range ruleFlags {
				filename := flag + "rules.mdc"
				var content []byte
				var err error
				if usingEmbedded {
					content, err = fs.ReadFile(embeddedRules, filepath.Join("rules", filename))
				} else {
					content, err = os.ReadFile(filepath.Join(rulesDir, filename))
				}
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not read %s: %v\n", filename, err)
					os.Exit(1)
				}
				merged = append(merged, content...)
				merged = append(merged, '\n')
			}
			outDir := filepath.Join(baseDir, destRulesPath)
			os.MkdirAll(outDir, 0755)
			outFile := filepath.Join(outDir, "consolidatedrules.mdc")
			if err := os.WriteFile(outFile, merged, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write consolidated file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Consolidated rules written to: %s\n", outFile)
			return
		}

		if usingEmbedded {
			// Copy from embedded rules into destRulesPath
			outDir := filepath.Join(baseDir, destRulesPath)
			os.MkdirAll(outDir, 0755)
			for _, flag := range ruleFlags {
				filename := flag + "rules.mdc"
				srcPath := filepath.Join("rules", filename)
				content, err := fs.ReadFile(embeddedRules, srcPath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Embedded rules file does not exist for '%s': %s\n", flag, srcPath)
					continue
				}
				dst := filepath.Join(outDir, filename)
				if !forceFlag {
					dstContent, err := os.ReadFile(dst)
					if err == nil {
						// File exists, check if content matches
						if string(dstContent) != string(content) {
							fmt.Fprintf(os.Stdout, "[ai-rules-link] Skipping %s: destination file has been modified by the user. Use --force to overwrite.\n", dst)
							continue
						}
					}
				}
				if err := os.WriteFile(dst, content, 0644); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to copy embedded rule for %s: %v\n", filename, err)
				} else {
					fmt.Fprintf(os.Stdout, "Copied embedded %s into %s\n", filename, outDir)
				}
			}
			return
		}

		opts := service.SymlinkOptions{
			Rules:         ruleFlags,
			CanonicalDir:  rulesDir,
			DestRulesPath: filepath.Join(baseDir, destRulesPath),
			Stdout:        os.Stdout,
			Stderr:        os.Stderr,
		}
		if err := service.SymlinkRules(cmd.Context(), opts); err != nil {
			fmt.Fprintf(os.Stderr, "Symlink error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rulesCmd.Flags().StringSliceVar(&ruleFlags, "rule", nil, "Rule(s) to symlink (e.g., --rule=go --rule=docker --rule=base)")
	rulesCmd.Flags().BoolVar(&consolidateFlag, "consolidate", false, "Merge all selected rules into one file instead of symlinking")
	rulesCmd.Flags().BoolVar(&globalFlag, "global", false, "Create rules in the home directory (~/) instead of the current directory")
	rulesCmd.Flags().BoolVar(&forceFlag, "force", false, "Overwrite destination files even if they have been modified by the user")
	rootCmd.AddCommand(rulesCmd)
}
