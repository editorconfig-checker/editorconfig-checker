# Changelog

## Unreleased
### BREAKING
### Added
### Deprecated
### Removed
### Fixed
* (Only care for not empty stringish passed files)[https://github.com/editorconfig-checker/editorconfig-checker/pull/88]
    Can lead to not correctly checking and aborting the program
* (Correctly initialize logger)[https://github.com/editorconfig-checker/editorconfig-checker/pull/90]
* (Correctly parse `--no-color`-flag)[https://github.com/editorconfig-checker/editorconfig-checker/pull/90]
### Security
### Misc
* (Use go 1.13 in ci)[https://github.com/editorconfig-checker/editorconfig-checker/pull/89]

## [2.0.2] - 2019-08-19
* (Update editorconfig-core-go to v2.1.1)[https://github.com/editorconfig-checker/editorconfig-checker/pull/77]

## [2.0.1] - 2019-08-18
* (Make allowed content types behave correctly)[https://github.com/editorconfig-checker/editorconfig-checker/pull/81]

## [2.0.0] - 2019-08-18
### BREAKING
* (Introduce `.ecrc` as a config file and not only to ignore files.)[https://github.com/editorconfig-checker/editorconfig-checker/pull/76]
* (Removed shorthand flags for clarity(see usage))[https://github.com/editorconfig-checker/editorconfig-checker/pull/76]
### Added
* (disable specific checks)[https://github.com/editorconfig-checker/editorconfig-checker/pull/71]
* (A config can be generated with the `init`-flag)[https://github.com/editorconfig-checker/editorconfig-checker/pull/76]
* (Files can be passed with `ec [<file>|<directory>]`)[https://github.com/editorconfig-checker/editorconfig-checker/pull/76]
* (Added `debug` flag (more debug output will be added while debugging))[https://github.com/editorconfig-checker/editorconfig-checker/pull/76]
### Changed
* (Print most output on `stdout` instead of `stderr` now for `grepability`)[https://github.com/editorconfig-checker/editorconfig-checker/pull/76]
### Fixed
* (fixed some golint and gocyclo issues)[https://github.com/editorconfig-checker/editorconfig-checker/pull/72]
### Misc
* (Put more structure into packages and rewrite most of internals.)[https://github.com/editorconfig-checker/editorconfig-checker/pull/76]

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
