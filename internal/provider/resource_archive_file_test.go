// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccArchiveFile_Resource_Basic(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string
	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: testAccArchiveFileResourceContentConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_md5", "ea35f0444ea9a3d5641d8760bc2815cc",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha", "019c79c4dc14dbe1edb3e467b2de6a6aad148717",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha256", "3fb55c931a048943b8d7558dde7c2e4bfc8e04be33b1b55691053d1352391fa7",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_base64sha256", "P7VckxoEiUO411WN3nwuS/yOBL4zsbVWkQU9E1I5H6c=",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha512", "57e2d073dce214609bd61113b90b0b2b7c75034047224d56e35f363c8f2662e3acd561eebf94826a67453411181eca7e1cbf15db1f2fdd496cf13df46b7848c3",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_base64sha512", "V+LQc9ziFGCb1hETuQsLK3x1A0BHIk1W4182PI8mYuOs1WHuv5SCamdFNBEYHsp+HL8V2x8v3Uls8T30a3hIww==",
					),
				),
			},
			{
				Config: testAccArchiveFileResourceFileConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_md5", "59fbc9e62af3cbc2f588f97498240dae",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha", "ce4ee1450ab93ac86e11446649e44cea907b6568",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha256", "5131387f97167da47aa741df3ab2c82f182f17c514c222538d34708d04e0756b",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_base64sha256", "UTE4f5cWfaR6p0HfOrLILxgvF8UUwiJTjTRwjQTgdWs=",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha512", "eb33eb0f8cd8efe1a5a0b99acbd22ed22dbebb80817f8de6e8fed15c21c52240838d9bb46fb0938846c74f694425551ba60829a6396f91fcfe49d21a1e3bb409",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_base64sha512", "6zPrD4zY7+GloLmay9Iu0i2+u4CBf43m6P7RXCHFIkCDjZu0b7CTiEbHT2lEJVUbpggppjlvkfz+SdIaHju0CQ==",
					),
				),
			},
			{
				Config: testAccArchiveFileResourceDirConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_md5", "b73f64a383716070aa4a29563b8b14d4",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha", "76d20a402eefd1cfbdc47886abd4e0909616c191",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha256", "c9d07cc2dabc9caf6f43bed51fa613c281e6ca58cae3a8d6fae2094b00b3369a",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_base64sha256", "ydB8wtq8nK9vQ77VH6YTwoHmyljK46jW+uIJSwCzNpo=",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha512", "b96ac6b9554a04473a733be36f190422bf7162b4afdb211a0f551713eadf4092459426750646c70383ce6c8b89171b88582a608e5841bfaaafa17004a2a2ca0a",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_base64sha512", "uWrGuVVKBEc6czvjbxkEIr9xYrSv2yEaD1UXE+rfQJJFlCZ1BkbHA4PObIuJFxuIWCpgjlhBv6qvoXAEoqLKCg==",
					),
				),
			},
			{
				Config: testAccArchiveFileResourceDirExcludesConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
				),
			},
			{
				Config: testAccArchiveFileResourceDirExcludesGlobConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
				),
			},
			{
				Config: testAccArchiveFileResourceMultiSourceConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
				),
			},
		},
	})
}

func TestResource_UpgradeFromVersion2_2_0_ContentConfig(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test_upgrade_content_config.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		Steps: []r.TestStep{
			{
				ExternalProviders: map[string]r.ExternalProvider{
					"archive": {
						VersionConstraint: "2.2.0",
						Source:            "hashicorp/archive",
					},
				},
				Config: testAccArchiveFileResourceContentConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_base64sha256", "P7VckxoEiUO411WN3nwuS/yOBL4zsbVWkQU9E1I5H6c=",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_md5", "ea35f0444ea9a3d5641d8760bc2815cc",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha", "019c79c4dc14dbe1edb3e467b2de6a6aad148717",
					),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileResourceContentConfig(f),
				PlanOnly:                 true,
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileResourceContentConfig(f),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_base64sha256", "P7VckxoEiUO411WN3nwuS/yOBL4zsbVWkQU9E1I5H6c=",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_md5", "ea35f0444ea9a3d5641d8760bc2815cc",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha", "019c79c4dc14dbe1edb3e467b2de6a6aad148717",
					),
				),
			},
		},
	})
}

