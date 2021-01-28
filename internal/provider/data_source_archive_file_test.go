package archive

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccArchiveFile_Basic(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string
	r.Test(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: testAccArchiveFileContentConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileExists(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),

					// We just check the hashes for syntax rather than exact
					// content since we don't want to break if the archive
					// library starts generating different bytes that are
					// functionally equivalent.
					r.TestMatchResourceAttr(
						"data.archive_file.foo", "output_base64sha256",
						regexp.MustCompile(`^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$`),
					),
					r.TestMatchResourceAttr(
						"data.archive_file.foo", "output_md5", regexp.MustCompile(`^[0-9a-f]{32}$`),
					),
					r.TestMatchResourceAttr(
						"data.archive_file.foo", "output_sha", regexp.MustCompile(`^[0-9a-f]{40}$`),
					),
				),
			},
			{
				Config: testAccArchiveFileFileConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileExists(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},
			{
				Config: testAccArchiveFileDirConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileExists(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},
			{
				Config: testAccArchiveFileDirExcludesConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileExists(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},

			{
				Config: testAccArchiveFileMultiConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileExists(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},

			{
				Config: testAccArchiveFileSensitiveConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileExists(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},
		},
	})
}

func testAccArchiveFileExists(filename string, fileSize *string) r.TestCheckFunc {
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
  type        = "zip"
  source_file = "test-fixtures/test-file.txt"
  output_path = "%s"
}
`, filepath.ToSlash(outputPath))
}

func testAccArchiveFileDirConfig(outputPath string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
  type        = "zip"
  source_dir  = "test-fixtures/test-dir"
  output_path = "%s"
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

func testAccArchiveFileMultiConfig(outputPath string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
	type = "zip"
	source {
		filename = "content.txt"
		content = "This is some content"
	}
	output_path = "%s"
}
`, filepath.ToSlash(outputPath))
}

func testAccArchiveFileSensitiveConfig(outputPath string) string {
	return fmt.Sprintf(`
data "archive_file" "foo" {
	type = "zip"
	sensitive_source {
		filename = "content.txt"
		content = "This is some content"
	}
	output_path = "%s"
}
`, filepath.ToSlash(outputPath))
}

func testTempDir(t *testing.T) string {
	tmp, err := ioutil.TempDir("", "tf")
	if err != nil {
		t.Fatal(err)
	}
	return tmp
}
