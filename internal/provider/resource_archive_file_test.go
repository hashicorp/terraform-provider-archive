package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Config: testAccArchiveFileResourceDirExcludesConfig(f),
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
				ExpectError: regexp.MustCompile(`.*At least one of these attributes must be configured:\n\[source_content_filename,source_file,source_dir]`),
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
  source_file      = "test-fixtures/test-file.txt"
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
  source_dir       = "test-fixtures/test-dir"
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
