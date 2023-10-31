# Archive a single file as tgz.

data "archive_file" "example" {
  type             = "tgz"
  output_path      = "${path.module}/main.tar.gz"
  output_file_mode = "0400"
  source_file      = "${path.module}/main.txt"
}
