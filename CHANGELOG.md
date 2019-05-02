# Changelog

## Unreleased
### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
### Misc

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
