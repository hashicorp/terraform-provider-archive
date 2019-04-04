package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	uint32max = (1 << 32) - 1
)

type ZipArchiver struct {
	filepath   string
	filewriter *os.File
	writer     *zip.Writer
}

func NewZipArchiver(filepath string) Archiver {
	return &ZipArchiver{
		filepath: filepath,
	}
}

func (a *ZipArchiver) ArchiveContent(content []byte, infilename string, normalizeFilesMetadata bool) error {
	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	var f io.Writer
	var err error

	if normalizeFilesMetadata {
		fh := prepareEmptyHeader(content, infilename)
		normalizeCompressingFile(fh)

		f, err = a.writer.CreateHeader(fh)
		if err != nil {
			return fmt.Errorf("error creating file inside archive: %s", err)
		}
	} else {
		f, err = a.writer.Create(filepath.ToSlash(infilename))
		if err != nil {
			return err
		}
	}

	_, err = f.Write(content)
	return err
}

func (a *ZipArchiver) ArchiveFile(infilename string, normalizeFilesMetadata bool) error {
	fi, err := assertValidFile(infilename)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(infilename)
	if err != nil {
		return err
	}

	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	fh, err := zip.FileInfoHeader(fi)
	if err != nil {
		return fmt.Errorf("error creating file header: %s", err)
	}
	fh.Name = filepath.ToSlash(fi.Name())

	if normalizeFilesMetadata {
		normalizeCompressingFile(fh)
	} else {
		fh.Method = zip.Deflate
		// fh.Modified alone isn't enough when using a zero value
		fh.SetModTime(time.Time{})
	}

	f, err := a.writer.CreateHeader(fh)
	if err != nil {
		return fmt.Errorf("error creating file inside archive: %s", err)
	}

	_, err = f.Write(content)
	return err
}

func checkMatch(fileName string, excludes []string) (value bool) {
	for _, exclude := range excludes {
		if exclude == "" {
			continue
		}

		if exclude == fileName {
			return true
		}
	}
	return false
}

// The basic file header is very simple. The UncompressedSize logic is not a real-world use case
// in this context, but "640K ought to be enough for anybody".
//
// For reference, see golang/src/archive/zip/struct.go.
func prepareEmptyHeader(content []byte, infilename string) *zip.FileHeader {
	fh := &zip.FileHeader{
		Name:               filepath.ToSlash(infilename),
		UncompressedSize64: uint64(len(content)),
	}

	if fh.UncompressedSize64 > uint32max {
		fh.UncompressedSize = uint32max
	} else {
		fh.UncompressedSize = uint32(fh.UncompressedSize64)
	}

	return fh
}

// Normalize the fields:
//
// - no compression, so the compressed stream is essentially a copy;
// - fixed date;
// - fixed file permissions.
//
func normalizeCompressingFile(fh *zip.FileHeader) {
	fh.Method = zip.Store
	fh.SetModTime(time.Date(1981, 4, 10, 0, 0, 0, 0, time.UTC))
	fh.SetMode(0644)
}

func (a *ZipArchiver) ArchiveDir(indirname string, excludes []string, normalizeFilesMetadata bool) error {
	_, err := assertValidDir(indirname)
	if err != nil {
		return err
	}

	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	return filepath.Walk(indirname, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return fmt.Errorf("error encountered during file walk: %s", err)
		}

		relname, err := filepath.Rel(indirname, path)
		if err != nil {
			return fmt.Errorf("error relativizing file for archival: %s", err)
		}

		isMatch := checkMatch(relname, excludes)

		if info.IsDir() {
			if isMatch {
				return filepath.SkipDir
			}
			return nil
		}

		if isMatch {
			return nil
		}

		if err != nil {
			return err
		}

		fh, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("error creating file header: %s", err)
		}
		fh.Name = filepath.ToSlash(relname)

		if normalizeFilesMetadata {
			normalizeCompressingFile(fh)
		} else {
			fh.Method = zip.Deflate
			// fh.Modified alone isn't enough when using a zero value
			fh.SetModTime(time.Time{})
		}

		f, err := a.writer.CreateHeader(fh)
		if err != nil {
			return fmt.Errorf("error creating file inside archive: %s", err)
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file for archival: %s", err)
		}
		_, err = f.Write(content)
		return err
	})
}

func (a *ZipArchiver) ArchiveMultiple(content map[string][]byte, normalizeFilesMetadata bool) error {
	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	// Ensure files are processed in the same order so hashes don't change
	keys := make([]string, len(content))
	i := 0
	for k := range content {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, filename := range keys {
		var f io.Writer
		var err error

		if normalizeFilesMetadata {
			fh := prepareEmptyHeader(content[filename], filename)
			normalizeCompressingFile(fh)

			f, err = a.writer.CreateHeader(fh)
			if err != nil {
				return fmt.Errorf("error creating file inside archive: %s", err)
			}
		} else {
			f, err = a.writer.Create(filepath.ToSlash(filename))
			if err != nil {
				return err
			}
		}

		_, err = f.Write(content[filename])
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *ZipArchiver) open() error {
	f, err := os.Create(a.filepath)
	if err != nil {
		return err
	}
	a.filewriter = f
	a.writer = zip.NewWriter(f)
	return nil
}

func (a *ZipArchiver) close() {
	if a.writer != nil {
		a.writer.Close()
		a.writer = nil
	}
	if a.filewriter != nil {
		a.filewriter.Close()
		a.filewriter = nil
	}
}
