# Archive a file to be used with Lambda using consistent file mode

data "archive_file" "example" {
  type             = "zip"
  output_path      = "${path.module}/main.zip"
  output_file_mode = "0666"
  source_file      = "${path.module}/main.py"
}
