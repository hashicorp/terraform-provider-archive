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
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string
	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: testAccArchiveFileContentConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_md5", "ea35f0444ea9a3d5641d8760bc2815cc",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha", "019c79c4dc14dbe1edb3e467b2de6a6aad148717",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha256", "3fb55c931a048943b8d7558dde7c2e4bfc8e04be33b1b55691053d1352391fa7",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha256", "P7VckxoEiUO411WN3nwuS/yOBL4zsbVWkQU9E1I5H6c=",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha512", "57e2d073dce214609bd61113b90b0b2b7c75034047224d56e35f363c8f2662e3acd561eebf94826a67453411181eca7e1cbf15db1f2fdd496cf13df46b7848c3",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha512", "V+LQc9ziFGCb1hETuQsLK3x1A0BHIk1W4182PI8mYuOs1WHuv5SCamdFNBEYHsp+HL8V2x8v3Uls8T30a3hIww==",
					),
				),
			},
			{
				Config: testAccArchiveFileFileConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_md5", "59fbc9e62af3cbc2f588f97498240dae",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha", "ce4ee1450ab93ac86e11446649e44cea907b6568",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha256", "5131387f97167da47aa741df3ab2c82f182f17c514c222538d34708d04e0756b",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha256", "UTE4f5cWfaR6p0HfOrLILxgvF8UUwiJTjTRwjQTgdWs=",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha512", "eb33eb0f8cd8efe1a5a0b99acbd22ed22dbebb80817f8de6e8fed15c21c52240838d9bb46fb0938846c74f694425551ba60829a6396f91fcfe49d21a1e3bb409",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha512", "6zPrD4zY7+GloLmay9Iu0i2+u4CBf43m6P7RXCHFIkCDjZu0b7CTiEbHT2lEJVUbpggppjlvkfz+SdIaHju0CQ==",
					),
				),
			},
			{
				Config: testAccArchiveFileDirConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_md5", "b73f64a383716070aa4a29563b8b14d4",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha", "76d20a402eefd1cfbdc47886abd4e0909616c191",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha256", "c9d07cc2dabc9caf6f43bed51fa613c281e6ca58cae3a8d6fae2094b00b3369a",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha256", "ydB8wtq8nK9vQ77VH6YTwoHmyljK46jW+uIJSwCzNpo=",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha512", "b96ac6b9554a04473a733be36f190422bf7162b4afdb211a0f551713eadf4092459426750646c70383ce6c8b89171b88582a608e5841bfaaafa17004a2a2ca0a",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha512", "uWrGuVVKBEc6czvjbxkEIr9xYrSv2yEaD1UXE+rfQJJFlCZ1BkbHA4PObIuJFxuIWCpgjlhBv6qvoXAEoqLKCg==",
					),
				),
			},
			{
				Config: testAccArchiveFileDirExcludesConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},
			{
				Config: testAccArchiveFileMultiSourceConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},
		},
	})
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
				Config: testAccArchiveFileContentConfig(f),
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
				Config:                   testAccArchiveFileContentConfig(f),
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
				Config: testAccArchiveFileFileConfig(f),
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
				Config:                   testAccArchiveFileFileConfig(f),
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
				Config: testAccArchiveFileDirConfig(f),
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
				Config:                   testAccArchiveFileDirConfig(f),
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
				Config: testAccArchiveFileDirExcludesConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					testExtractResourceAttr("data.archive_file.foo", "output_sha", &outputSha),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileDirExcludesConfig(f),
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
				Config: testAccArchiveFileMultiSourceConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					testExtractResourceAttr("data.archive_file.foo", "output_sha", &outputSha),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccArchiveFileMultiSourceConfig(f),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_sha", &outputSha),
				),
			},
		},
	})
}

func TestAccArchiveFile_SourceConfigMissing(t *testing.T) {
	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config:      testAccArchiveSourceConfigMissing(),
				ExpectError: regexp.MustCompile(`.*At least one of these attributes must be configured:\n\[source,source_content_filename,source_file,source_dir\]`),
			},
		},
	})
}

func TestAccArchiveFile_SourceConfigConflicting(t *testing.T) {
	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config:      testAccArchiveSourceConfigConflicting(),
				ExpectError: regexp.MustCompile(`.*Attribute "source_dir" cannot be specified when "source" is specified`),
			},
		},
	})
}

func TestAccArchiveFile_Symlinks(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkFile := "test-fixtures/test-dir-with-symlink-file/test-symlink.txt"
	symlinkFileAbs, err := filepath.Abs(symlinkFile)
	if err != nil {
		t.Fatal(err)
	}

	symlinkDirWithRegularFiles := "test-fixtures/test-symlink-dir"
	symlinkDirWithRegularFilesAbs, err := filepath.Abs(symlinkDirWithRegularFiles)
	if err != nil {
		t.Fatal(err)
	}

	symlinkDirWithSymlinkFiles := "test-fixtures/test-dir-with-symlink-file"
	symlinkDirWithSymlinkFilesAbs, err := filepath.Abs(symlinkDirWithSymlinkFiles)
	if err != nil {
		t.Fatal(err)
	}

	symlinkFileInSymlinkDir := "test-fixtures/test-symlink-dir-with-symlink-file/test-symlink.txt"
	symlinkFileInSymlinkDirAbs, err := filepath.Abs(symlinkFileInSymlinkDir)
	if err != nil {
		t.Fatal(err)
	}

	symlinkDirInRegularDir := "test-fixtures/test-dir-with-symlink-dir"
	symlinkDirInRegularDirAbs, err := filepath.Abs(symlinkDirInRegularDir)
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
			`, filepath.ToSlash(symlinkFile), filepath.ToSlash(f)),
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
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkDirWithRegularFiles), filepath.ToSlash(f)),
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
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkDirWithSymlinkFiles), filepath.ToSlash(f)),
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
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkFileInSymlinkDir), filepath.ToSlash(f)),
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
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkDirInRegularDir), filepath.ToSlash(f)),
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
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, symlinkDirInRegularDirAbs, filepath.ToSlash(f)),
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

func testAccArchiveFileContentConfig(outputPath string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
  type                    = "zip"
  source_content          = "This is some content"
  source_content_filename = "content.txt"
  output_path             = "%s"
}
`, filepath.ToSlash(outputPath))
}

func testAccArchiveFileFileConfig(outputPath string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
  type             = "zip"
  source_file      = "test-fixtures/test-dir/test-file.txt"
  output_path      = "%s"
  output_file_mode = "0666"
}
`, filepath.ToSlash(outputPath))
}

func testAccArchiveFileDirConfig(outputPath string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
  type             = "zip"
  source_dir       = "test-fixtures/test-dir/test-dir1"
  output_path      = "%s"
  output_file_mode = "0666"
}
`, filepath.ToSlash(outputPath))
}

func testAccArchiveFileDirExcludesConfig(outputPath string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
  type        = "zip"
  source_dir  = "test-fixtures/test-dir/test-dir1"
  excludes    = ["test-fixtures/test-dir/test-dir1/file2.txt"]
  output_path = "%s"
}
`, filepath.ToSlash(outputPath))
}

func testAccArchiveFileMultiSourceConfig(outputPath string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
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

func testAccArchiveSourceConfigMissing() string {
	return `
data "archive_file" "foo" {
  type                    = "zip"
  output_path             = "path"
}
`
}

func testAccArchiveSourceConfigConflicting() string {
	return `
data "archive_file" "foo" {
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
