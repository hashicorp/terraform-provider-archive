// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
	"time"
)

func TestTarArchiver_Content(t *testing.T) {
	tarFilePath := filepath.Join(t.TempDir(), "archive-content.tar.gz")

	archiver := NewTarGzArchiver(tarFilePath)
	if err := archiver.ArchiveContent([]byte("This is some content"), "content.txt"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureTarContents(t, tarFilePath, map[string][]byte{
		"content.txt": []byte("This is some content"),
	})
}

func TestTarArchiver_File(t *testing.T) {
	tarFilePath := filepath.Join(t.TempDir(), "archive-file.tar.gz")

	archiver := NewTarGzArchiver(tarFilePath)
	if err := archiver.ArchiveFile("./test-fixtures/test-dir/test-file.txt"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureTarContents(t, tarFilePath, map[string][]byte{
		"test-file.txt": []byte("This is test content"),
	})
}

func TestTarArchiver_FileMode(t *testing.T) {
	file, err := os.CreateTemp("", "archive-file-mode-test.tar.gz")
	if err != nil {
		t.Fatal(err)
	}

	var (
		tarFilePath = file.Name()
		toTarPath   = filepath.FromSlash("./test-fixtures/test-dir/test-file.txt")
	)

	stringArray := [5]string{"0444", "0644", "0666", "0744", "0777"}
	for _, element := range stringArray {
		archiver := NewTarGzArchiver(tarFilePath)
		archiver.SetOutputFileMode(element)
		if err := archiver.ArchiveFile(toTarPath); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		ensureTarFileMode(t, tarFilePath, element)
	}
}

func TestTarArchiver_FileModified(t *testing.T) {
	var (
		tarFilePath = filepath.Join(t.TempDir(), "archive-file-modified.tar.gz")
		toTarPath   = filepath.FromSlash("./test-fixtures/test-dir/test-file.txt")
	)

	var tarFunc = func() {
		archiver := NewTarGzArchiver(tarFilePath)
		if err := archiver.ArchiveFile(toTarPath); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}

	tarFunc()

	expectedContents, err := os.ReadFile(tarFilePath)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	// touch file modified, in the future just in case of weird race issues
	newTime := time.Now().Add(1 * time.Hour)
	if err := os.Chtimes(toTarPath, newTime, newTime); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	tarFunc()

	actualContents, err := os.ReadFile(tarFilePath)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !bytes.Equal(expectedContents, actualContents) {
		t.Fatalf("tar contents do not match, potentially a modified time issue")
	}
}

func TestTarArchiver_Dir(t *testing.T) {
	tarFilePath := filepath.Join(t.TempDir(), "archive-dir.tar.gz")

	archiver := NewTarGzArchiver(tarFilePath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir/test-dir1", ArchiveDirOpts{}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureTarContents(t, tarFilePath, map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file2.txt": []byte("This is file 2"),
		"file3.txt": []byte("This is file 3"),
	})
}

func TestTarArchiver_Dir_Exclude(t *testing.T) {
	tarFilePath := filepath.Join(t.TempDir(), "archive-dir-exclude.tar.gz")

	archiver := NewTarGzArchiver(tarFilePath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir/test-dir1", ArchiveDirOpts{
		Excludes: []string{"file2.txt"},
	}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureTarContents(t, tarFilePath, map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file3.txt": []byte("This is file 3"),
	})
}

func TestTarArchiver_Dir_Exclude_With_Directory(t *testing.T) {
	tarFilePath := filepath.Join(t.TempDir(), "archive-dir-exclude-dir.tar.gz")

	archiver := NewTarGzArchiver(tarFilePath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir", ArchiveDirOpts{
		Excludes: []string{"test-dir1", "test-dir2/file2.txt"},
	}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureTarContents(t, tarFilePath, map[string][]byte{
		"test-dir2/file1.txt": []byte("This is file 1"),
		"test-dir2/file3.txt": []byte("This is file 3"),
		"test-file.txt":       []byte("This is test content"),
	})
}

func TestTarArchiver_Multiple(t *testing.T) {
	tarFilePath := filepath.Join(t.TempDir(), "archive-content.tar.gz")

	content := map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file2.txt": []byte("This is file 2"),
		"file3.txt": []byte("This is file 3"),
	}

	archiver := NewTarGzArchiver(tarFilePath)
	if err := archiver.ArchiveMultiple(content); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureTarContents(t, tarFilePath, content)
}

