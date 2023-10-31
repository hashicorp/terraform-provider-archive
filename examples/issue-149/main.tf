data "archive_file" "input_archive" {
  type        = "zip"
  source_file = "input.txt"
  output_path = "input.zip"
}

data "archive_file" "output_archive" {
  type        = "zip"
  output_path = "output.zip"

  source {
    content  = filebase64(data.archive_file.input_archive.output_path)
    filename = data.archive_file.input_archive.output_path
  }
}
