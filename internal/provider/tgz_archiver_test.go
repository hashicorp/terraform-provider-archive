// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestTgzArchiver_Content(t *testing.T) {
	tgzFilePath := filepath.Join(t.TempDir(), "archive-content.tgz")

	archiver := NewTgzArchiver(tgzFilePath)
	if err := archiver.ArchiveContent([]byte("This is some content"), "content.txt"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureTgzContents(t, tgzFilePath, map[string][]byte{
		"content.txt": []byte("This is some content"),
	})
}

func TestTgzArchiver_File(t *testing.T) {
	tgzFilePath := filepath.Join(t.TempDir(), "archive-file.tgz")

	archiver := NewTgzArchiver(tgzFilePath)
	if err := archiver.ArchiveFile("./test-fixtures/test-dir/test-file.txt"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureTgzContents(t, tgzFilePath, map[string][]byte{
		"test-file.txt": []byte("This is test content"),
	})
}

func TestTgzArchiver_FileMode(t *testing.T) {
	file, err := os.CreateTemp("", "archive-file-mode-test.tgz")
	if err != nil {
		t.Fatal(err)
	}

	var (
		tgzFilePath = file.Name()
		toTgzPath   = filepath.FromSlash("./test-fixtures/test-dir/test-file.txt")
	)

	for _, element := range []string{"0444", "0644", "0666", "0744", "0777"} {
		archiver := NewTgzArchiver(tgzFilePath)
		archiver.SetOutputFileMode(element)
		if err := archiver.ArchiveFile(toTgzPath); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		ensureTgzFileMode(t, tgzFilePath, element)
	}
}

func TestTgzArchiver_Dir(t *testing.T) {
	tgzFilePath := filepath.Join(t.TempDir(), "archive-dir.tgz")

	archiver := NewTgzArchiver(tgzFilePath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir/test-dir1", ArchiveDirOpts{}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureTgzContents(t, tgzFilePath, map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file2.txt": []byte("This is file 2"),
		"file3.txt": []byte("This is file 3"),
	})
}

func TestTgzArchiver_Multiple(t *testing.T) {
	tgzFilePath := filepath.Join(t.TempDir(), "archive-content.tgz")

	content := map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file2.txt": []byte("This is file 2"),
		"file3.txt": []byte("This is file 3"),
	}

	archiver := NewTgzArchiver(tgzFilePath)
	if err := archiver.ArchiveMultiple(content); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureTgzContents(t, tgzFilePath, content)
}

func ensureTgzContents(t *testing.T, tgzfilepath string, wants map[string][]byte) {
	t.Helper()
	f, err := os.Open(tgzfilepath)
	if err != nil {
		t.Fatalf("could not open tgz file: %s", err)
	}
	defer f.Close()

	gzipreader, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("could not open tgz reader: %s", err)
	}
	defer gzipreader.Close()

	tarreader := tar.NewReader(gzipreader)
	count := 0

	for {
		hdr, err := tarreader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error iterating tgz file: %s", err)
		}
		ensureTgzContent(t, wants, hdr, tarreader)

		count++
	}

	if count != len(wants) {
		t.Errorf("mismatched file count, got %d, want %d", count, len(wants))
	}
}

func ensureTgzContent(t *testing.T, wants map[string][]byte, hdr *tar.Header, got *tar.Reader) {
	t.Helper()
	want, ok := wants[hdr.Name]
	if !ok {
		t.Errorf("additional file in tgz: %s", hdr.Name)
		return
	}

	gotContentBytes, err := io.ReadAll(got)
	if err != nil {
		t.Errorf("could not read file: %s", err)
	}

	wantContent := string(want)
	gotContent := string(gotContentBytes)
	if gotContent != wantContent {
		t.Errorf("mismatched content\ngot\n%s\nwant\n%s", gotContent, wantContent)
	}
}

func ensureTgzFileMode(t *testing.T, tgzfilepath string, outputFileMode string) {
	t.Helper()
	f, err := os.Open(tgzfilepath)
	if err != nil {
		t.Fatalf("could not open tgz file: %s", err)
	}
	defer f.Close()

	gzipreader, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("could not open tgz reader: %s", err)
	}
	defer gzipreader.Close()

	filemode, err := strconv.ParseInt(outputFileMode, 0, 64)
	if err != nil {
		t.Fatalf("error parsing outputFileMode value: %s", outputFileMode)
	}

	tarreader := tar.NewReader(gzipreader)
	for {
		hdr, err := tarreader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error iterating tgz file: %s", err)
		}

		if hdr.Typeflag == tar.TypeDir {
			continue
		}

		if hdr.Mode != filemode {
			t.Fatalf("Expected filemode \"%d\" but was \"%d\"", filemode, hdr.Mode)
		}
	}
}
