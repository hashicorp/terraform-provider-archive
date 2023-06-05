# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Archive a single file.

data "archive_file" "init" {
  type        = "zip"
  source_file = "${path.module}/init.tpl"
  output_path = "${path.module}/files/init.zip"
}
