data "archive_file" "example" {
  type        = "zip"
  output_path = "main.zip"
  source_dir  = "${path.module}/src"
}
