package filestore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalStore(t *testing.T) {
	root := t.TempDir()
	Init(root, "/manage/export/download")

	relative := "20250101/test.xlsx"
	path, err := FilePath(relative)
	if err != nil {
		t.Fatalf("FilePath() error = %v", err)
	}
	if !strings.HasPrefix(path, root) {
		t.Fatalf("FilePath() = %q, want under %q", path, root)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("xlsx"), 0o644); err != nil {
		t.Fatal(err)
	}

	opened, err := Open(relative)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if opened != path {
		t.Fatalf("Open() = %q, want %q", opened, path)
	}

	if got := URL(relative); got != "/manage/export/download/20250101/test.xlsx" {
		t.Fatalf("URL() = %q", got)
	}

	if _, err := FilePath("../escape.xlsx"); err == nil {
		t.Fatal("FilePath() should reject directory traversal")
	}
}
