// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccZipEphemeralArchiveFile_Basic(t *testing.T) {
	td := t.TempDir()
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string
	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.10.0"))), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`

ephemeral "archive_file" "foo" {
  type                    = "%s"
  source_content          = "This is some content"
  source_content_filename = "content.txt"
  output_path             = "%s"
}
`, "zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttr(
						"echo.test1", "data.output_md5", "ea35f0444ea9a3d5641d8760bc2815cc",
					),
					r.TestCheckResourceAttr(
						"echo.test1", "data.output_sha", "019c79c4dc14dbe1edb3e467b2de6a6aad148717",
					),
					r.TestCheckResourceAttr(
						"echo.test1", "data.output_sha256", "3fb55c931a048943b8d7558dde7c2e4bfc8e04be33b1b55691053d1352391fa7",
					),
					r.TestCheckResourceAttr(
						"echo.test1", "data.output_base64sha256", "P7VckxoEiUO411WN3nwuS/yOBL4zsbVWkQU9E1I5H6c=",
					),
					r.TestCheckResourceAttr(
						"echo.test1", "data.output_sha512", "57e2d073dce214609bd61113b90b0b2b7c75034047224d56e35f363c8f2662e3acd561eebf94826a67453411181eca7e1cbf15db1f2fdd496cf13df46b7848c3",
					),
					r.TestCheckResourceAttr(
						"echo.test1", "data.output_base64sha512", "V+LQc9ziFGCb1hETuQsLK3x1A0BHIk1W4182PI8mYuOs1WHuv5SCamdFNBEYHsp+HL8V2x8v3Uls8T30a3hIww==",
					),
				),
			},
			{
				Config: addEchoConfig("test2") + testAccArchiveFileFileConfig("ephemeral", "zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test2", "data.output_size", &fileSize),
					r.TestCheckResourceAttr(
						"echo.test2", "data.output_md5", "59fbc9e62af3cbc2f588f97498240dae",
					),
					r.TestCheckResourceAttr(
						"echo.test2", "data.output_sha", "ce4ee1450ab93ac86e11446649e44cea907b6568",
					),
					r.TestCheckResourceAttr(
						"echo.test2", "data.output_sha256", "5131387f97167da47aa741df3ab2c82f182f17c514c222538d34708d04e0756b",
					),
					r.TestCheckResourceAttr(
						"echo.test2", "data.output_base64sha256", "UTE4f5cWfaR6p0HfOrLILxgvF8UUwiJTjTRwjQTgdWs=",
					),
					r.TestCheckResourceAttr(
						"echo.test2", "data.output_sha512", "eb33eb0f8cd8efe1a5a0b99acbd22ed22dbebb80817f8de6e8fed15c21c52240838d9bb46fb0938846c74f694425551ba60829a6396f91fcfe49d21a1e3bb409",
					),
					r.TestCheckResourceAttr(
						"echo.test2", "data.output_base64sha512", "6zPrD4zY7+GloLmay9Iu0i2+u4CBf43m6P7RXCHFIkCDjZu0b7CTiEbHT2lEJVUbpggppjlvkfz+SdIaHju0CQ==",
					),
				),
			},
			{
				Config: addEchoConfig("test3") + testAccArchiveFileDirConfig("ephemeral", "zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test3", "data.output_size", &fileSize),
					r.TestCheckResourceAttr(
						"echo.test3", "data.output_md5", "b73f64a383716070aa4a29563b8b14d4",
					),
					r.TestCheckResourceAttr(
						"echo.test3", "data.output_sha", "76d20a402eefd1cfbdc47886abd4e0909616c191",
					),
					r.TestCheckResourceAttr(
						"echo.test3", "data.output_sha256", "c9d07cc2dabc9caf6f43bed51fa613c281e6ca58cae3a8d6fae2094b00b3369a",
					),
					r.TestCheckResourceAttr(
						"echo.test3", "data.output_base64sha256", "ydB8wtq8nK9vQ77VH6YTwoHmyljK46jW+uIJSwCzNpo=",
					),
					r.TestCheckResourceAttr(
						"echo.test3", "data.output_sha512", "b96ac6b9554a04473a733be36f190422bf7162b4afdb211a0f551713eadf4092459426750646c70383ce6c8b89171b88582a608e5841bfaaafa17004a2a2ca0a",
					),
					r.TestCheckResourceAttr(
						"echo.test3", "data.output_base64sha512", "uWrGuVVKBEc6czvjbxkEIr9xYrSv2yEaD1UXE+rfQJJFlCZ1BkbHA4PObIuJFxuIWCpgjlhBv6qvoXAEoqLKCg==",
					),
				),
			},
			{
				Config: addEchoConfig("test4") + testAccArchiveFileDirExcludesConfig("ephemeral", "zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test4", "data.output_size", &fileSize),
				),
			},
			{
				Config: addEchoConfig("test5") + testAccArchiveFileDirExcludesGlobConfig("ephemeral", "zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test5", "data.output_size", &fileSize),
				),
			},
			{
				Config: addEchoConfig("test6") + fmt.Sprintf(`
			%s "archive_file" "foo" {
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
			`, "ephemeral", "zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test6", "data.output_size", &fileSize),
				),
			},
		},
	})
}

func TestAccZipEphemeralArchiveFile_SourceConfigMissing(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config:      addEchoConfig("test1") + testAccArchiveSourceConfigMissing("ephemeral", "zip"),
				ExpectError: regexp.MustCompile(`.*At least one of these attributes must be configured:\n\[source,source_content_filename,source_file,source_dir]`),
			},
		},
	})
}

func TestAccZipEphemeralArchiveFile_SourceConfigConflicting(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config:      testAccArchiveSourceConfigConflicting("ephemeral", "zip"),
				ExpectError: regexp.MustCompile(`.*Attribute "source_dir" cannot be specified when "source" is specified`),
			},
		},
	})
}

// TestAccZipEphemeralArchiveFile_SymlinkFile_Relative verifies that a symlink to a file using a relative path generates an
// archive which includes the file.
func TestAccZipEphemeralArchiveFile_SymlinkFile_Relative(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_SymlinkFile_Absolute verifies that a symlink to a file using an absolute path generates an
// archive which includes the file.
func TestAccZipEphemeralArchiveFile_SymlinkFile_Absolute(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkFileAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-file/test-symlink.txt")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkFileAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_SymlinkDirectory_Relative verifies that a symlink to a directory using a relative path
// generates an archive which includes the files in the directory.
func TestAccZipEphemeralArchiveFile_SymlinkDirectory_Relative(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-symlink-dir"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_SymlinkDirectory_Absolute verifies that a symlink to a directory using an absolute path
// generates an archive which includes the files in the directory.
func TestAccZipEphemeralArchiveFile_SymlinkDirectory_Absolute(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirWithRegularFilesAbs, err := filepath.Abs("test-fixtures/test-symlink-dir")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkDirWithRegularFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_DirectoryWithSymlinkFile_Relative verifies that a relative path to a directory containing
// a symlink file generates an archive which includes the files in the directory.
func TestAccZipEphemeralArchiveFile_DirectoryWithSymlinkFile_Relative(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_DirectoryWithSymlinkFile_Absolute verifies that an absolute path to a directory containing
// a symlink file generates an archive which includes the files in the directory.
func TestAccZipEphemeralArchiveFile_DirectoryWithSymlinkFile_Absolute(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirWithSymlinkFilesAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-file")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkDirWithSymlinkFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative verifies that a relative path to a symlink
// file in a symlink directory generates an archive which includes the files in the directory.
func TestAccZipEphemeralArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-symlink-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute verifies that an absolute path to a symlink
// file in a symlink directory generates an archive which includes the files in the directory.
func TestAccZipEphemeralArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkFileInSymlinkDirAbs, err := filepath.Abs("test-fixtures/test-symlink-dir-with-symlink-file/test-symlink.txt")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkFileInSymlinkDirAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_DirectoryWithSymlinkDirectory_Relative verifies that a relative path to a
// directory containing a symlink to a directory generates an archive which includes the directory.
func TestAccZipEphemeralArchiveFile_DirectoryWithSymlinkDirectory_Relative(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-dir"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute verifies that an absolute path to a
// directory containing a symlink to a directory generates an archive which includes the directory.
func TestAccZipEphemeralArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirInRegularDirAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-dir")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkDirInRegularDirAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_Multiple_Relative verifies that a relative path to a directory containing multiple
// directories including symlink directories generates an archive which includes the directories and files.
func TestAccZipEphemeralArchiveFile_Multiple_Relative(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_Multiple_Absolute verifies that an absolute path to a directory containing multiple
// directories including symlink directories generates an archive which includes the directories and files.
func TestAccZipEphemeralArchiveFile_Multiple_Absolute(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	multipleDirsAndFilesAbs, err := filepath.Abs("test-fixtures")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(multipleDirsAndFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_SymlinkFile_Relative_ExcludeSymlinkDirectories verifies that a symlink to a file using a relative
// path generates an archive which includes the file.
func TestAccZipEphemeralArchiveFile_SymlinkFile_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_SymlinkFile_Absolute_ExcludeSymlinkDirectories verifies that a symlink to a file using an absolute
// path generates an archive which includes the file.
func TestAccZipEphemeralArchiveFile_SymlinkFile_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkFileAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-file/test-symlink.txt")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkFileAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_SymlinkDirectory_Relative_ExcludeSymlinkDirectories verifies that an empty archive
// is generated when trying to archive a directory which only contains a symlink to a directory.
func TestAccZipEphemeralArchiveFile_SymlinkDirectory_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
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

// TestAccZipEphemeralArchiveFile_SymlinkDirectory_Absolute_ExcludeSymlinkDirectories verifies that an empty archive
// is generated when trying to archive a directory which only contains a symlink to a directory.
func TestAccZipEphemeralArchiveFile_SymlinkDirectory_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirWithRegularFilesAbs, err := filepath.Abs("test-fixtures/test-symlink-dir")
	if err != nil {
		t.Fatal(err)
	}

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
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

// TestAccZipEphemeralArchiveFile_DirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories verifies that a relative path to a
// directory containing a symlink file generates an archive which includes the files in the directory.
func TestAccZipEphemeralArchiveFile_DirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_DirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories verifies that an absolute path to a
// directory containing a symlink file generates an archive which includes the files in the directory.
func TestAccZipEphemeralArchiveFile_DirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirWithSymlinkFilesAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-file")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkDirWithSymlinkFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories verifies that a relative path
// to a symlink file in a symlink directory generates an archive which includes the files in the directory.
func TestAccZipEphemeralArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-symlink-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories verifies that an absolute path
// to a symlink file in a symlink directory generates an archive which includes the files in the directory.
func TestAccZipEphemeralArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkFileInSymlinkDirAbs, err := filepath.Abs("test-fixtures/test-symlink-dir-with-symlink-file/test-symlink.txt")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkFileInSymlinkDirAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_DirectoryWithSymlinkDirectory_Relative_ExcludeSymlinkDirectories verifies that an empty archive
// is generated when trying to archive a directory which only contains a symlink to a directory.
func TestAccZipEphemeralArchiveFile_DirectoryWithSymlinkDirectory_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
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

// TestAccZipEphemeralArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute_ExcludeSymlinkDirectories verifies that an empty archive
// is generated when trying to archive a directory which only contains a symlink to a directory.
func TestAccZipEphemeralArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirInRegularDirAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-dir")
	if err != nil {
		t.Fatal(err)
	}

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
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

// TestAccZipEphemeralArchiveFile_Multiple_Relative_ExcludeSymlinkDirectories verifies that
// symlinked directories are excluded.
func TestAccZipEphemeralArchiveFile_Multiple_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

// TestAccZipEphemeralArchiveFile_Multiple_Relative_ExcludeSymlinkDirectories verifies that
// symlinked directories are excluded.
func TestAccZipEphemeralArchiveFile_Multiple_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_VERSION", "1.10.0-rc1")

	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	multipleDirsAndFilesAbs, err := filepath.Abs("test-fixtures")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // echo provider is protocol version 6
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []r.TestStep{
			{
				Config: addEchoConfig("test1") + fmt.Sprintf(`
			ephemeral "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(multipleDirsAndFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("echo.test1", "data.output_size", &fileSize),
					r.TestCheckResourceAttrWith("echo.test1", "data.output_path", func(value string) error {
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

func addEchoConfig(echoName string) string {
	return fmt.Sprintf(`
	provider "echo" {
		data = ephemeral.archive_file.foo
	}
	resource "echo" "%s" {}
	`, echoName)
}