func TestResource_UpgradeFromVersion2_2_0_FileConfig(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test_upgrade_file_config.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		Steps: []r.TestStep{
			{
				ExternalProviders: map[string]r.ExternalProvider{
					"archive": {
						VersionConstraint: "2.2.0",
						Source:            "hashicorp/archive",
					},
				},
				Config: testAccArchiveFileResourceFileConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_base64sha256", "UTE4f5cWfaR6p0HfOrLILxgvF8UUwiJTjTRwjQTgdWs=",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_md5", "59fbc9e62af3cbc2f588f97498240dae",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha", "ce4ee1450ab93ac86e11446649e44cea907b6568",
					),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileResourceFileConfig(f),
				PlanOnly:                 true,
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileResourceFileConfig(f),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_base64sha256", "UTE4f5cWfaR6p0HfOrLILxgvF8UUwiJTjTRwjQTgdWs=",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_md5", "59fbc9e62af3cbc2f588f97498240dae",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha", "ce4ee1450ab93ac86e11446649e44cea907b6568",
					),
				),
			},
		},
	})
}

func TestResource_UpgradeFromVersion2_2_0_DirConfig(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test_upgrade_dir_config.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		Steps: []r.TestStep{
			{
				ExternalProviders: map[string]r.ExternalProvider{
					"archive": {
						VersionConstraint: "2.2.0",
						Source:            "hashicorp/archive",
					},
				},
				Config: testAccArchiveFileResourceDirConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_base64sha256", "ydB8wtq8nK9vQ77VH6YTwoHmyljK46jW+uIJSwCzNpo=",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_md5", "b73f64a383716070aa4a29563b8b14d4",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha", "76d20a402eefd1cfbdc47886abd4e0909616c191",
					),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileResourceDirConfig(f),
				PlanOnly:                 true,
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileResourceDirConfig(f),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_base64sha256", "ydB8wtq8nK9vQ77VH6YTwoHmyljK46jW+uIJSwCzNpo=",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_md5", "b73f64a383716070aa4a29563b8b14d4",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha", "76d20a402eefd1cfbdc47886abd4e0909616c191",
					),
				),
			},
		},
	})
}

func TestResource_UpgradeFromVersion2_2_0_DirExcludesConfig(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test_upgrade_dir_excludes.zip")

	var fileSize, outputSha string

	r.ParallelTest(t, r.TestCase{
		Steps: []r.TestStep{
			{
				ExternalProviders: map[string]r.ExternalProvider{
					"archive": {
						VersionConstraint: "2.2.0",
						Source:            "hashicorp/archive",
					},
				},
				Config: testAccArchiveFileResourceDirExcludesConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					testExtractResourceAttr("archive_file.foo", "output_sha", &outputSha),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileResourceDirExcludesConfig(f),
				PlanOnly:                 true,
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileResourceDirExcludesConfig(f),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_sha", &outputSha),
				),
			},
		},
	})
}

func TestResource_UpgradeFromVersion2_2_0_SourceConfig(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test_upgrade_source.zip")

	var fileSize, outputSha string

	r.ParallelTest(t, r.TestCase{
		Steps: []r.TestStep{
			{
				ExternalProviders: map[string]r.ExternalProvider{
					"archive": {
						VersionConstraint: "2.2.0",
						Source:            "hashicorp/archive",
					},
				},
				Config: testAccArchiveFileResourceMultiSourceConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					testExtractResourceAttr("archive_file.foo", "output_sha", &outputSha),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileResourceMultiSourceConfig(f),
				PlanOnly:                 true,
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileResourceMultiSourceConfig(f),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_sha", &outputSha),
				),
			},
		},
	})
}

// TestResource_FileConfig_ModifiedContents tests that archive_file resource replaces the resource on every read.
// The contents of the source file are altered, but no aspect of the Terraform configuration is changed.
// The change in the output hashes demonstrates that the resource Read function is replacing the resource.
func TestResource_FileConfig_ModifiedContents(t *testing.T) {
	td := t.TempDir()

	sourceFilePath := filepath.Join(td, "sourceFile")
	outputFilePath := filepath.Join(td, "zip_file_acc_test_upgrade_file_config.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		Steps: []r.TestStep{
			{
				PreConfig: func() {
					alterFileContents("content", sourceFilePath)
				},
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileResourceFileSourceFileConfig(sourceFilePath, outputFilePath),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(outputFilePath, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_base64sha256", "8inOmQJB12dXqCyRTdaRO63yP22Rmuube/A1DLDii10=",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_md5", "20d9c8096f99174d128e5042279fe576",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha", "119d3169ec43fe95bbd38a3824f4e477f4e8d4e7",
					),
				),
			},
			{
				PreConfig: func() {
					alterFileContents("modified content", sourceFilePath)
				},
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileResourceFileSourceFileConfig(sourceFilePath, outputFilePath),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(outputFilePath, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_base64sha256", "OnzXDJ3jda5RPuINpxKHQsZ+jSNOupxShSmW3iUWw7Q=",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_md5", "ce31b13da062764f2975d1ef08ee56fe",
					),
					r.TestCheckResourceAttr(
						"archive_file.foo", "output_sha", "12c51ec24fc5dc10570abbc0b56ac5a3b4141b83",
					),
				),
			},
		},
	})
}

