data "archive_file" "example" {
  type        = "zip"
  output_path = "main.zip"
  source {
    content  = "bar"
    filename = "main.txt"
  }
}
