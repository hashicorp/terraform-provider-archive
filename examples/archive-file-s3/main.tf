# Archive a single file to s3.

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.19.0"
    }
  }
}

data "archive_file" "example" {
  type        = "zip"
  output_path = "main.zip"
  source_file = "main.py"
}

resource "aws_s3_bucket" "example" {
  bucket = "bad-lambda-layer-bucket"
  tags = {
    "adsk:moniker" = "AMPS-C-UW2"
  }
}

resource "aws_s3_object" "example" {
  bucket = aws_s3_bucket.example.id
  key    = data.archive_file.example.output_path
  source = data.archive_file.example.output_path
}
