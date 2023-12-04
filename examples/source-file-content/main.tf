# Archive content.

data "archive_file" "example" {
  type                    = "zip"
  output_path             = "${path.module}/main.zip"
  source_content_filename = "example.txt"
  source_content          = "example"
}