func TestTarArchiver_Multiple_NoChange(t *testing.T) {
	tarFilePath := filepath.Join(t.TempDir(), "archive-content.tar.gz")

	content := map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file2.txt": []byte("This is file 2"),
		"file3.txt": []byte("This is file 3"),
	}

	archiver := NewTarGzArchiver(tarFilePath)
	if err := archiver.ArchiveMultiple(content); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expectedContents, err := os.ReadFile(tarFilePath)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	time.Sleep(1 * time.Second)

	archiver = NewTarGzArchiver(tarFilePath)
	if err := archiver.ArchiveMultiple(content); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	actualContents, err := os.ReadFile(tarFilePath)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !bytes.Equal(expectedContents, actualContents) {
		t.Fatalf("tar contents do not match, potentially a modified time issue")
	}
}

func TestTarArchiver_Dir_With_Symlink_File(t *testing.T) {
	tarFilePath := filepath.Join(t.TempDir(), "archive-dir-with-symlink-file.tar.gz")

	archiver := NewTarGzArchiver(tarFilePath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir-with-symlink-file", ArchiveDirOpts{}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureTarContents(t, tarFilePath, map[string][]byte{
		"test-file.txt":    []byte("This is test content"),
		"test-symlink.txt": []byte("This is test content"),
	})
}

