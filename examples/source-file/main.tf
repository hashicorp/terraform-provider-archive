# Archive a single file.

data "archive_file" "example" {
  type             = "zip"
  output_path      = "${path.module}/main.zip"
  output_file_mode = "0666"
  source_file      = "${path.module}/main.txt"
}
