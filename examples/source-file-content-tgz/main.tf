# Archive content as tgz.

data "archive_file" "example" {
  type                    = "tgz"
  output_path             = "${path.module}/main.tar.gz"
  source_content_filename = "example.txt"
  source_content          = "example"
}
