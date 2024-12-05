## 2.7.0 (December 05, 2024)

FEATURES:

* resource/archive_file: Remove `deprecated` status ([#218](https://github.com/hashicorp/terraform-provider-archive/issues/218))

## 2.6.0 (September 09, 2024)

FEATURES:

* data-source/archive_file: Add support for creating `tar.gz` archive files. ([#277](https://github.com/hashicorp/terraform-provider-archive/issues/277))
* resource/archive_file: Add support for creating `tar.gz` archive files. ([#277](https://github.com/hashicorp/terraform-provider-archive/issues/277))

## 2.5.0 (July 31, 2024)

ENHANCEMENTS:

* data-source/archive_file: Add glob pattern matching support to the `excludes` attribute. ([#354](https://github.com/hashicorp/terraform-provider-archive/issues/354))
* resource/archive_file: Add glob pattern matching support to the `excludes` attribute. ([#354](https://github.com/hashicorp/terraform-provider-archive/issues/354))

## 2.4.2 (January 24, 2024)

BUG FIXES:

* data-source/archive_file: Prevent error when generating archive from source containing symbolically linked directories, and `exclude_symlink_directories` is set to true ([#298](https://github.com/hashicorp/terraform-provider-archive/issues/298))
* resource/archive_file: Prevent error when generating archive from source containing symbolically linked directories, and `exclude_symlink_directories` is set to true ([#298](https://github.com/hashicorp/terraform-provider-archive/issues/298))
* resource/archive_file: Return error when generated archive would be empty ([#298](https://github.com/hashicorp/terraform-provider-archive/issues/298))
* data-source/archive_file: Return error when generated archive would be empty ([#298](https://github.com/hashicorp/terraform-provider-archive/issues/298))

## 2.4.1 (December 18, 2023)

NOTES:

* This release introduces no functional changes. It does however include dependency updates which address upstream CVEs. ([#287](https://github.com/hashicorp/terraform-provider-archive/issues/287))

## 2.4.0 (June 07, 2023)

NOTES:

* This Go module has been updated to Go 1.19 per the [Go support policy](https://golang.org/doc/devel/release.html#policy). Any consumers building on earlier Go versions may experience errors. ([#200](https://github.com/hashicorp/terraform-provider-archive/issues/200))

ENHANCEMENTS:

* data-source/archive_file: Added attribute `exclude_symlink_directories` which will exclude symbolically linked directories from the archive when set to true. Defaults to false ([#183](https://github.com/hashicorp/terraform-provider-archive/issues/183))
* resource/archive_file: Added attribute `exclude_symlink_directories` which will exclude symbolically linked directories from the archive when set to true. Defaults to false ([#183](https://github.com/hashicorp/terraform-provider-archive/issues/183))

BUG FIXES:

* data-source/archive_file: Symbolically linked directories are included in archives by default rather than generating an error ([#183](https://github.com/hashicorp/terraform-provider-archive/issues/183))
* resource/archive_file: Symbolically linked directories are included in archives by default rather than generating an error ([#183](https://github.com/hashicorp/terraform-provider-archive/issues/183))
## 2.3.0 (January 18, 2023)

NOTES:

* Provider has been re-written using the new [`terraform-plugin-framework`](https://www.terraform.io/plugin/framework) ([#170](https://github.com/hashicorp/terraform-provider-archive/pull/170)).

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
