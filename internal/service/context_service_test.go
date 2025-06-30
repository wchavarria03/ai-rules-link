package service

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"strings"
	"testing"
)

type memFS map[string][]byte

func (m memFS) Open(name string) (fs.File, error) {
	data, ok := m[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return &memFile{data: data, name: name}, nil
}

func (m memFS) ReadFile(name string) ([]byte, error) {
	data, ok := m[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return data, nil
}

type memFile struct {
	data []byte
	name string
	pos  int64
}

func (f *memFile) Stat() (os.FileInfo, error) { return nil, nil }
func (f *memFile) Read(b []byte) (int, error) {
	if f.pos >= int64(len(f.data)) {
		return 0, errors.New("EOF")
	}
	n := copy(b, f.data[f.pos:])
	f.pos += int64(n)
	return n, nil
}
func (f *memFile) Close() error                                 { return nil }
func (f *memFile) Seek(offset int64, whence int) (int64, error) { return 0, nil }
func (f *memFile) ReadDir(count int) ([]os.DirEntry, error)     { return nil, nil }

func TestGeneratePrompt_Success(t *testing.T) {
	fsys := memFS{
		"rules/baserules.mdc": []byte("base\n"),
		"rules/prompt.go.mdc": []byte("go\n"),
	}
	svc := &ContextService{RulesFS: fsys}
	out, err := svc.GeneratePrompt(context.Background(), "go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "base\ngo\n" {
		t.Errorf("unexpected output: %q", string(out))
	}
}

func TestGeneratePrompt_MissingBase(t *testing.T) {
	fsys := memFS{"rules/prompt.go.mdc": []byte("go\n")}
	svc := &ContextService{RulesFS: fsys}
	_, err := svc.GeneratePrompt(context.Background(), "go")
	if err == nil || !strings.Contains(err.Error(), "read base prompt") {
		t.Errorf("expected error for missing base, got: %v", err)
	}
}

func TestGeneratePrompt_MissingTech(t *testing.T) {
	fsys := memFS{"rules/baserules.mdc": []byte("base\n")}
	svc := &ContextService{RulesFS: fsys}
	_, err := svc.GeneratePrompt(context.Background(), "go")
	if err == nil || !strings.Contains(err.Error(), "read tech prompt") {
		t.Errorf("expected error for missing tech, got: %v", err)
	}
}

func TestGeneratePromptFlexible_BaseOnly(t *testing.T) {
	fsys := memFS{"rules/baserules.mdc": []byte("base\n")}
	svc := &ContextService{RulesFS: fsys}
	out, err := svc.GeneratePromptFlexible(context.Background(), "go", true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "base\n" {
		t.Errorf("unexpected output: %q", string(out))
	}
}

func TestGeneratePromptFlexible_LangOnly(t *testing.T) {
	fsys := memFS{"rules/prompt.go.mdc": []byte("go\n")}
	svc := &ContextService{RulesFS: fsys}
	out, err := svc.GeneratePromptFlexible(context.Background(), "go", false, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "go\n" {
		t.Errorf("unexpected output: %q", string(out))
	}
}

func TestGeneratePromptFlexible_BothTrue(t *testing.T) {
	fsys := memFS{}
	svc := &ContextService{RulesFS: fsys}
	_, err := svc.GeneratePromptFlexible(context.Background(), "go", true, true)
	if err == nil || !strings.Contains(err.Error(), "cannot set both baseOnly and langOnly to true") {
		t.Errorf("expected error for both true, got: %v", err)
	}
}
