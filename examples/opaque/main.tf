data "archive_file" "example" {
  type        = "zip"
  output_path = "example.zip"
  source_file = "example.txt"
}

data "archive_file" "main" {
  type        = "opaque"
  output_path = "${path.module}/main.zip"
  source_file = "${path.module}/example.zip"
}
