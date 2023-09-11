---
page_title: "archive_file Data Source - terraform-provider-archive"
subcategory: ""
description: |-
  Generates an archive from content, a file, or directory of files.
---


<!-- Please do not edit this file, it is generated. -->
# archive_file (Data Source)

Generates an archive from content, a file, or directory of files.

## Example Usage

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { DataArchiveFile } from "./.gen/providers/archive/data-archive-file";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    new DataArchiveFile(this, "init", {
      outputPath: "${path.module}/files/init.zip",
      sourceFile: "${path.module}/init.tpl",
      type: "zip",
    });
  }
}

```

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { Token, TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { DataArchiveFile } from "./.gen/providers/archive/data-archive-file";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    new DataArchiveFile(this, "dotfiles", {
      excludes: ["${path.module}/unwanted.zip"],
      outputPath: "${path.module}/files/dotfiles.zip",
      source: [
        {
          content: Token.asString(vimrc.rendered),
          filename: ".vimrc",
        },
        {
          content: Token.asString(sshConfig.rendered),
          filename: ".ssh/config",
        },
      ],
      type: "zip",
    });
  }
}

```

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { DataArchiveFile } from "./.gen/providers/archive/data-archive-file";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    new DataArchiveFile(this, "lambda_my_function", {
      outputFileMode: "0666",
      outputPath: "${path.module}/files/lambda-my-function.js.zip",
      sourceFile: "${path.module}/../lambda/my-function/index.js",
      type: "zip",
    });
  }
}

```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `outputPath` (String) The output of the archive file.
- `type` (String) The type of archive to generate. NOTE: `zip` is supported.

### Optional

- `excludeSymlinkDirectories` (Boolean) Boolean flag indicating whether symbolically linked directories should be excluded during the creation of the archive. Defaults to `false`.
- `excludes` (Set of String) Specify files to ignore when reading the `sourceDir`.
- `outputFileMode` (String) String that specifies the octal file mode for all archived files. For example: `"0666"`. Setting this will ensure that cross platform usage of this module will not vary the modes of archived files (and ultimately checksums) resulting in more deterministic behavior.
- `source` (Block Set) Specifies attributes of a single source file to include into the archive. One and only one of `source`, `sourceContentFilename` (with `sourceContent`), `sourceFile`, or `sourceDir` must be specified. (see [below for nested schema](#nestedblock--source))
- `sourceContent` (String) Add only this content to the archive with `sourceContentFilename` as the filename. One and only one of `source`, `sourceContentFilename` (with `sourceContent`), `sourceFile`, or `sourceDir` must be specified.
- `sourceContentFilename` (String) Set this as the filename when using `sourceContent`. One and only one of `source`, `sourceContentFilename` (with `sourceContent`), `sourceFile`, or `sourceDir` must be specified.
- `sourceDir` (String) Package entire contents of this directory into the archive. One and only one of `source`, `sourceContentFilename` (with `sourceContent`), `sourceFile`, or `sourceDir` must be specified.
- `sourceFile` (String) Package this file into the archive. One and only one of `source`, `sourceContentFilename` (with `sourceContent`), `sourceFile`, or `sourceDir` must be specified.

### Read-Only

- `id` (String) The sha1 checksum hash of the output.
- `outputBase64Sha256` (String) Base64 Encoded SHA256 checksum of output file
- `outputBase64Sha512` (String) Base64 Encoded SHA512 checksum of output file
- `outputMd5` (String) MD5 of output file
- `outputSha` (String) SHA1 checksum of output file
- `outputSha256` (String) SHA256 checksum of output file
- `outputSha512` (String) SHA512 checksum of output file
- `outputSize` (Number) The byte size of the output archive file.

<a id="nestedblock--source"></a>
### Nested Schema for `source`

Required:

- `content` (String) Add this content to the archive with `filename` as the filename.
- `filename` (String) Set this as the filename when declaring a `source`.

<!-- cache-key: cdktf-0.18.0 input-d3e4d9d868812b91ebfd286a3305959de07442df7fcbac30222b0c2a83a87bdf -->