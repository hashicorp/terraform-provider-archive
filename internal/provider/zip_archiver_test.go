package archive

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestZipArchiver_Content(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-content.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveContent([]byte("This is some content"), "content.txt"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, map[string][]byte{
		"content.txt": []byte("This is some content"),
	})
}

func TestZipArchiver_File(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-file.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveFile("./test-fixtures/test-dir/test-file.txt"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, map[string][]byte{
		"test-file.txt": []byte("This is test content"),
	})
}

func TestZipArchiver_FileMode(t *testing.T) {
	file, err := os.CreateTemp("", "archive-file-mode-test.zip")
	if err != nil {
		t.Fatal(err)
	}

	var (
		zipFilePath = file.Name()
		toZipPath   = filepath.FromSlash("./test-fixtures/test-dir/test-file.txt")
	)

	stringArray := [5]string{"0444", "0644", "0666", "0744", "0777"}
	for _, element := range stringArray {
		archiver := NewZipArchiver(zipFilePath)
		archiver.SetOutputFileMode(element)
		if err := archiver.ArchiveFile(toZipPath); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		ensureFileMode(t, zipFilePath, element)
	}
}

func TestZipArchiver_FileModified(t *testing.T) {
	var (
		zipFilePath = filepath.Join(t.TempDir(), "archive-file-modified.zip")
		toZipPath   = filepath.FromSlash("./test-fixtures/test-dir/test-file.txt")
	)

	var zipFunc = func() {
		archiver := NewZipArchiver(zipFilePath)
		if err := archiver.ArchiveFile(toZipPath); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}

	zipFunc()

	expectedContents, err := os.ReadFile(zipFilePath)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	//touch file modified, in the future just in case of weird race issues
	newTime := time.Now().Add(1 * time.Hour)
	if err := os.Chtimes(toZipPath, newTime, newTime); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	zipFunc()

	actualContents, err := os.ReadFile(zipFilePath)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !bytes.Equal(expectedContents, actualContents) {
		t.Fatalf("zip contents do not match, potentially a modified time issue")
	}
}

func TestZipArchiver_Dir(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir/test-dir1", []string{""}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file2.txt": []byte("This is file 2"),
		"file3.txt": []byte("This is file 3"),
	})
}

func TestZipArchiver_Dir_Exclude(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir-exclude.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir/test-dir1", []string{"file2.txt"}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file3.txt": []byte("This is file 3"),
	})
}

func TestZipArchiver_Dir_Exclude_With_Directory(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir-exclude-dir.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir", []string{"test-dir1", "test-dir2/file2.txt"}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, map[string][]byte{
		"test-dir2/file1.txt": []byte("This is file 1"),
		"test-dir2/file3.txt": []byte("This is file 3"),
		"test-file.txt":       []byte("This is test content"),
	})
}

func TestZipArchiver_Multiple(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-content.zip")

	content := map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file2.txt": []byte("This is file 2"),
		"file3.txt": []byte("This is file 3"),
	}

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveMultiple(content); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, content)
}

func TestZipArchiver_Dir_With_Symlink_File(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir-with-symlink-file.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir-with-symlink-file/target", []string{""}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, map[string][]byte{
		"file1.txt":               []byte("This is file 1"),
		"test-symlink1.txt":       []byte("This is test content"),
		"test-nested-symlink.txt": []byte("This is test content"),
	})
}

func TestZipArchiver_Dir_With_Symlink_Dir(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir-with-symlink-dir.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveDir("./test-fixtures/test-symlink-dir1/target", []string{""}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, map[string][]byte{
		"file1.txt":                                                 []byte("This is file 1"),
		"symlink-to-sample/file1.txt":                               []byte("This is file 1"),
		"symlink-to-sample/file2.txt":                               []byte("This is file 2"),
		"symlink-to-sample/file3.txt":                               []byte("This is file 3"),
		"symlink-to-symlink-to-sample2/file1.txt":                   []byte("This is file 1"),
		"symlink-to-symlink-to-sample2/symlink-to-sample-file1.txt": []byte("This is file 1"),
	})
}

func TestZipArchiver_Dir_With_Symlink_Dir_Exclude_With_Directory(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir-with-symlink-dir.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveDir("./test-fixtures/test-symlink-dir1/target", []string{
		"symlink-to-sample/file1.txt",
		"symlink-to-symlink-to-sample2/symlink-to-sample-file1.txt",
	}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, map[string][]byte{
		"file1.txt":                               []byte("This is file 1"),
		"symlink-to-sample/file2.txt":             []byte("This is file 2"),
		"symlink-to-sample/file3.txt":             []byte("This is file 3"),
		"symlink-to-symlink-to-sample2/file1.txt": []byte("This is file 1"),
	})
}

func ensureContents(t *testing.T, zipfilepath string, wants map[string][]byte) {
	t.Helper()
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
	t.Helper()
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
	gotContentBytes, err := io.ReadAll(r)
	if err != nil {
		t.Errorf("could not read file: %s", err)
	}

	wantContent := string(want)
	gotContent := string(gotContentBytes)
	if gotContent != wantContent {
		t.Errorf("mismatched content\ngot\n%s\nwant\n%s", gotContent, wantContent)
	}
}

func ensureFileMode(t *testing.T, zipfilepath string, outputFileMode string) {
	t.Helper()
	r, err := zip.OpenReader(zipfilepath)
	if err != nil {
		t.Fatalf("could not open zip file: %s", err)
	}
	defer r.Close()

	filemode, err := strconv.ParseUint(outputFileMode, 0, 32)
	if err != nil {
		t.Fatalf("error parsing outputFileMode value: %s", outputFileMode)
	}
	var osfilemode = os.FileMode(filemode)

	for _, cf := range r.File {
		if cf.FileInfo().IsDir() {
			continue
		}

		if cf.Mode() != osfilemode {
			t.Fatalf("Expected filemode \"%s\" but was \"%s\"", osfilemode, cf.Mode())
		}
	}
}
