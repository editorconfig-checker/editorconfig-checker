# Changelog

## Unreleased
### Added
* (disable specific checks)[https://github.com/editorconfig-checker/editorconfig-checker/pull/71]
### Changed
### Deprecated
### Removed
### Fixed
### Security
### Misc

## [1.3.0] - 2019-08-06
### Added
* (allow spaces after tabs)[https://github.com/editorconfig-checker/editorconfig-checker/pull/67] with flag `spaces-after-tabs`
### Misc
* Some code refactoring (together with (allow spaces after tabs)[https://github.com/editorconfig-checker/editorconfig-checker/pull/67])
* (updated editorconfig-core version to v2 which uses go modules now)[https://github.com/editorconfig-checker/editorconfig-checker/pull/68]

## [1.2.1] - 2019-07-06
### Fixed
* fix type to correctly ignore `jpeg` (`jepg` before)
* allow every content type which contains `text/`

## [1.2.0] - 2019-05-02
### Added
* (`dry-run`-flag)[https://github.com/editorconfig-checker/editorconfig-checker/pull/60]
### Misc
* Switch to `go mod`

## [1.1.3] - 2019-04-20
### Fixed
* `insert_final_newline` behaviour according to specification (https://github.com/editorconfig-checker/editorconfig-checker/pull/56)
* Check if current branch is master and up to date with remote on release

## [1.1.2] - 2019-03-16
### Changed
* use Go 1.12 in travis
### Fixed
* use `CGO_ENABLED=0` to let the binary run on alpine
* correctly use go vet in travis

## [1.1.1] - 2019-03-01
### Fixed
* Use `.exe` extension for windows binaries

## [1.1.0] - 2019-02-27
### Added
* Changelog
* disable lines inline with `editorconfig-checker-disable-line` see https://github.com/editorconfig-checker/editorconfig-checker/pull/43
* disable files inline with `editorconfig-checker-disable-file` on first line see https://github.com/editorconfig-checker/editorconfig-checker/pull/43

## [1.0.0] - 2019-02-08
* initial release
