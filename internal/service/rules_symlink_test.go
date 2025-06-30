package service

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSymlinkRules_NoRulesSpecified(t *testing.T) {
	opts := SymlinkOptions{
		Rules:         nil,
		CanonicalDir:  t.TempDir(),
		DestRulesPath: t.TempDir(),
		Stdout:        os.Stdout,
		Stderr:        os.Stderr,
	}
	err := SymlinkRules(context.Background(), opts)
	if err == nil || !strings.Contains(err.Error(), "no rules specified") {
		t.Errorf("expected error for no rules specified, got: %v", err)
	}
}

func TestSymlinkRules_CanonicalFileMissing(t *testing.T) {
	dest := t.TempDir()
	opts := SymlinkOptions{
		Rules:         []string{"missing"},
		CanonicalDir:  t.TempDir(),
		DestRulesPath: dest,
		Stdout:        os.Stdout,
		Stderr:        os.Stderr,
	}
	err := SymlinkRules(context.Background(), opts)
	if err != nil {
		t.Errorf("expected no error for missing canonical file, got: %v", err)
	}
	// Should not create a symlink
	if _, err := os.Lstat(filepath.Join(dest, "missingrules.mdc")); !os.IsNotExist(err) {
		t.Errorf("expected no symlink for missing file, got: %v", err)
	}
}

func TestSymlinkRules_SymlinkCreated(t *testing.T) {
	canon := t.TempDir()
	dest := t.TempDir()
	filename := "gorules.mdc"
	canonFile := filepath.Join(canon, filename)
	os.WriteFile(canonFile, []byte("test"), 0644)
	opts := SymlinkOptions{
		Rules:         []string{"go"},
		CanonicalDir:  canon,
		DestRulesPath: dest,
		Stdout:        os.Stdout,
		Stderr:        os.Stderr,
	}
	err := SymlinkRules(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	link := filepath.Join(dest, filename)
	info, err := os.Lstat(link)
	if err != nil {
		t.Fatalf("expected symlink, got error: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Errorf("expected a symlink, got mode: %v", info.Mode())
	}
	target, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("could not read symlink: %v", err)
	}
	if target != canonFile {
		t.Errorf("symlink target mismatch: got %s, want %s", target, canonFile)
	}
}

func TestSymlinkRules_SymlinkAlreadyCorrect(t *testing.T) {
	canon := t.TempDir()
	dest := t.TempDir()
	filename := "gorules.mdc"
	canonFile := filepath.Join(canon, filename)
	os.WriteFile(canonFile, []byte("test"), 0644)
	link := filepath.Join(dest, filename)
	os.MkdirAll(dest, 0755)
	os.Symlink(canonFile, link)
	opts := SymlinkOptions{
		Rules:         []string{"go"},
		CanonicalDir:  canon,
		DestRulesPath: dest,
		Stdout:        os.Stdout,
		Stderr:        os.Stderr,
	}
	err := SymlinkRules(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should not change the symlink
	target, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("could not read symlink: %v", err)
	}
	if target != canonFile {
		t.Errorf("symlink target mismatch: got %s, want %s", target, canonFile)
	}
}
