# Archive a single directory as tgz.

data "archive_file" "example" {
  type             = "tgz"
  output_path      = "${path.module}/main.tar.gz"
  source_dir       = "${path.module}/dir"
  output_file_mode = "0400"
  excludes         = ["exclude.txt"]
}