func TestResource_SourceConfigMissing(t *testing.T) {
	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config:      testResourceSourceConfigMissing(),
				ExpectError: regexp.MustCompile(`.*At least one of these attributes must be configured:\n\[source,source_content_filename,source_file,source_dir]`),
			},
		},
	})
}

func TestResource_SourceConfigConflicting(t *testing.T) {
	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config:      testResourceSourceConfigConflicting(),
				ExpectError: regexp.MustCompile(`.*Attribute "source_dir" cannot be specified when "source" is specified`),
			},
		},
	})
}

// TestResource_ArchiveFile_SymlinkFile_Relative verifies that a symlink to a file using a relative path generates an
// archive which includes the file.
func TestResource_ArchiveFile_SymlinkFile_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_SymlinkFile_Absolute verifies that a symlink to a file using an absolute path generates an
// archive which includes the file.
func TestResource_ArchiveFile_SymlinkFile_Absolute(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkFileAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-file/test-symlink.txt")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkFileAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_SymlinkDirectory_Relative verifies that a symlink to a directory using a relative path
// generates an archive which includes the files in the directory.
func TestResource_ArchiveFile_SymlinkDirectory_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-symlink-dir"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"file1.txt": []byte(`This is file 1`),
							"file2.txt": []byte(`This is file 2`),
							"file3.txt": []byte(`This is file 3`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_SymlinkDirectory_Absolute verifies that a symlink to a directory using an absolute path
// generates an archive which includes the files in the directory.
func TestResource_ArchiveFile_SymlinkDirectory_Absolute(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirWithRegularFilesAbs, err := filepath.Abs("test-fixtures/test-symlink-dir")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkDirWithRegularFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"file1.txt": []byte(`This is file 1`),
							"file2.txt": []byte(`This is file 2`),
							"file3.txt": []byte(`This is file 3`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_DirectoryWithSymlinkFile_Relative verifies that a relative path to a directory containing
// a symlink file generates an archive which includes the files in the directory.
func TestResource_ArchiveFile_DirectoryWithSymlinkFile_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-file.txt":    []byte(`This is test content`),
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_DirectoryWithSymlinkFile_Absolute verifies that an absolute path to a directory containing
// a symlink file generates an archive which includes the files in the directory.
func TestResource_ArchiveFile_DirectoryWithSymlinkFile_Absolute(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirWithSymlinkFilesAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-file")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkDirWithSymlinkFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-file.txt":    []byte(`This is test content`),
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative verifies that a relative path to a symlink
// file in a symlink directory generates an archive which includes the files in the directory.
func TestResource_ArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-symlink-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute verifies that an absolute path to a symlink
// file in a symlink directory generates an archive which includes the files in the directory.
func TestResource_ArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkFileInSymlinkDirAbs, err := filepath.Abs("test-fixtures/test-symlink-dir-with-symlink-file/test-symlink.txt")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkFileInSymlinkDirAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_DirectoryWithSymlinkDirectory_Relative verifies that a relative path to a
// directory containing a symlink to a directory generates an archive which includes the directory.
func TestResource_ArchiveFile_DirectoryWithSymlinkDirectory_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-dir"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink-dir/file1.txt": []byte("This is file 1"),
							"test-symlink-dir/file2.txt": []byte("This is file 2"),
							"test-symlink-dir/file3.txt": []byte("This is file 3"),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute verifies that an absolute path to a
// directory containing a symlink to a directory generates an archive which includes the directory.
func TestResource_ArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirInRegularDirAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-dir")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkDirInRegularDirAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink-dir/file1.txt": []byte("This is file 1"),
							"test-symlink-dir/file2.txt": []byte("This is file 2"),
							"test-symlink-dir/file3.txt": []byte("This is file 3"),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_Multiple_Relative verifies that a relative path to a directory containing multiple
// directories including symlink directories generates an archive which includes the directories and files.
func TestResource_ArchiveFile_Multiple_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
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
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_Multiple_Absolute verifies that an absolute path to a directory containing multiple
// directories including symlink directories generates an archive which includes the directories and files.
func TestResource_ArchiveFile_Multiple_Absolute(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	multipleDirsAndFilesAbs, err := filepath.Abs("test-fixtures")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(multipleDirsAndFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
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
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_SymlinkFile_Relative_ExcludeSymlinkDirectories verifies that a symlink to a file using a relative
// path generates an archive which includes the file.
func TestResource_ArchiveFile_SymlinkFile_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_SymlinkFile_Absolute_ExcludeSymlinkDirectories verifies that a symlink to a file using an absolute
// path generates an archive which includes the file.
func TestResource_ArchiveFile_SymlinkFile_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkFileAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-file/test-symlink.txt")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkFileAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_SymlinkDirectory_Relative_ExcludeSymlinkDirectories verifies that an empty archive
// is generated when trying to archive a directory which only contains a symlink to a directory.
func TestResource_ArchiveFile_SymlinkDirectory_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-symlink-dir"), filepath.ToSlash(f)),
				ExpectError: regexp.MustCompile(`.*error creating archive: error archiving directory: archive has not been\ncreated as it would be empty`),
			},
		},
	})
}

// TestResource_ArchiveFile_SymlinkDirectory_Absolute_ExcludeSymlinkDirectories verifies that an empty archive
// is generated when trying to archive a directory which only contains a symlink to a directory.
func TestResource_ArchiveFile_SymlinkDirectory_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirWithRegularFilesAbs, err := filepath.Abs("test-fixtures/test-symlink-dir")
	if err != nil {
		t.Fatal(err)
	}

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkDirWithRegularFilesAbs), filepath.ToSlash(f)),
				ExpectError: regexp.MustCompile(`.*error creating archive: error archiving directory: archive has not been\ncreated as it would be empty`),
			},
		},
	})
}

