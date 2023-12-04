// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type TgzArchiver struct {
	filepath       string
	outputFileMode string
	filewriter     *os.File
	gzipwriter     *gzip.Writer
	tarwriter      *tar.Writer
}

func NewTgzArchiver(filepath string) Archiver {
	return &TgzArchiver{
		filepath: filepath,
	}
}

func (a *TgzArchiver) ArchiveContent(content []byte, infilename string) error {
	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	if err := a.tarwriter.WriteHeader(&tar.Header{
		Name: infilename,
		Mode: 0600,
		Size: int64(len(content)),
	}); err != nil {
		return err
	}
	if _, err := a.tarwriter.Write(content); err != nil {
		return err
	}

	return nil
}

func (a *TgzArchiver) ArchiveFile(infilename string) error {
	fi, err := assertValidFile(infilename)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(infilename)
	if err != nil {
		return err
	}

	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	fih, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return fmt.Errorf("error creating file header: %s", err)
	}

	if a.outputFileMode != "" {
		filemode, err := strconv.ParseInt(a.outputFileMode, 0, 64)
		if err != nil {
			return fmt.Errorf("error parsing output_file_mode value: %s", a.outputFileMode)
		}
		fih.Mode = filemode
	}

	if err := a.tarwriter.WriteHeader(fih); err != nil {
		return fmt.Errorf("error creating file inside archive: %s", err)
	}

	if _, err = a.tarwriter.Write(content); err != nil {
		return err
	}

	return nil
}

func (a *TgzArchiver) ArchiveDir(indirname string, opts ArchiveDirOpts) error {
	if err := assertValidDir(indirname); err != nil {
		return err
	}

	for i := range opts.Excludes {
		opts.Excludes[i] = filepath.FromSlash(opts.Excludes[i])
	}

	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	return filepath.Walk(indirname, a.createWalkFunc("", indirname, opts))
}

func (a *TgzArchiver) createWalkFunc(basePath string, indirname string, opts ArchiveDirOpts) func(path string, info os.FileInfo, err error) error {
	return func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error encountered during file walk: %s", err)
		}

		relName, err := filepath.Rel(indirname, path)
		if err != nil {
			return fmt.Errorf("error relativizing file for archival: %s", err)
		}

		archivePath := filepath.Join(basePath, relName)

		isExcluded, err := checkMatch(archivePath, opts.Excludes)
		if err != nil {
			return fmt.Errorf("error matching file for archival: %s", err)
		}

		if fi.IsDir() {
			if isExcluded {
				return filepath.SkipDir
			}
			return nil
		}

		if isExcluded {
			return nil
		}

		if err != nil {
			return err
		}

		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			if !opts.ExcludeSymlinkDirectories {
				realPath, err := filepath.EvalSymlinks(path)
				if err != nil {
					return err
				}

				realFileInfo, err := os.Stat(realPath)
				if err != nil {
					return err
				}

				if realFileInfo.IsDir() {
					return filepath.Walk(realPath, a.createWalkFunc(archivePath, realPath, opts))
				}

				fi = realFileInfo
			}
		}

		fih, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return fmt.Errorf("error creating file header: %s", err)
		}

		if a.outputFileMode != "" {
			filemode, err := strconv.ParseInt(a.outputFileMode, 0, 64)
			if err != nil {
				return fmt.Errorf("error parsing output_file_mode value: %s", a.outputFileMode)
			}
			fih.Mode = filemode
		}

		err = a.tarwriter.WriteHeader(fih)
		if err != nil {
			return fmt.Errorf("error creating file inside archive: %s", err)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file for archival: %s", err)
		}
		_, err = a.tarwriter.Write(content)
		return err
	}
}

func (a *TgzArchiver) ArchiveMultiple(content map[string][]byte) error {
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
		if err := a.tarwriter.WriteHeader(&tar.Header{
			Name: filename,
			Mode: 0600,
			Size: int64(len(content[filename])),
		}); err != nil {
			return err
		}
		if _, err := a.tarwriter.Write(content[filename]); err != nil {
			return err
		}
	}
	return nil
}

func (a *TgzArchiver) SetOutputFileMode(outputFileMode string) {
	a.outputFileMode = outputFileMode
}

func (a *TgzArchiver) open() error {
	f, err := os.Create(a.filepath)
	if err != nil {
		return err
	}
	a.filewriter = f
	a.gzipwriter = gzip.NewWriter(f)
	a.tarwriter = tar.NewWriter(a.gzipwriter)
	return nil
}

func (a *TgzArchiver) close() {
	if a.tarwriter != nil {
		a.tarwriter.Close()
		a.tarwriter = nil
	}
	if a.gzipwriter != nil {
		a.gzipwriter.Close()
		a.gzipwriter = nil
	}
	if a.filewriter != nil {
		a.filewriter.Close()
		a.filewriter = nil
	}
}
