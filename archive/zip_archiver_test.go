package archive

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestZipArchiver_Content(t *testing.T) {
	zipfilepath := "archive-content.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveContent([]byte("This is some content"), "content.txt", false); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, map[string][]byte{
		"content.txt": []byte("This is some content"),
	})
}

func TestZipArchiver_Content_WithNormalizedFilesMetadata(t *testing.T) {
	zipfilepath := "archive-content.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveContent([]byte("This is some content"), "content.txt", true); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureFileChecksum(t, zipfilepath, "952e89afb0435cd5e01e3e4cdf22c5b0")
}

func TestZipArchiver_File(t *testing.T) {
	zipfilepath := "archive-file.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveFile("./test-fixtures/test-file.txt", false); err != nil {
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
		if err := archiver.ArchiveFile(toZipPath, false); err != nil {
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

func TestZipArchiver_File_WithNormalizedFilesMetadata(t *testing.T) {
	zipfilepath := "archive-file.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveFile("./test-fixtures/test-file.txt", true); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureFileChecksum(t, zipfilepath, "86f7cb871bc437b8174fca96bf7a464f")
}

func TestZipArchiver_Dir(t *testing.T) {
	zipfilepath := "archive-dir.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir", []string{""}, false); err != nil {
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
	if err := archiver.ArchiveDir("./test-fixtures/test-dir", []string{"file2.txt"}, false); err != nil {
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
	if err := archiver.ArchiveDir("./test-fixtures/", []string{"test-dir", "test-dir2/file2.txt"}, false); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, map[string][]byte{
		"test-dir2/file1.txt": []byte("This is file 1"),
		"test-dir2/file3.txt": []byte("This is file 3"),
		"test-file.txt":       []byte("This is test content"),
	})
}

func TestZipArchiver_Dir_WithNormalizedFilesMetadata(t *testing.T) {
	zipfilepath := "archive-dir.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir", []string{""}, true); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureFileChecksum(t, zipfilepath, "dfb9a8da8c73034f51a5c3c5d822e64b")
}

func TestZipArchiver_Multiple(t *testing.T) {
	zipfilepath := "archive-content.zip"
	content := map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file2.txt": []byte("This is file 2"),
		"file3.txt": []byte("This is file 3"),
	}

	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveMultiple(content, false); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, content)

}

func TestZipArchiver_Multiple_WithNormalizedFilesMetadata(t *testing.T) {
	zipfilepath := "archive-content.zip"
	content := map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file2.txt": []byte("This is file 2"),
		"file3.txt": []byte("This is file 3"),
	}

	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveMultiple(content, true); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureFileChecksum(t, zipfilepath, "dfb9a8da8c73034f51a5c3c5d822e64b")
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

// There are different approaches to testing the functionality. Testing the checksum is a simple yet
// functional one, since, as long as we're assured that a normalized file has a fixed content (which
// a checksum guarantees), we don't need to know/test the normalization inner details.
func ensureFileChecksum(t *testing.T, zipfilepath string, expectedChecksum string) {
	file, err := os.Open(zipfilepath)
	if err != nil {
		t.Errorf("could not open file: %s", err)
	}

	defer file.Close()

	hashWriter := md5.New()

	if _, err := io.Copy(hashWriter, file); err != nil {
		t.Errorf("could not open file: %s", err)
	}

	fileHash := hex.EncodeToString(hashWriter.Sum(nil)[:16])

	if expectedChecksum != fileHash {
		t.Errorf("the file actual checksum (%s) didn't match the expected one (%s)", fileHash, expectedChecksum)
	}
}
