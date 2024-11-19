# Archive a file to be used with Lambda using consistent file mode

resource "archive_file" "lambda_my_function" {
  type             = "zip"
  source_file      = "${path.module}/../lambda/my-function/index.js"
  output_file_mode = "0666"
  output_path      = "${path.module}/files/lambda-my-function.js.zip"
}
