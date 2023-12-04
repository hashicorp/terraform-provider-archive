data "archive_file" "foo" {
  type        = "zip"
  source_dir  = "${path.module}/foo"
  output_path = "${path.module}/foo/bar.zip"
}
