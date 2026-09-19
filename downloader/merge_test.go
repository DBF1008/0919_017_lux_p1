package downloader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMergeMultiPart(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.mp4")
	parts := []*FilePartMeta{
		{Index: 0, Start: 0, End: 3, Cur: 4},
		{Index: 1, Start: 4, End: 7, Cur: 8},
	}
	contents := []string{"0123", "4567"}
	for i, part := range parts {
		file, err := os.Create(filePartPath(filePath, part))
		if err != nil {
			t.Fatal(err)
		}
		if err := writeFilePartMeta(file, part); err != nil {
			t.Fatal(err)
		}
		if _, err := file.WriteString(contents[i]); err != nil {
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}

	if err := mergeMultiPart(filePath, parts); err != nil {
		t.Fatal(err)
	}
	merged, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(merged) != "01234567" {
		t.Fatalf("unexpected merged content: %q", merged)
	}
	for _, part := range parts {
		if _, err := os.Stat(filePartPath(filePath, part)); !os.IsNotExist(err) {
			t.Errorf("part file %f should have been removed", part.Index)
		}
	}
}

func TestMergeMultiPartMissingPart(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.mp4")
	parts := []*FilePartMeta{
		{Index: 0, Start: 0, End: 3, Cur: 4},
	}
	// The part file does not exist, mergeMultiPart must return an error
	// instead of panicking.
	if err := mergeMultiPart(filePath, parts); err == nil {
		t.Fatal("expected an error for the missing part file")
	}
}
