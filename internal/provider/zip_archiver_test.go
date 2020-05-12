package archive

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestZipArchiver_Content(t *testing.T) {
	zipfilepath := "archive-content.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveContent([]byte("This is some content"), "content.txt"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, map[string][]byte{
		"content.txt": []byte("This is some content"),
	})
}

func TestZipArchiver_File(t *testing.T) {
	zipfilepath := "archive-file.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveFile("./test-fixtures/test-file.txt"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, map[string][]byte{
		"test-file.txt": []byte("This is test content"),
	})
}
func TestZipArchiver_FileModified(t *testing.T) {
	var (
		zipFilePath = filepath.FromSlash("archive-file.zip")
		toZipPath   = filepath.FromSlash("./test-fixtures/test-file.txt")
	)

	var zip = func() {
		archiver := NewZipArchiver(zipFilePath)
		if err := archiver.ArchiveFile(toZipPath); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}

	zip()

	expectedContents, err := ioutil.ReadFile(zipFilePath)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	//touch file modified, in the future just in case of weird race issues
	newTime := time.Now().Add(1 * time.Hour)
	if err := os.Chtimes(toZipPath, newTime, newTime); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	zip()

	actualContents, err := ioutil.ReadFile(zipFilePath)
	if err != nil {
		t.Fatalf("unexpecte error: %s", err)
	}

	if !bytes.Equal(expectedContents, actualContents) {
		t.Fatalf("zip contents do not match, potentially a modified time issue")
	}
}

func TestZipArchiver_Dir(t *testing.T) {
	zipfilepath := "archive-dir.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir", []string{""}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file2.txt": []byte("This is file 2"),
		"file3.txt": []byte("This is file 3"),
	})
}

func TestZipArchiver_Dir_Exclude(t *testing.T) {
	zipfilepath := "archive-dir.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir", []string{"file2.txt"}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file3.txt": []byte("This is file 3"),
	})
}

func TestZipArchiver_Dir_Exclude_With_Directory(t *testing.T) {
	zipfilepath := "archive-dir.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir3", []string{"file1.txt", "test-subdir/file2.txt"}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, map[string][]byte{
		"test-subdir/file1.txt": []byte("This is file 1"),
		"test-subdir/file3.txt": []byte("This is file 3"),
		"file2.txt":             []byte("This is file 2"),
		"file3.txt":             []byte("This is file 3"),
	})
}

func TestZipArchiver_Dir_With_Symlink_File(t *testing.T) {
	zipfilepath := "archive-dir-with-symlink-file.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir-with-symlink-file", []string{""}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, map[string][]byte{
		"file1.txt":               []byte("This is file 1"),
		"test-symlink1.txt":       []byte("This is test content"),
		"test-nested-symlink.txt": []byte("This is test content"),
	})
}

func TestZipArchiver_Dir_With_Symlink_Dir(t *testing.T) {
	zipfilepath := "archive-dir-with-symlink-dir.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveDir("./test-fixtures/test-symlink-dir1", []string{""}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, map[string][]byte{
		"file1.txt":                                  []byte("This is file 1"),
		"symlink-dir1/file1.txt":                     []byte("This is file 1"),
		"symlink-dir1/file2.txt":                     []byte("This is file 2"),
		"symlink-dir1/file3.txt":                     []byte("This is file 3"),
		"nested-symlink-dir/file1.txt":               []byte("This is file 1"),
		"nested-symlink-dir/test-nested-symlink.txt": []byte("This is test content"),
		"nested-symlink-dir/test-symlink1.txt":       []byte("This is test content"),
	})
}

func TestZipArchiver_Dir_With_Symlink_Dir_Exclude_With_Directory(t *testing.T) {
	zipfilepath := "archive-dir-with-symlink-dir.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveDir("./test-fixtures/test-symlink-dir1", []string{
		"symlink-dir1/file1.txt",
		"nested-symlink-dir/test-nested-symlink.txt",
	}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"symlink-dir1/file2.txt":       []byte("This is file 2"),
		"symlink-dir1/file3.txt":       []byte("This is file 3"),
		"nested-symlink-dir/file1.txt": []byte("This is file 1"),
		"nested-symlink-dir/test-symlink1.txt": []byte("This is test content"),
	})
}

func TestZipArchiver_Multiple(t *testing.T) {
	zipfilepath := "archive-content.zip"
	content := map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file2.txt": []byte("This is file 2"),
		"file3.txt": []byte("This is file 3"),
	}

	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveMultiple(content); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, content)

}

func ensureContents(t *testing.T, zipfilepath string, wants map[string][]byte) {
	r, err := zip.OpenReader(zipfilepath)
	if err != nil {
		t.Fatalf("could not open zip file: %s", err)
	}
	defer r.Close()

	if len(r.File) != len(wants) {
		t.Errorf("mismatched file count, got %d, want %d", len(r.File), len(wants))
	}
	for _, cf := range r.File {
		ensureContent(t, wants, cf)
	}
}

func ensureContent(t *testing.T, wants map[string][]byte, got *zip.File) {
	want, ok := wants[got.Name]
	if !ok {
		t.Errorf("additional file in zip: %s", got.Name)
		return
	}

	r, err := got.Open()
	if err != nil {
		t.Errorf("could not open file: %s", err)
	}
	defer r.Close()
	gotContentBytes, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("could not read file: %s", err)
	}

	wantContent := string(want)
	gotContent := string(gotContentBytes)
	if gotContent != wantContent {
		t.Errorf("mismatched content\ngot\n%s\nwant\n%s", gotContent, wantContent)
	}
}
