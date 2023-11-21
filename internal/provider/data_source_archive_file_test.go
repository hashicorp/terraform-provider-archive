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
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccArchiveFile_Basic(t *testing.T) {
	hashes := map[string]map[string]map[string]string{
		"content.txt": {
			"zip": {
				"md5":          "ea35f0444ea9a3d5641d8760bc2815cc",
				"sha":          "019c79c4dc14dbe1edb3e467b2de6a6aad148717",
				"sha256":       "3fb55c931a048943b8d7558dde7c2e4bfc8e04be33b1b55691053d1352391fa7",
				"base64sha256": "P7VckxoEiUO411WN3nwuS/yOBL4zsbVWkQU9E1I5H6c=",
				"sha512":       "57e2d073dce214609bd61113b90b0b2b7c75034047224d56e35f363c8f2662e3acd561eebf94826a67453411181eca7e1cbf15db1f2fdd496cf13df46b7848c3",
				"base64sha512": "V+LQc9ziFGCb1hETuQsLK3x1A0BHIk1W4182PI8mYuOs1WHuv5SCamdFNBEYHsp+HL8V2x8v3Uls8T30a3hIww==",
			},
			"tar.gz": {
				"md5":          "a09ee39e708c38ccd9ba44cc39e7cacc",
				"sha":          "6c84af188d367644731196007301c9dc93914b0e",
				"sha256":       "ae18ec27af576dd62f29cec7ae0df130e7487c1a3cddefdec9f27d5ed3a4ca95",
				"base64sha256": "rhjsJ69XbdYvKc7Hrg3xMOdIfBo83e/eyfJ9XtOkypU=",
				"sha512":       "ce93a8ba072bad42656a41def0c9ed160f9109c1a7087fb0dcf0b9fce9effc25477f9cdbf9cbc5aa593f3ded0e0db11d2c8cf67dc8d2693ff4069aa01071e68d",
				"base64sha512": "zpOougcrrUJlakHe8MntFg+RCcGnCH+w3PC5/Onv/CVHf5zb+cvFqlk/Pe0ODbEdLIz2fcjSaT/0BpqgEHHmjQ==",
			},
		},
		"test-file.txt": {
			"zip": {
				"md5":          "59fbc9e62af3cbc2f588f97498240dae",
				"sha":          "ce4ee1450ab93ac86e11446649e44cea907b6568",
				"sha256":       "5131387f97167da47aa741df3ab2c82f182f17c514c222538d34708d04e0756b",
				"base64sha256": "UTE4f5cWfaR6p0HfOrLILxgvF8UUwiJTjTRwjQTgdWs=",
				"sha512":       "eb33eb0f8cd8efe1a5a0b99acbd22ed22dbebb80817f8de6e8fed15c21c52240838d9bb46fb0938846c74f694425551ba60829a6396f91fcfe49d21a1e3bb409",
				"base64sha512": "6zPrD4zY7+GloLmay9Iu0i2+u4CBf43m6P7RXCHFIkCDjZu0b7CTiEbHT2lEJVUbpggppjlvkfz+SdIaHju0CQ==",
			},
			"tar.gz": {
				"md5":          "39948f8ddedc8914ac2e42dd18dd3c06",
				"sha":          "d2f26a69cbb920715f797f81c1477c41d8fc9195",
				"sha256":       "6f1b1a5d17e42fae154f0bdf9301a0ad43394d7fe8485b64dfdb533a0cf07784",
				"base64sha256": "bxsaXRfkL64VTwvfkwGgrUM5TX/oSFtk39tTOgzwd4Q=",
				"sha512":       "e0dc636c3d11095a5b5e97c598185ee3d8a0ed7d1accb69cc28419aeeaeda22b2e774a260f71892a2e85efae1a3aee36669b61dafaed9ac0886d2ca8c5add6e9",
				"base64sha512": "4NxjbD0RCVpbXpfFmBhe49ig7X0azLacwoQZrurtoisud0omD3GJKi6F764aOu42Zpth2vrtmsCIbSyoxa3W6Q==",
			},
		},
		"test-dir1": {
			"zip": {
				"md5":          "b73f64a383716070aa4a29563b8b14d4",
				"sha":          "76d20a402eefd1cfbdc47886abd4e0909616c191",
				"sha256":       "c9d07cc2dabc9caf6f43bed51fa613c281e6ca58cae3a8d6fae2094b00b3369a",
				"base64sha256": "ydB8wtq8nK9vQ77VH6YTwoHmyljK46jW+uIJSwCzNpo=",
				"sha512":       "b96ac6b9554a04473a733be36f190422bf7162b4afdb211a0f551713eadf4092459426750646c70383ce6c8b89171b88582a608e5841bfaaafa17004a2a2ca0a",
				"base64sha512": "uWrGuVVKBEc6czvjbxkEIr9xYrSv2yEaD1UXE+rfQJJFlCZ1BkbHA4PObIuJFxuIWCpgjlhBv6qvoXAEoqLKCg==",
			},
			"tar.gz": {
				"md5":          "56d3b5da9e46596a7037008224a3535d",
				"sha":          "d1934f7c5cd325b2c1b3128621e3041c931d5706",
				"sha256":       "0496affbe2b4dcde9cd3e49665b351a1ef7633faa35bb908de617b155f74f123",
				"base64sha256": "BJav++K03N6c0+SWZbNRoe92M/qjW7kI3mF7FV908SM=",
				"sha512":       "5ddacbb55be1dc431cef74ab65ad8059a51990eb4cf108f679b79a59004a586cdd0b1bb0d869fa1fe679fa71a1d85391c19413a2d79050df68703aa3390eabee",
				"base64sha512": "XdrLtVvh3EMc73SrZa2AWaUZkOtM8Qj2ebeaWQBKWGzdCxuw2Gn6H+Z5+nGh2FORwZQToteQUN9ocDqjOQ6r7g==",
			},
		},
	}

	for _, format := range []string{"zip", "tar.gz"} {
		t.Run(format, func(t *testing.T) {
			td := t.TempDir()

			f := filepath.Join(td, "file_acc_test."+format)

			var fileSize string
			r.ParallelTest(t, r.TestCase{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Steps: []r.TestStep{
					{
						Config: testAccArchiveFileContentConfig(format, f),
						Check: r.ComposeTestCheckFunc(
							testAccArchiveFileSize(f, &fileSize),
							r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_md5", hashes["content.txt"][format]["md5"],
							),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_sha", hashes["content.txt"][format]["sha"],
							),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_sha256", hashes["content.txt"][format]["sha256"],
							),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_base64sha256", hashes["content.txt"][format]["base64sha256"],
							),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_sha512", hashes["content.txt"][format]["sha512"],
							),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_base64sha512", hashes["content.txt"][format]["base64sha512"],
							),
						),
					},
					{
						Config: testAccArchiveFileFileConfig(format, f),
						Check: r.ComposeTestCheckFunc(
							testAccArchiveFileSize(f, &fileSize),
							r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_md5", hashes["test-file.txt"][format]["md5"],
							),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_sha", hashes["test-file.txt"][format]["sha"],
							),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_sha256", hashes["test-file.txt"][format]["sha256"],
							),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_base64sha256", hashes["test-file.txt"][format]["base64sha256"],
							),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_sha512", hashes["test-file.txt"][format]["sha512"],
							),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_base64sha512", hashes["test-file.txt"][format]["base64sha512"],
							),
						),
					},
					{
						Config: testAccArchiveFileDirConfig(format, f),
						Check: r.ComposeTestCheckFunc(
							testAccArchiveFileSize(f, &fileSize),
							r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_md5", hashes["test-dir1"][format]["md5"],
							),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_sha", hashes["test-dir1"][format]["sha"],
							),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_sha256", hashes["test-dir1"][format]["sha256"],
							),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_base64sha256", hashes["test-dir1"][format]["base64sha256"],
							),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_sha512", hashes["test-dir1"][format]["sha512"],
							),
							r.TestCheckResourceAttr(
								"data.archive_file.foo", "output_base64sha512", hashes["test-dir1"][format]["base64sha512"],
							),
						),
					},
					{
						Config: testAccArchiveFileDirExcludesConfig(format, f),
						Check: r.ComposeTestCheckFunc(
							testAccArchiveFileSize(f, &fileSize),
							r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
						),
					},
					{
						Config: testAccArchiveFileMultiSourceConfig(format, f),
						Check: r.ComposeTestCheckFunc(
							testAccArchiveFileSize(f, &fileSize),
							r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
						),
					},
				},
			})
		})
	}
}

