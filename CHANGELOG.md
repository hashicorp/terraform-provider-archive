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