func TestTarArchiver_Dir_DoNotExcludeSymlinkDirectories(t *testing.T) {
	tarFilePath := filepath.Join(t.TempDir(), "archive-dir-with-symlink-dir.tar.gz")

	archiver := NewTarGzArchiver(tarFilePath)
	if err := archiver.ArchiveDir("./test-fixtures", ArchiveDirOpts{}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureTarContents(t, tarFilePath, map[string][]byte{
		"test-dir/test-dir1/file1.txt":                         []byte("This is file 1"),
		"test-dir/test-dir1/file2.txt":                         []byte("This is file 2"),
		"test-dir/test-dir1/file3.txt":                         []byte("This is file 3"),
		"test-dir/test-dir2/file1.txt":                         []byte("This is file 1"),
		"test-dir/test-dir2/file2.txt":                         []byte("This is file 2"),
		"test-dir/test-dir2/file3.txt":                         []byte("This is file 3"),
		"test-dir/test-file.txt":                               []byte("This is test content"),
		"test-dir-with-symlink-dir/test-symlink-dir/file1.txt": []byte("This is file 1"),
		"test-dir-with-symlink-dir/test-symlink-dir/file2.txt": []byte("This is file 2"),
		"test-dir-with-symlink-dir/test-symlink-dir/file3.txt": []byte("This is file 3"),
		"test-dir-with-symlink-file/test-file.txt":             []byte("This is test content"),
		"test-dir-with-symlink-file/test-symlink.txt":          []byte("This is test content"),
		"test-symlink-dir/file1.txt":                           []byte("This is file 1"),
		"test-symlink-dir/file2.txt":                           []byte("This is file 2"),
		"test-symlink-dir/file3.txt":                           []byte("This is file 3"),
		"test-symlink-dir-with-symlink-file/test-file.txt":     []byte("This is test content"),
		"test-symlink-dir-with-symlink-file/test-symlink.txt":  []byte("This is test content"),
	})
}

func TestTarArchiver_Dir_ExcludeSymlinkDirectories(t *testing.T) {
	tarFilePath := filepath.Join(t.TempDir(), "archive-dir-with-symlink-dir.tar.gz")

	archiver := NewTarGzArchiver(tarFilePath)
	err := archiver.ArchiveDir("./test-fixtures", ArchiveDirOpts{
		ExcludeSymlinkDirectories: true,
	})

	regex := regexp.MustCompile(`error reading file for archival: read test-fixtures(\/|\\)test-dir-with-symlink-dir(\/|\\)test-symlink-dir: `)
	found := regex.Match([]byte(err.Error()))

	if !found {
		t.Fatalf("expedted error to match %q, got: %s", regex.String(), err.Error())
	}
}

func TestTarArchiver_Dir_Exclude_DoNotExcludeSymlinkDirectories(t *testing.T) {
	tarFilePath := filepath.Join(t.TempDir(), "archive-dir-with-symlink-dir.tar.gz")

	archiver := NewTarGzArchiver(tarFilePath)
	if err := archiver.ArchiveDir("./test-fixtures", ArchiveDirOpts{
		Excludes: []string{
			"test-symlink-dir/file1.txt",
			"test-symlink-dir-with-symlink-file/test-symlink.txt",
		},
	}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureTarContents(t, tarFilePath, map[string][]byte{
		"test-dir/test-dir1/file1.txt":                         []byte("This is file 1"),
		"test-dir/test-dir1/file2.txt":                         []byte("This is file 2"),
		"test-dir/test-dir1/file3.txt":                         []byte("This is file 3"),
		"test-dir/test-dir2/file1.txt":                         []byte("This is file 1"),
		"test-dir/test-dir2/file2.txt":                         []byte("This is file 2"),
		"test-dir/test-dir2/file3.txt":                         []byte("This is file 3"),
		"test-dir/test-file.txt":                               []byte("This is test content"),
		"test-dir-with-symlink-dir/test-symlink-dir/file1.txt": []byte("This is file 1"),
		"test-dir-with-symlink-dir/test-symlink-dir/file2.txt": []byte("This is file 2"),
		"test-dir-with-symlink-dir/test-symlink-dir/file3.txt": []byte("This is file 3"),
		"test-dir-with-symlink-file/test-file.txt":             []byte("This is test content"),
		"test-dir-with-symlink-file/test-symlink.txt":          []byte("This is test content"),
		"test-symlink-dir/file2.txt":                           []byte("This is file 2"),
		"test-symlink-dir/file3.txt":                           []byte("This is file 3"),
		"test-symlink-dir-with-symlink-file/test-file.txt":     []byte("This is test content"),
	})
}

func TestTarArchiver_Dir_Exclude_ExcludeSymlinkDirectories(t *testing.T) {
	tarFilePath := filepath.Join(t.TempDir(), "archive-dir-with-symlink-dir.tar.gz")

	archiver := NewTarGzArchiver(tarFilePath)
	err := archiver.ArchiveDir("./test-fixtures", ArchiveDirOpts{
		Excludes: []string{
			"test-dir/test-dir1/file1.txt",
			"test-symlink-dir-with-symlink-file/test-symlink.txt",
		},
		ExcludeSymlinkDirectories: true,
	})

	regex := regexp.MustCompile(`error reading file for archival: read test-fixtures(\/|\\)test-dir-with-symlink-dir(\/|\\)test-symlink-dir: `)
	found := regex.Match([]byte(err.Error()))

	if !found {
		t.Fatalf("expedted error to match %q, got: %s", regex.String(), err.Error())
	}
}

func ensureTarContents(t *testing.T, tarfilepath string, wants map[string][]byte) {
	t.Helper()

	f, err := os.Open(tarfilepath)
	if err != nil {
		t.Fatalf("could not open tar.gz file: %s", err)
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("could not open tar.gz file: %s", err)
	}

	tarReader := tar.NewReader(gzf)

	i := 0
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			t.Fatal(err)
		}

		if len(wants) < i+1 {
			t.Fatalf("unexpect file count in tar. expect %d", len(wants))
		}

		name := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			buf := bytes.Buffer{}
			io.Copy(&buf, tarReader)

			wantFile, ok := wants[name]
			if !ok {
				t.Fatalf("missing file %s in tar", name)
			}

			wantContent := string(wantFile)
			gotContent := buf.String()
			if gotContent != wantContent {
				t.Errorf("mismatched content\ngot\n%s\nwant\n%s", gotContent, wantContent)
			}

		default:
			t.Fatalf("%s : %c %s %s\n",
				"Yikes! Unable to figure out type",
				header.Typeflag,
				"in file",
				name,
			)
		}

		i++
	}
}

func ensureTarFileMode(t *testing.T, tarfilepath string, outputFileMode string) {
	t.Helper()

	f, err := os.Open(tarfilepath)
	if err != nil {
		t.Fatalf("could not open tar.gz file: %s", err)
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("could not open tar.gz file: %s", err)
	}

	tarReader := tar.NewReader(gzf)

	filemode, err := strconv.ParseUint(outputFileMode, 0, 32)
	if err != nil {
		t.Fatalf("error parsing outputFileMode value: %s", outputFileMode)
	}

	var osfilemode = os.FileMode(filemode)

	i := 0
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			t.Fatal(err)
		}

		name := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			if header.FileInfo().Mode() != osfilemode {
				t.Fatalf("Expected filemode \"%s\" but was \"%s\"", osfilemode, header.FileInfo().Mode())
			}
		default:
			t.Fatalf("%s : %c %s %s\n",
				"Yikes! Unable to figure out type",
				header.Typeflag,
				"in file",
				name,
			)
		}

		i++
	}
}
