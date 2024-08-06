// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type TarCompressionType int

const (
	TarCompressionGz TarCompressionType = iota
)

type TarArchiver struct {
	compression       TarCompressionType
	filepath          string
	outputFileMode    string // Default value "" means unset
	fileWriter        *os.File
	tarWriter         *tar.Writer
	compressionWriter io.WriteCloser
}

func NewTarGzArchiver(filepath string) Archiver {
	return NewTarArchiver(filepath, TarCompressionGz)
}

func NewTarArchiver(filepath string, compression TarCompressionType) Archiver {
	return &TarArchiver{
		filepath:    filepath,
		compression: compression,
	}
}

func (a *TarArchiver) ArchiveContent(content []byte, infilename string) error {
	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	return a.addContent(content, &tar.Header{
		Name:    infilename,
		Size:    int64(len(content)),
		ModTime: time.Time{},
	})
}

func (a *TarArchiver) ArchiveFile(infilename string) error {
	fi, err := assertValidFile(infilename)
	if err != nil {
		return err
	}

	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	header := &tar.Header{
		Name:    filepath.ToSlash(fi.Name()),
		Size:    fi.Size(),
		Mode:    int64(fi.Mode()),
		ModTime: time.Time{},
	}

	if err := a.addFile(infilename, header); err != nil {
		return err
	}

	return err
}

func (a *TarArchiver) ArchiveDir(indirname string, opts ArchiveDirOpts) error {
	err := assertValidDir(indirname)
	if err != nil {
		return err
	}

	// ensure exclusions are OS compatible paths
	for i := range opts.Excludes {
		opts.Excludes[i] = filepath.FromSlash(opts.Excludes[i])
	}

	// Determine whether an empty archive would be generated.
	isArchiveEmpty := true

	err = filepath.Walk(indirname, a.createWalkFunc("", indirname, opts, &isArchiveEmpty, true))
	if err != nil {
		return err
	}

	// Return an error if an empty archive would be generated.
	if isArchiveEmpty {
		return fmt.Errorf("archive has not been created as it would be empty")
	}

	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	return filepath.Walk(indirname, a.createWalkFunc("", indirname, opts, &isArchiveEmpty, false))
}

func (a *TarArchiver) createWalkFunc(basePath, indirname string, opts ArchiveDirOpts, isArchiveEmpty *bool, dryRun bool) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error encountered during file walk: %s", err)
		}

		relname, err := filepath.Rel(indirname, path)
		if err != nil {
			return fmt.Errorf("error relativizing file for archival: %s", err)
		}

		archivePath := filepath.Join(basePath, relname)

		isMatch, err := checkMatch(archivePath, opts.Excludes)
		if err != nil {
			return fmt.Errorf("error checking excludes matches: %w", err)
		}

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

		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}

			realInfo, err := os.Stat(realPath)
			if err != nil {
				return err
			}

			if realInfo.IsDir() {
				if !opts.ExcludeSymlinkDirectories {
					return filepath.Walk(realPath, a.createWalkFunc(archivePath, realPath, opts, isArchiveEmpty, dryRun))
				} else {
					return filepath.SkipDir
				}
			}

			info = realInfo
		}

		*isArchiveEmpty = false

		if dryRun {
			return nil
		}

		header := &tar.Header{
			Name:    filepath.ToSlash(archivePath),
			Size:    info.Size(),
			Mode:    int64(info.Mode()),
			ModTime: time.Time{},
		}

		return a.addFile(path, header)
	}
}

func (a *TarArchiver) ArchiveMultiple(content map[string][]byte) error {
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
		header := &tar.Header{
			Name:    filepath.ToSlash(filename),
			Size:    int64(len(content[filename])),
			ModTime: time.Time{},
		}

		if err := a.addContent(content[filename], header); err != nil {
			return err
		}
	}
	return nil
}

func (a *TarArchiver) SetOutputFileMode(outputFileMode string) {
	a.outputFileMode = outputFileMode
}

func (a *TarArchiver) open() error {
	var err error

	a.fileWriter, err = os.Create(filepath.ToSlash(a.filepath))
	if err != nil {
		return err
	}

	switch a.compression {
	case TarCompressionGz:
		a.compressionWriter = gzip.NewWriter(a.fileWriter)
	}

	a.tarWriter = tar.NewWriter(a.compressionWriter)
	return nil
}

func (a *TarArchiver) close() {
	if a.tarWriter != nil {
		err := a.tarWriter.Close()
		if err != nil {
			fmt.Printf("error closing tarwriter : %s\n\n", err)
		}
		a.tarWriter = nil
	}
	if a.compressionWriter != nil {
		err := a.compressionWriter.Close()
		if err != nil {
			fmt.Printf("error closing compressionWriter : %s\n\n", err)
		}
		a.compressionWriter = nil
	}
	if a.fileWriter != nil {
		err := a.fileWriter.Close()
		if err != nil {
			fmt.Printf("error closing fileWriter: %s\n\n", err)
		}
		a.fileWriter = nil
	}
}

func (a *TarArchiver) addFile(filePath string, header *tar.Header) error {
	if header == nil {
		return fmt.Errorf("tar.Header is nil")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file '%s', got error '%w'", filePath, err)
	}
	defer file.Close()

	if a.outputFileMode != "" {
		fileMode, err := strconv.ParseInt(a.outputFileMode, 0, 32)
		if err != nil {
			return fmt.Errorf("error parsing output_file_mode value: %s", a.outputFileMode)
		}
		header.Mode = fileMode
	}

	err = a.tarWriter.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("could not write header for file '%s', got error '%w'", filePath, err)
	}

	_, err = io.Copy(a.tarWriter, file)
	if err != nil {
		return fmt.Errorf("error reading file for archival: %s", err)
	}

	return nil
}

func (a *TarArchiver) addContent(content []byte, header *tar.Header) error {
	if header == nil {
		return errors.New("tar.Header is nil")
	}

	if a.outputFileMode != "" {
		filemode, err := strconv.ParseInt(a.outputFileMode, 0, 32)
		if err != nil {
			return fmt.Errorf("error parsing output_file_mode value: %s", a.outputFileMode)
		}
		header.Mode = filemode
	}

	if err := a.tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("could not write header, got error '%w'", err)
	}

	_, err := a.tarWriter.Write(content)
	if err != nil {
		return fmt.Errorf("could not copy data to the tarball, got error '%w'", err)
	}

	return nil
}