func TestDataSource_UpgradeFromVersion2_2_0_ContentConfig(t *testing.T) {
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
				Config: testAccArchiveFileContentConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha256", "P7VckxoEiUO411WN3nwuS/yOBL4zsbVWkQU9E1I5H6c=",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_md5", "ea35f0444ea9a3d5641d8760bc2815cc",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha", "019c79c4dc14dbe1edb3e467b2de6a6aad148717",
					),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileContentConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha256", "P7VckxoEiUO411WN3nwuS/yOBL4zsbVWkQU9E1I5H6c=",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_md5", "ea35f0444ea9a3d5641d8760bc2815cc",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha", "019c79c4dc14dbe1edb3e467b2de6a6aad148717",
					),
				),
			},
		},
	})
}

func TestDataSource_UpgradeFromVersion2_2_0_FileConfig(t *testing.T) {
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
				Config: testAccArchiveFileFileConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha256", "UTE4f5cWfaR6p0HfOrLILxgvF8UUwiJTjTRwjQTgdWs=",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_md5", "59fbc9e62af3cbc2f588f97498240dae",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha", "ce4ee1450ab93ac86e11446649e44cea907b6568",
					),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileFileConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha256", "UTE4f5cWfaR6p0HfOrLILxgvF8UUwiJTjTRwjQTgdWs=",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_md5", "59fbc9e62af3cbc2f588f97498240dae",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha", "ce4ee1450ab93ac86e11446649e44cea907b6568",
					),
				),
			},
		},
	})
}

