package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SymlinkOptions struct {
	Rules         []string
	CanonicalDir  string
	DestRulesPath string
	Stdout        *os.File // for testability, or use io.Writer
	Stderr        *os.File
}

// SymlinkRules creates symlinks for the specified rules from CanonicalDir to DestRulesPath.
func SymlinkRules(ctx context.Context, opts SymlinkOptions) error {
	if len(opts.Rules) == 0 {
		fmt.Fprintln(opts.Stderr, "No rules specified. Use --rule for each rule you want to symlink (e.g., --rule=go --rule=base)")
		return fmt.Errorf("no rules specified")
	}
	if err := os.MkdirAll(opts.DestRulesPath, 0755); err != nil {
		return fmt.Errorf("error creating %s: %w", opts.DestRulesPath, err)
	}
	for _, flag := range opts.Rules {
		flag = strings.ToLower(flag)
		filename := flag + "rules.mdc"
		src := filepath.Join(opts.CanonicalDir, filename)
		dst := filepath.Join(opts.DestRulesPath, filename)
		if _, err := os.Stat(src); os.IsNotExist(err) {
			fmt.Fprintf(opts.Stderr, "Canonical rules file does not exist for '%s': %s\n", flag, src)
			continue
		}
		info, err := os.Lstat(dst)
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(dst)
			if err == nil && target == src {
				fmt.Fprintf(opts.Stdout, "Symlink for %s already exists and is correct.\n", filename)
				continue
			}
		}
		os.Remove(dst)
		if err := os.Symlink(src, dst); err != nil {
			fmt.Fprintf(opts.Stderr, "Failed to create symlink for %s: %v\n", filename, err)
		} else {
			fmt.Fprintf(opts.Stdout, "Symlinked %s into %s\n", filename, opts.DestRulesPath)
		}
	}
	return nil
}
