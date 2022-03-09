---
layout: "archive"
page_title: "Archive: archive_file"
description: |-
  Generates an archive from content, a file, or directory of files.
---

# archive_file

Generates an archive from content, a file, or directory of files.

## Example Usage

```hcl
# Archive a single file.

data "archive_file" "init" {
  type        = "zip"
  source_file = "${path.module}/init.tpl"
  output_path = "${path.module}/files/init.zip"
}

# Archive multiple files and exclude file.

data "archive_file" "dotfiles" {
  type        = "zip"
  output_path = "${path.module}/files/dotfiles.zip"
  excludes    = [ "${path.module}/unwanted.zip" ]

  source {
    content  = "${data.template_file.vimrc.rendered}"
    filename = ".vimrc"
  }

  source {
    content  = "${data.template_file.ssh_config.rendered}"
    filename = ".ssh/config"
  }
}

# Archive a file to be used with Lambda using consistent file mode

data "archive_file" "lambda_my_function" {
  type             = "zip"
  source_file      = "${path.module}/../lambda/my-function/index.js"
  output_file_mode = "0666"
  output_path      = "${path.module}/files/lambda-my-function.js.zip"
}

```

~> **Note regarding symbolic links**: Due to a bug, the `archive_file` data
  source does not currently create proper zip archives when the source includes
  symbolic links (also known as "symlinks"). Please see [GitHub Issue
  #6](https://github.com/terraform-providers/terraform-provider-archive/issues/6)
  for more details and workaround options. This message will be removed when the
  bug is fixed.

## Argument Reference

The following arguments are supported:

NOTE: One of `source`, `source_content_filename` (with `source_content`), `source_file`, or `source_dir` must be specified.

* `type` - (Required) The type of archive to generate.
  NOTE: `zip` is supported.

* `output_path` - (Required) The output of the archive file.

* `output_file_mode` (Optional) String that specifies the octal file mode for all archived files. For example: `"0666"`. Setting this will ensure that cross platform usage of this module will not vary the modes of archived files (and ultimately checksums) resulting in more deterministic behavior.

* `source_content` - (Optional) Add only this content to the archive with `source_content_filename` as the filename.

* `source_content_filename` - (Optional) Set this as the filename when using `source_content`.

* `source_file` - (Optional) Package this file into the archive.

* `source_dir` - (Optional) Package entire contents of this directory into the archive.

* `source` - (Optional) Specifies attributes of a single source file to include into the archive.

* `excludes` - (Optional) Specify files to ignore when reading the `source_dir`.

The `source` block supports the following:

* `content` - (Required) Add this content to the archive with `filename` as the filename.

* `filename` - (Required) Set this as the filename when declaring a `source`.

## Attributes Reference

The following attributes are exported:

* `output_size` - The size of the output archive file.

* `output_sha` - The SHA1 checksum of output archive file.

* `output_base64sha256` - The base64-encoded SHA256 checksum of output archive file.

* `output_md5` - The MD5 checksum of output archive file.