func TestDataSource_UpgradeFromVersion2_2_0_DirConfig(t *testing.T) {
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
				Config: testAccArchiveFileDirConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha256", "ydB8wtq8nK9vQ77VH6YTwoHmyljK46jW+uIJSwCzNpo=",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_md5", "b73f64a383716070aa4a29563b8b14d4",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha", "76d20a402eefd1cfbdc47886abd4e0909616c191",
					),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileDirConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha256", "ydB8wtq8nK9vQ77VH6YTwoHmyljK46jW+uIJSwCzNpo=",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_md5", "b73f64a383716070aa4a29563b8b14d4",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha", "76d20a402eefd1cfbdc47886abd4e0909616c191",
					),
				),
			},
		},
	})
}

func TestDataSource_UpgradeFromVersion2_2_0_DirExcludesConfig(t *testing.T) {
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
				Config: testAccArchiveFileDirExcludesConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					testExtractResourceAttr("data.archive_file.foo", "output_sha", &outputSha),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileDirExcludesConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_sha", &outputSha),
				),
			},
		},
	})
}

func TestDataSource_UpgradeFromVersion2_2_0_SourceConfig(t *testing.T) {
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
				Config: testAccArchiveFileMultiSourceConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					testExtractResourceAttr("data.archive_file.foo", "output_sha", &outputSha),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileMultiSourceConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_sha", &outputSha),
				),
			},
		},
	})
}

func TestAccArchiveFile_SourceConfigMissing(t *testing.T) {
	for _, format := range []string{"zip", "tar.gz"} {
		t.Run(format, func(t *testing.T) {
			r.ParallelTest(t, r.TestCase{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Steps: []r.TestStep{
					{
						Config:      testAccArchiveSourceConfigMissing(format),
						ExpectError: regexp.MustCompile(`.*At least one of these attributes must be configured:\n\[source,source_content_filename,source_file,source_dir\]`),
					},
				},
			})
		})
	}
}

func TestAccArchiveFile_SourceConfigConflicting(t *testing.T) {
	for _, format := range []string{"zip", "tar.gz"} {
		t.Run(format, func(t *testing.T) {
			r.ParallelTest(t, r.TestCase{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Steps: []r.TestStep{
					{
						Config:      testAccArchiveSourceConfigConflicting(format),
						ExpectError: regexp.MustCompile(`.*Attribute "source_dir" cannot be specified when "source" is specified`),
					},
				},
			})
		})
	}
}

