# Archive a single file.

resource "archive_file" "init" {
  type        = "zip"
  source_file = "${path.module}/init.tpl"
  output_path = "${path.module}/files/init.zip"

  lifecycle {
    replace_triggered_by = [
      # Replace `archive_file` each time this instance of
      # the `random_integer` is replaced.
      random_integer.example.result
    ]
  }
}
