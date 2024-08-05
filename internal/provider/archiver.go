// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"fmt"
	"os"
)

type ArchiveDirOpts struct {
	Excludes                  []string
	ExcludeSymlinkDirectories bool
}

type Archiver interface {
	ArchiveContent(content []byte, infilename string) error
	ArchiveFile(infilename string) error
	ArchiveDir(indirname string, opts ArchiveDirOpts) error
	ArchiveMultiple(content map[string][]byte) error
	SetOutputFileMode(outputFileMode string)
}

type ArchiverBuilder func(outputPath string) Archiver

var archiverBuilders = map[string]ArchiverBuilder{
	"zip":    NewZipArchiver,
	"tar.gz": NewTarGzArchiver,
}

func getArchiver(archiveType string, outputPath string) Archiver {
	if builder, ok := archiverBuilders[archiveType]; ok {
		return builder(outputPath)
	}
	return nil
}

func assertValidFile(infilename string) (os.FileInfo, error) {
	fi, err := os.Stat(infilename)
	if err != nil && os.IsNotExist(err) {
		return fi, fmt.Errorf("could not archive missing file: %s", infilename)
	}
	return fi, err
}

func assertValidDir(indirname string) error {
	fi, err := os.Stat(indirname)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("could not archive missing directory: %s", indirname)
		}
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("could not archive directory that is a file: %s", indirname)
	}

	return nil
}