// TestAccArchiveFile_SymlinkFile_Relative verifies that a symlink to a file using a relative path generates an
// archive which includes the file.
func TestAccArchiveFile_SymlinkFile_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_SymlinkFile_Absolute verifies that a symlink to a file using an absolute path generates an
// archive which includes the file.
func TestAccArchiveFile_SymlinkFile_Absolute(t *testing.T) {
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
			data "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkFileAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_SymlinkDirectory_Relative verifies that a symlink to a directory using a relative path
// generates an archive which includes the files in the directory.
func TestAccArchiveFile_SymlinkDirectory_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-symlink-dir"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_SymlinkDirectory_Absolute verifies that a symlink to a directory using an absolute path
// generates an archive which includes the files in the directory.
func TestAccArchiveFile_SymlinkDirectory_Absolute(t *testing.T) {
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
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkDirWithRegularFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_DirectoryWithSymlinkFile_Relative verifies that a relative path to a directory containing
// a symlink file generates an archive which includes the files in the directory.
func TestAccArchiveFile_DirectoryWithSymlinkFile_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_DirectoryWithSymlinkFile_Absolute verifies that an absolute path to a directory containing
// a symlink file generates an archive which includes the files in the directory.
func TestAccArchiveFile_DirectoryWithSymlinkFile_Absolute(t *testing.T) {
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
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkDirWithSymlinkFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative verifies that a relative path to a symlink
// file in a symlink directory generates an archive which includes the files in the directory.
func TestAccArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-symlink-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute verifies that an absolute path to a symlink
// file in a symlink directory generates an archive which includes the files in the directory.
func TestAccArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute(t *testing.T) {
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
			data "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkFileInSymlinkDirAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_DirectoryWithSymlinkDirectory_Relative verifies that a relative path to a
// directory containing a symlink to a directory generates an archive which includes the directory.
func TestAccArchiveFile_DirectoryWithSymlinkDirectory_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-dir"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute verifies that an absolute path to a
// directory containing a symlink to a directory generates an archive which includes the directory.
func TestAccArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute(t *testing.T) {
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
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkDirInRegularDirAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_Multiple_Relative verifies that a relative path to a directory containing multiple
// directories including symlink directories generates an archive which includes the directories and files.
func TestAccArchiveFile_Multiple_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_Multiple_Absolute verifies that an absolute path to a directory containing multiple
// directories including symlink directories generates an archive which includes the directories and files.
func TestAccArchiveFile_Multiple_Absolute(t *testing.T) {
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
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(multipleDirsAndFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_SymlinkFile_Relative_ExcludeSymlinkDirectories verifies that a symlink to a file using a relative
// path generates an archive which includes the file.
func TestAccArchiveFile_SymlinkFile_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_SymlinkFile_Absolute_ExcludeSymlinkDirectories verifies that a symlink to a file using an absolute
// path generates an archive which includes the file.
func TestAccArchiveFile_SymlinkFile_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
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
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkFileAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_SymlinkDirectory_Relative_ExcludeSymlinkDirectories verifies that an error is generated when
// trying to use a symlink to a directory.
func TestAccArchiveFile_SymlinkDirectory_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-symlink-dir"), filepath.ToSlash(f)),
				ExpectError: regexp.MustCompile(`.*error creating archive: error archiving directory: error reading file for`),
			},
		},
	})
}

// TestAccArchiveFile_SymlinkDirectory_Absolute_ExcludeSymlinkDirectories verifies that an error is generated when
// trying to use a symlink to a directory.
func TestAccArchiveFile_SymlinkDirectory_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
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
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkDirWithRegularFilesAbs), filepath.ToSlash(f)),
				ExpectError: regexp.MustCompile(`.*error creating archive: error archiving directory: error reading file for`),
			},
		},
	})
}

// TestAccArchiveFile_DirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories verifies that a relative path to a
// directory containing a symlink file generates an archive which includes the files in the directory.
func TestAccArchiveFile_DirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_DirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories verifies that an absolute path to a
// directory containing a symlink file generates an archive which includes the files in the directory.
func TestAccArchiveFile_DirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
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
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkDirWithSymlinkFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories verifies that a relative path
// to a symlink file in a symlink directory generates an archive which includes the files in the directory.
func TestAccArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-symlink-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories verifies that an absolute path
// to a symlink file in a symlink directory generates an archive which includes the files in the directory.
func TestAccArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
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
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkFileInSymlinkDirAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
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

