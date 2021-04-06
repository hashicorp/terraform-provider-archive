package archive

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
					testAccArchiveFileExists(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha256", "0pc5VnHEiYEymXLjbKzuGXOxiztmeQEohwrIsqKmyCc=",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_md5", "2e7e38508e1a38efde0f6f794c185d32",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha", "d3bc49809357cf0c092211d16b0fa01d2b18684a",
					),
				),
			},
			{
				Config: testAccArchiveFileDirConfig(f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileExists(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha256", "dUpFaNIDnZC0Pp/v7iPOARsGFlEoI42v94vYHB3lggw=",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_md5", "fe4b7c8000a518c5b9e8c1769f1aacc1",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha", "1236a4cf5e93ee0cf78c8406d05ec52a9ccb9540",
					),
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

func testTempDir(t *testing.T) string {
	tmp, err := ioutil.TempDir("", "tf")
	if err != nil {
		t.Fatal(err)
	}
	return tmp
}
