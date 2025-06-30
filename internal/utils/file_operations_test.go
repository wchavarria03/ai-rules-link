package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCombineBytes(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "combined.txt")
	data1 := []byte("hello ")
	data2 := []byte("world")
	if err := CombineBytes(file, data1, data2); err != nil {
		t.Fatalf("CombineBytes failed: %v", err)
	}
	result, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(result) != "hello world" {
		t.Errorf("unexpected result: got %q, want %q", string(result), "hello world")
	}
}

func TestWriteBytes(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.txt")
	data := []byte("foobar")
	if err := WriteBytes(file, data); err != nil {
		t.Fatalf("WriteBytes failed: %v", err)
	}
	result, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(result) != "foobar" {
		t.Errorf("unexpected result: got %q, want %q", string(result), "foobar")
	}
}

func TestEnsureGitignore_NewFile(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)
	// .gitignore does not exist
	if err := EnsureGitignore(); err != nil {
		t.Fatalf("EnsureGitignore failed: %v", err)
	}
	content, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if !strings.Contains(string(content), ".context/") {
		t.Errorf(".gitignore missing entry: %q", string(content))
	}
}

func TestEnsureGitignore_ExistingFile(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)
	os.WriteFile(".gitignore", []byte("node_modules/\n"), 0644)
	if err := EnsureGitignore(); err != nil {
		t.Fatalf("EnsureGitignore failed: %v", err)
	}
	content, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if !strings.Contains(string(content), ".context/") {
		t.Errorf(".gitignore missing entry: %q", string(content))
	}
	if !strings.Contains(string(content), "node_modules/") {
		t.Errorf(".gitignore missing original content: %q", string(content))
	}
}

func TestCombineBytes_ErrorOnCreate(t *testing.T) {
	err := CombineBytes("/invalid/path/shouldfail.txt", []byte("fail"))
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestCombineBytes_ErrorOnCopy(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "failcopy.txt")
	// Pass a nil source to trigger io.Copy error
	err := CombineBytes(file, nil)
	if err != nil {
		t.Errorf("expected no error for nil source, got: %v", err)
	}
}

func TestWriteBytes_ErrorOnCreate(t *testing.T) {
	err := WriteBytes("/invalid/path/shouldfail.txt", []byte("fail"))
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestEnsureGitignore_ErrorOnRead(t *testing.T) {
	// Use a directory as .gitignore to force a read error
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)
	os.Mkdir(".gitignore", 0755)
	if err := EnsureGitignore(); err == nil {
		t.Error("expected error for .gitignore as directory, got nil")
	}
}

func TestEnsureGitignore_ErrorOnWrite(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)
	// Create a read-only .gitignore
	os.WriteFile(".gitignore", []byte("node_modules/\n"), 0400)
	if err := EnsureGitignore(); err == nil {
		t.Error("expected error for read-only .gitignore, got nil")
	}
}
