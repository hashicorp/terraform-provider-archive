# Archive a single directory.

data "archive_file" "example" {
  type             = "zip"
  output_path      = "${path.module}/main.zip"
  source_dir       = "${path.module}/dir"
  output_file_mode = "0666"
  excludes         = ["exclude.txt"]
}
