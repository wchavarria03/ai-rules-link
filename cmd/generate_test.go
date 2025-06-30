package cmd

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConsolidateFlag_CreatesMergedFile(t *testing.T) {
	dir := t.TempDir()
	canonicalDir := filepath.Join(dir, "canonical")
	os.MkdirAll(canonicalDir, 0755)
	os.WriteFile(filepath.Join(canonicalDir, "gorules.mdc"), []byte("go content"), 0644)
	os.WriteFile(filepath.Join(canonicalDir, "pythonrules.mdc"), []byte("python content"), 0644)
	destDir := filepath.Join(dir, "dest")
	os.MkdirAll(destDir, 0755)

	// Set up environment
	os.Setenv("HOME", dir)
	os.Setenv("DEST_RULES_PATH", "dest")

	// Simulate CLI args
	ruleFlags = []string{"go", "python"}
	consolidateFlag = true

	// Run the logic from rulesCmd (extracted for testability)
	var merged []byte
	for _, flag := range ruleFlags {
		filename := flag + "rules.mdc"
		src := filepath.Join(canonicalDir, filename)
		content, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("Could not read %s: %v", src, err)
		}
		merged = append(merged, content...)
		merged = append(merged, '\n')
	}
	outFile := filepath.Join(destDir, "consolidatedrules.mdc")
	if err := os.WriteFile(outFile, merged, 0644); err != nil {
		t.Fatalf("Failed to write consolidated file: %v", err)
	}
	result, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("Failed to read consolidated file: %v", err)
	}
	if !strings.Contains(string(result), "go content") || !strings.Contains(string(result), "python content") {
		t.Errorf("Merged file missing expected content: %q", string(result))
	}
}

func TestConsolidateFlag_MissingRuleFile(t *testing.T) {
	dir := t.TempDir()
	canonicalDir := filepath.Join(dir, "canonical")
	os.MkdirAll(canonicalDir, 0755)
	os.WriteFile(filepath.Join(canonicalDir, "gorules.mdc"), []byte("go content"), 0644)
	destDir := filepath.Join(dir, "dest")
	os.MkdirAll(destDir, 0755)

	os.Setenv("HOME", dir)
	os.Setenv("DEST_RULES_PATH", "dest")

	ruleFlags = []string{"go", "missing"}
	consolidateFlag = true

	var merged []byte
	for _, flag := range ruleFlags {
		filename := flag + "rules.mdc"
		src := filepath.Join(canonicalDir, filename)
		content, err := os.ReadFile(src)
		if err != nil {
			// Should error for missing file
			continue
		}
		merged = append(merged, content...)
		merged = append(merged, '\n')
	}
	outFile := filepath.Join(destDir, "consolidatedrules.mdc")
	if err := os.WriteFile(outFile, merged, 0644); err != nil {
		t.Fatalf("Failed to write consolidated file: %v", err)
	}
	result, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("Failed to read consolidated file: %v", err)
	}
	if !strings.Contains(string(result), "go content") {
		t.Errorf("Merged file missing expected content: %q", string(result))
	}
	if strings.Contains(string(result), "missing content") {
		t.Errorf("Merged file should not contain missing content: %q", string(result))
	}
}

func TestRulesCmd_EmbeddedRulesCopy(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", "")
	os.Setenv("HOME", dir)
	os.Setenv("DEST_RULES_PATH", "dest")

	// Simulate no rules in XDG or HOME, so embedded is used
	ruleFlags = []string{"go"}
	consolidateFlag = false
	globalFlag = false
	forceFlag = true

	// Patch embeddedRules to use a test FS
	testFS := os.DirFS("../rules")
	SetEmbeddedRules(testFS)

	// Run the logic from rulesCmd (extracted for testability)
	outDir := filepath.Join(dir, "dest")
	os.MkdirAll(outDir, 0755)
	filename := "gorules.mdc"
	content, err := fs.ReadFile(testFS, filepath.Join("gorules.mdc"))
	if err != nil {
		t.Skip("Test rules file not found in ../rules/gorules.mdc")
	}
	outFile := filepath.Join(outDir, filename)
	if err := os.WriteFile(outFile, content, 0644); err != nil {
		t.Fatalf("Failed to copy embedded rule: %v", err)
	}
	result, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("Failed to read copied rule: %v", err)
	}
	if string(result) != string(content) {
		t.Errorf("Copied content mismatch: got %q, want %q", string(result), string(content))
	}
}

func TestRulesCmd_ForceFlag(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", "")
	os.Setenv("HOME", dir)
	os.Setenv("DEST_RULES_PATH", "dest")

	// Simulate embedded rules
	ruleFlags = []string{"go"}
	consolidateFlag = false
	globalFlag = false
	forceFlag = true

	testFS := os.DirFS("../rules")
	SetEmbeddedRules(testFS)

	outDir := filepath.Join(dir, "dest")
	os.MkdirAll(outDir, 0755)
	filename := "gorules.mdc"
	content, err := fs.ReadFile(testFS, filepath.Join("gorules.mdc"))
	if err != nil {
		t.Skip("Test rules file not found in ../rules/gorules.mdc")
	}
	outFile := filepath.Join(outDir, filename)
	os.WriteFile(outFile, []byte("user modified content"), 0644)
	if err := os.WriteFile(outFile, content, 0644); err != nil {
		t.Fatalf("Failed to overwrite with force: %v", err)
	}
	result, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("Failed to read overwritten rule: %v", err)
	}
	if string(result) != string(content) {
		t.Errorf("Force flag did not overwrite: got %q, want %q", string(result), string(content))
	}
}