// TestResource_ArchiveFile_DirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories verifies that a relative path to a
// directory containing a symlink file generates an archive which includes the files in the directory.
func TestResource_ArchiveFile_DirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-file.txt":    []byte(`This is test content`),
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_DirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories verifies that an absolute path to a
// directory containing a symlink file generates an archive which includes the files in the directory.
func TestResource_ArchiveFile_DirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirWithSymlinkFilesAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-file")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkDirWithSymlinkFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-file.txt":    []byte(`This is test content`),
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories verifies that a relative path
// to a symlink file in a symlink directory generates an archive which includes the files in the directory.
func TestResource_ArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-symlink-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories verifies that an absolute path
// to a symlink file in a symlink directory generates an archive which includes the files in the directory.
func TestResource_ArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkFileInSymlinkDirAbs, err := filepath.Abs("test-fixtures/test-symlink-dir-with-symlink-file/test-symlink.txt")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkFileInSymlinkDirAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_DirectoryWithSymlinkDirectory_Relative_ExcludeSymlinkDirectories verifies that an empty archive
// is generated when trying to archive a directory which only contains a symlink to a directory.
func TestResource_ArchiveFile_DirectoryWithSymlinkDirectory_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-dir"), filepath.ToSlash(f)),
				ExpectError: regexp.MustCompile(`.*error creating archive: error archiving directory: archive has not been\ncreated as it would be empty`),
			},
		},
	})
}

// TestResource_ArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute_ExcludeSymlinkDirectories verifies that an empty archive
// is generated when trying to archive a directory which only contains a symlink to a directory.
func TestResource_ArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirInRegularDirAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-dir")
	if err != nil {
		t.Fatal(err)
	}

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkDirInRegularDirAbs), filepath.ToSlash(f)),
				ExpectError: regexp.MustCompile(`.*error creating archive: error archiving directory: archive has not been\ncreated as it would be empty`),
			},
		},
	})
}

// TestResource_ArchiveFile_Multiple_Relative_ExcludeSymlinkDirectories verifies that
// symlinked directories are excluded.
func TestResource_ArchiveFile_Multiple_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-dir/test-dir1/file1.txt":                []byte("This is file 1"),
							"test-dir/test-dir1/file2.txt":                []byte("This is file 2"),
							"test-dir/test-dir1/file3.txt":                []byte("This is file 3"),
							"test-dir/test-dir2/file1.txt":                []byte("This is file 1"),
							"test-dir/test-dir2/file2.txt":                []byte("This is file 2"),
							"test-dir/test-dir2/file3.txt":                []byte("This is file 3"),
							"test-dir/test-file.txt":                      []byte("This is test content"),
							"test-dir-with-symlink-file/test-file.txt":    []byte("This is test content"),
							"test-dir-with-symlink-file/test-symlink.txt": []byte("This is test content"),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestResource_ArchiveFile_Multiple_Relative_ExcludeSymlinkDirectories verifies that
