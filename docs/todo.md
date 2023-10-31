## List of issues

https://github.com/hashicorp/terraform-provider-archive/issues

#### 149

https://github.com/hashicorp/terraform-provider-archive/issues/149
Issue archiving base64 encoded content w/ source block

As designed, the call to filebase64() returns encoded data which is written to the zip entry

#### 161

https://github.com/hashicorp/terraform-provider-archive/issues/161
archive_file doesn't re-create the archive upon content change

Basics seem to work. Need to investigate use of templatefile()

#### 173

https://github.com/hashicorp/terraform-provider-archive/issues/173
Generated archive contents include an extra (empty) file when output_path is configured within same directory as source_dir.

Fix by excluding output_path if inside source_dir.

#### 172

https://github.com/hashicorp/terraform-provider-archive/issues/172
Zip file created by terraform archive_file cannot be properly read by python

Fixed by including directories (along with files) in archive.
This will change the zip file output. i.e. file size and output sha
See TestResource_UpgradeFromVersion2_2_0_DirExcludesConfig

#### 221

https://github.com/hashicorp/terraform-provider-archive/issues/221
Error generated during the execution of acceptance test on archive_file resource

This was addressed.

#### 218

https://github.com/hashicorp/terraform-provider-archive/issues/218
archive_file data source gets created during "terraform plan" vs "terraform apply" and also is not deleted during destroy

This is by design.

#### 175

https://github.com/hashicorp/terraform-provider-archive/pull/175
Remove zip files that were generated as a result of a test.

This was addressed using t.TempDir()

#### https://github.com/hashicorp/terraform-provider-archive/pull/86

https://github.com/hashicorp/terraform-provider-archive/pull/86
Support glob matching for zip excludes

Added support in checkMatch() for filepath.Match()

#### 4

https://github.com/hashicorp/terraform-provider-archive/issues/4
gzip support for archive_file

https://github.com/hashicorp/terraform-provider-archive/issues/241
Support Additional Compression Types(Ex: tar.gz format)

https://github.com/hashicorp/terraform-provider-archive/pull/29
Support and array of compression formats
zip, tar, tar.gz, base64, tar.bz2, tar.xz, tar.lz4, tar.sz

Added support for tgz type.

* zip
* tgz

#### 2

https://github.com/hashicorp/terraform-provider-archive/issues/2
Feature request - add feature to add file to pre-existing archive in Archive provider

https://github.com/hashicorp/terraform/pull/9924

#### 64

https://github.com/hashicorp/terraform-provider-archive/issues/64
Documentation missing excludes

https://github.com/hashicorp/terraform-provider-archive/issues/35
Docs are missing exclude information

Addressed in https://registry.terraform.io/providers/hashicorp/archive/latest/docs

#### Professional thoughts on the future of this provider

https://github.com/hashicorp/terraform-provider-archive/pull/29#issuecomment-406760298

## opentofu

Port to opentofu here https://github.com/opentofu/terraform-provider-archive
