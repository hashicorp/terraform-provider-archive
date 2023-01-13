package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
  source_file      = "test-fixtures/test-file.txt"
  output_path      = "%s"
  output_file_mode = "0666"
}
`, filepath.ToSlash(outputPath))
}

func testAccArchiveFileDirConfig(outputPath string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
  type             = "zip"
  source_dir       = "test-fixtures/test-dir"
  output_path      = "%s"
  output_file_mode = "0666"
}
`, filepath.ToSlash(outputPath))
}

func testAccArchiveFileDirExcludesConfig(outputPath string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
  type        = "zip"
  source_dir  = "test-fixtures/test-dir"
  excludes    = ["test-fixtures/test-dir/file2.txt"]
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
