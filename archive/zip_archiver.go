package archive

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
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

func (a *ZipArchiver) ArchiveContent(content []byte, infilename string) error {
	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	f, err := a.writer.Create(infilename)
	if err != nil {
		return err
	}

	_, err = f.Write(content)
	return err
}

func (a *ZipArchiver) ArchiveFile(infilename string) error {
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
	fh.Name = fi.Name()
	fh.Method = zip.Deflate
	// fh.Modified alone isn't enough when using a zero value
	fh.SetModTime(time.Time{})

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
		matchCheck, _ := regexp.MatchString("^[/\\\\]?"+exclude+"([/\\\\].*)?$", fileName)
		if matchCheck {
			return true
		}
	}
	return false
}

func (a *ZipArchiver) ArchiveDir(indirname string, excludes []string) error {
	_, err := assertValidDir(indirname)
	if err != nil {
		return err
	}

	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	return filepath.Walk(indirname, func(path string, info os.FileInfo, err error) error {
		tmpFilePath, _ := filepath.Abs(path)
		tmpLongFileName, _ := filepath.Abs(indirname)
		shortFileName := strings.TrimPrefix(tmpFilePath, tmpLongFileName)
		isMatch := checkMatch(shortFileName, excludes)
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

		relname, err := filepath.Rel(indirname, path)
		if err != nil {
			return fmt.Errorf("error relativizing file for archival: %s", err)
		}
		fh, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("error creating file header: %s", err)
		}
		fh.Name = relname
		fh.Method = zip.Deflate
		// fh.Modified alone isn't enough when using a zero value
		fh.SetModTime(time.Time{})

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

func (a *ZipArchiver) ArchiveMultiple(content map[string][]byte) error {
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
		f, err := a.writer.Create(filename)
		if err != nil {
			return err
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
