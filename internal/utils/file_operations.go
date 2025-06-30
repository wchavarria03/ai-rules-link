package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// CombineBytes creates or truncates the destination file and writes all source byte slices to it in order.
func CombineBytes(dest string, sources ...[]byte) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, src := range sources {
		if _, err := io.Copy(out, bytes.NewReader(src)); err != nil {
			return err
		}
	}
	return nil
}

// WriteBytes writes the given byte slice to the specified file path, creating or truncating the file.
func WriteBytes(dest string, data []byte) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// EnsureGitignore ensures that .context/ is in .gitignore. Returns error for testability.
func EnsureGitignore() error {
	gitignorePath := ".gitignore"
	entry := ".context/"

	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			// .gitignore does not exist, create it with the entry
			newContent := fmt.Sprintf("\n# AI-generated context files\n%s\n", entry)
			if err := os.WriteFile(gitignorePath, []byte(newContent), 0644); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if !strings.Contains(string(content), entry) {
		// Entry does not exist, append it
		f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := f.WriteString(fmt.Sprintf("\n# AI-generated context files\n%s\n", entry)); err != nil {
			return err
		}
	}
	return nil
}
