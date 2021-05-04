## 2.2.0 (May 04, 2021)

ENHANCEMENTS:

* New opt-in flag to specify the `output_file_mode` to produce more deterministic behavior across operating systems. ([#90](https://github.com/terraform-providers/terraform-provider-archive/issues/90))

DEPENDENCIES:

* Update `github.com/hashicorp/terraform-plugin-sdk/v2` to `v2.6.1` ([#95](https://github.com/terraform-providers/terraform-provider-archive/issues/95))

NOTES:

Changelogs now list all dependency updates in a separate section. These are understood to have no user-facing changes except those detailed in earlier sections.

## 2.1.0 (February 19, 2021)

Binary releases of this provider now include the darwin-arm64 platform. This version contains no further changes.

## 2.0.0 (October 14, 2020)

Binary releases of this provider now include the linux-arm64 platform.

BREAKING CHANGES:

* Upgrade to version 2 of the Terraform Plugin SDK, which drops support for Terraform 0.11. This provider will continue to work as expected for users of Terraform 0.11, which will not download the new version. ([#72](https://github.com/terraform-providers/terraform-provider-archive/issues/72))

BUG FIXES:

* Fixed path bug with exclusions on Windows ([#71](https://github.com/terraform-providers/terraform-provider-archive/issues/71))

## 1.3.0 (September 30, 2019)

NOTES:

* The provider has switched to the standalone TF SDK, there should be no noticeable impact on compatibility. ([#50](https://github.com/terraform-providers/terraform-provider-archive/issues/50))

## 1.2.2 (April 30, 2019)

* This release includes another Terraform SDK upgrade intended to align with that being used for other providers as we prepare for the Core v0.12.0 release. It should have no significant changes in behavior for this provider.

## 1.2.1 (April 12, 2019)

* This release includes only a Terraform SDK upgrade intended to align with that being used for other providers as we prepare for the Core v0.12.0 release. It should have no significant changes in behavior for this provider.

## 1.2.0 (March 20, 2019)

IMPROVEMENTS:

* The provider is now compatible with Terraform v0.12, while retaining compatibility with prior versions.

## 1.1.0 (July 30, 2018)

ENHANCEMENTS:

* Add `excludes` to the `archive_file` data source to exclude files when using `source_dir` ([#18](https://github.com/terraform-providers/terraform-provider-archive/issues/18))

BUG FIXES:

* Fix zip file path names to use forward slash on Windows ([#25](https://github.com/terraform-providers/terraform-provider-archive/issues/25))
* Fix panic in `filepath.Walk` call ([#26](https://github.com/terraform-providers/terraform-provider-archive/issues/26))

## 1.0.3 (March 23, 2018)

BUG FIXES:

* Fix modified time affecting zip contents and causing spurious diffs ([#16](https://github.com/terraform-providers/terraform-provider-archive/issues/16))

## 1.0.2 (March 16, 2018)

BUG FIXES:

* Fix issue with flags not being copied on a single file and regression introduced in 1.0.1 ([#13](https://github.com/terraform-providers/terraform-provider-archive/issues/13))

## 1.0.1 (March 13, 2018)

BUG FIXES:

* Fix issue with flags not being copied in to archive ([#9](https://github.com/terraform-providers/terraform-provider-archive/issues/9))

## 1.0.0 (September 15, 2017)

* No changes from 0.1.0; just adjusting to [the new version numbering scheme](https://www.hashicorp.com/blog/hashicorp-terraform-provider-versioning/).

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