// symlinked directories are excluded.
func TestResource_ArchiveFile_Multiple_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	multipleDirsAndFilesAbs, err := filepath.Abs("test-fixtures")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			resource "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(multipleDirsAndFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-dir/test-dir1/file1.txt":                []byte("This is file 1"),
							"test-dir/test-dir1/file2.txt":                []byte("This is file 2"),
							"test-dir/test-dir1/file3.txt":                []byte("This is file 3"),
							"test-dir/test-dir2/file1.txt":                []byte("This is file 1"),
							"test-dir/test-dir2/file2.txt":                []byte("This is file 2"),
							"test-dir/test-dir2/file3.txt":                []byte("This is file 3"),
							"test-dir/test-file.txt":                      []byte("This is test content"),
							"test-dir-with-symlink-file/test-file.txt":    []byte("This is test content"),
							"test-dir-with-symlink-file/test-symlink.txt": []byte("This is test content"),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

func alterFileContents(content, path string) {
	f, err := os.Create(path)
	if err != nil {
		panic(fmt.Sprintf("error creating file: %s", err))
	}

	defer f.Close()

	_, err = f.Write([]byte(content))
	if err != nil {
		panic(fmt.Sprintf("error writing file: %s", err))
	}
}

func testAccArchiveFileResourceContentConfig(outputPath string) string {
	return fmt.Sprintf(`
resource "archive_file" "foo" {
  type                    = "zip"
  source_content          = "This is some content"
  source_content_filename = "content.txt"
  output_path             = "%s"
}
`, filepath.ToSlash(outputPath))
}

func testAccArchiveFileResourceFileConfig(outputPath string) string {
	return fmt.Sprintf(`
resource "archive_file" "foo" {
  type             = "zip"
  source_file      = "test-fixtures/test-dir/test-file.txt"
  output_path      = "%s"
  output_file_mode = "0666"
}
`, filepath.ToSlash(outputPath))
}

func testAccArchiveFileResourceFileSourceFileConfig(sourceFile, outputPath string) string {
	return fmt.Sprintf(`
resource "archive_file" "foo" {
  type             = "zip"
  source_file      = "%s"
  output_path      = "%s"
  output_file_mode = "0666"
}
`,
		filepath.ToSlash(sourceFile),
		filepath.ToSlash(outputPath))
}

func testAccArchiveFileResourceDirConfig(outputPath string) string {
	return fmt.Sprintf(`
resource "archive_file" "foo" {
  type             = "zip"
  source_dir       = "test-fixtures/test-dir/test-dir1"
  output_path      = "%s"
  output_file_mode = "0666"
}
`, filepath.ToSlash(outputPath))
}

func testAccArchiveFileResourceDirExcludesConfig(outputPath string) string {
	return fmt.Sprintf(`
resource "archive_file" "foo" {
  type        = "zip"
  source_dir  = "test-fixtures/test-dir"
  excludes    = ["test-fixtures/test-dir/file2.txt"]
  output_path = "%s"
}
`, filepath.ToSlash(outputPath))
}

func testAccArchiveFileResourceDirExcludesGlobConfig(outputPath string) string {
	return fmt.Sprintf(`
resource "archive_file" "foo" {
  type        = "zip"
  source_dir  = "test-fixtures/test-dir"
  excludes    = ["test-fixtures/test-dir/file2.txt", "**/file[2-3].txt"]
  output_path = "%s"
}
`, filepath.ToSlash(outputPath))
}

func testAccArchiveFileResourceMultiSourceConfig(outputPath string) string {
	return fmt.Sprintf(`
resource "archive_file" "foo" {
  type = "zip"
  source {
    filename = "content_1.txt"
    content = "This is the content for content_1.txt"
  }
  source {
    filename = "content_2.txt"
    content = "This is the content for content_2.txt"
  }
  output_path = "%s"
}
`, filepath.ToSlash(outputPath))
}

func testResourceSourceConfigMissing() string {
	return `
resource "archive_file" "foo" {
  type                    = "zip"
  output_path             = "path"
}
`
}

func testResourceSourceConfigConflicting() string {
	return `
resource "archive_file" "foo" {
  type                    = "zip"
  source {
    filename = "content_1.txt"
    content = "This is the content for content_1.txt"
  }
  source_dir  = "test-fixtures/test-dir"
  output_path             = "path"
}
`
}