// TestAccArchiveFile_DirectoryWithSymlinkDirectory_Relative_ExcludeSymlinkDirectories verifies that an error is
// generated when trying to a directory which contains a symlink to a directory.
func TestAccArchiveFile_DirectoryWithSymlinkDirectory_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-dir"), filepath.ToSlash(f)),
				ExpectError: regexp.MustCompile(`.*error creating archive: error archiving directory: error reading file for`),
			},
		},
	})
}

// TestAccArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute_ExcludeSymlinkDirectories verifies that an error is
// generated when trying to a directory which contains a symlink to a directory.
func TestAccArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
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
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkDirInRegularDirAbs), filepath.ToSlash(f)),
				ExpectError: regexp.MustCompile(`.*error creating archive: error archiving directory: error reading file for`),
			},
		},
	})
}

// TestAccArchiveFile_Multiple_Relative_ExcludeSymlinkDirectories verifies that an error is
// generated when trying to a directory which contains a symlink to a directory.
func TestAccArchiveFile_Multiple_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures"), filepath.ToSlash(f)),
				ExpectError: regexp.MustCompile(`.*error creating archive: error archiving directory: error reading file for`),
			},
		},
	})
}

// TestAccArchiveFile_Multiple_Absolute_ExcludeSymlinkDirectories verifies that an error is
// generated when trying to a directory which contains a symlink to a directory.
func TestAccArchiveFile_Multiple_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	multipleDirsAndFilesAbs, err := filepath.Abs("test-fixtures")
	if err != nil {
		t.Fatal(err)
	}

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(multipleDirsAndFilesAbs), filepath.ToSlash(f)),
				ExpectError: regexp.MustCompile(`.*error creating archive: error archiving directory: error reading file for`),
			},
		},
	})
}

func testAccArchiveFileSize(filename string, fileSize *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		*fileSize = ""
		fi, err := os.Stat(filename)
		if err != nil {
			return err
		}
		*fileSize = fmt.Sprintf("%d", fi.Size())
		return nil
	}
}

func testAccArchiveFileContentConfig(format, outputPath string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
  type                    = "%s"
  source_content          = "This is some content"
  source_content_filename = "content.txt"
  output_path             = "%s"
}
`, format, filepath.ToSlash(outputPath))
}

func testAccArchiveFileFileConfig(format, outputPath string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
  type             = "%s"
  source_file      = "test-fixtures/test-dir/test-file.txt"
  output_path      = "%s"
  output_file_mode = "0666"
}
`, format, filepath.ToSlash(outputPath))
}

func testAccArchiveFileDirConfig(format, outputPath string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
  type             = "%s"
  source_dir       = "test-fixtures/test-dir/test-dir1"
  output_path      = "%s"
  output_file_mode = "0666"
}
`, format, filepath.ToSlash(outputPath))
}

func testAccArchiveFileDirExcludesConfig(format, outputPath string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
  type        = "%s"
  source_dir  = "test-fixtures/test-dir/test-dir1"
  excludes    = ["test-fixtures/test-dir/test-dir1/file2.txt"]
  output_path = "%s"
}
`, format, filepath.ToSlash(outputPath))
}

func testAccArchiveFileMultiSourceConfig(format, outputPath string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
  type = "%s"
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
`, format, filepath.ToSlash(outputPath))
}

func testAccArchiveSourceConfigMissing(format string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
  type                    = "%s"
  output_path             = "path"
}
`, format)
}

func testAccArchiveSourceConfigConflicting(format string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
  type                    = "%s"
  source {
    filename = "content_1.txt"
    content = "This is the content for content_1.txt"
  }
  source_dir  = "test-fixtures/test-dir"
  output_path             = "path"
}
`, format)
}

//nolint:unparam
func testExtractResourceAttr(resourceName string, attributeName string, attributeValue *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource name %s not found in state", resourceName)
		}

		attrValue, ok := rs.Primary.Attributes[attributeName]
		if !ok {
			return fmt.Errorf("attribute %s not found in resource %s state", attributeName, resourceName)
		}

		*attributeValue = attrValue

		return nil
	}
}
