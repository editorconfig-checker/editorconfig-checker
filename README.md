# editorconfig-checker

<a href="https://www.buymeacoffee.com/mstruebing" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>

[![Build Status](https://travis-ci.org/editorconfig-checker/editorconfig-checker.svg?branch=master)](https://travis-ci.org/editorconfig-checker/editorconfig-checker)
[![codecov](https://codecov.io/gh/editorconfig-checker/editorconfig-checker/branch/master/graph/badge.svg)](https://codecov.io/gh/editorconfig-checker/editorconfig-checker)
[![Hits-of-Code](https://hitsofcode.com/github/editorconfig-checker/editorconfig-checker)](https://hitsofcode.com/view/github/editorconfig-checker/editorconfig-checker)
[![Go Report Card](https://goreportcard.com/badge/github.com/editorconfig-checker/editorconfig-checker)](https://goreportcard.com/report/github.com/editorconfig-checker/editorconfig-checker)

![Logo](https://raw.githubusercontent.com/editorconfig-checker/editorconfig-checker/master/docs/logo.png "Logo")

1. [What](#what)
2. [Quickstart](#quickstart)
3. [Installation](#installation)
4. [Usage](#usage)
5. [Configuration](#configuration)
6. [Excluding](#excluding)
6.1 [Excluding Lines](#excluding-lines)
6.2 [Excluding Files](#excluding-files)
6.2.1 [Inline](#inline)
6.2.2 [Default Excludes](#default-excludes)
6.2.3 [Manually Excluding](#manually-excluding)
6.2.4 [via configuration](#via-configuration)
6.2.5 [via arguments](#via-arguments)
6.2.6 [Generally](#generally)
7. [Docker](#docker)
8. [Continuous Integration](#continous-integration)
8. [Support](#support)


## What?

![Example Screenshot](https://raw.githubusercontent.com/editorconfig-checker/editorconfig-checker/master/docs/screenshot.png "Example Screenshot")

This is a tool to check if your files consider your `.editorconfig`-rules.
Most tools - like linters for example - only test one filetype and need an extra configuration.
This tool only needs your `.editorconfig` to check all files.

If you don't know about editorconfig already you can read about it here: [editorconfig.org](https://editorconfig.org/).

Currently implemented editorconfig features are:
* `end_of_line`
* `insert_final_newline`
* `trim_trailing_whitespace`
* `indent_style`
* `indent_size`
* `max_line_length`

Unsupported features are:
* `charset`

## Quickstart

```bash
VERSION="2.4.0"
OS="linux"
ARCH="amd64"
curl -O -L -C - https://github.com/editorconfig-checker/editorconfig-checker/releases/download/$VERSION/ec-$OS-$ARCH.tar.gz && \
tar xzf ec-$OS-$ARCH.tar.gz && \
./bin/ec-$OS-$ARCH
```

## Installation

Grab a binary from the [release page](https://github.com/editorconfig-checker/editorconfig-checker/releases).

If you have go installed you can run `go get github.com/editorconfig-checker/editorconfig-checker` and run `make build` inside the project folder.
This will place a binary called `ec` into the `bin` directory.

If you are using Arch Linux, you can use [pacman](https://wiki.archlinux.org/title/Pacman) to install from [community repository](https://archlinux.org/packages/community/x86_64/editorconfig-checker/):

```
pacman -S editorconfig-checker
```

Also, development (VCS) package is available in the [AUR](https://aur.archlinux.org/packages/editorconfig-checker-git):

```
<favourite-aur-helper> <install-command> editorconfig-checker-git

# i.e.
paru -S editorconfig-checker-git
```

If go 1.16 or greater is installed, you can also install it globally via `go install`:

```bash
go install github.com/editorconfig-checker/editorconfig-checker/cmd/editorconfig-checker@latest
```

## Usage

```
USAGE:
  -config string
        config
  -debug
        print debugging information
  -disable-end-of-line
        disables the trailing whitespace check
  -disable-indent-size
        disables only the indent-size check
  -disable-indentation
        disables the indentation check
  -disable-insert-final-newline
        disables the final newline check
  -disable-trim-trailing-whitespace
        disables the trailing whitespace check
  -dry-run
        show which files would be checked
  -exclude string
        a regex which files should be excluded from checking - needs to be a valid regular expression
  -h    print the help
  -help
        print the help
  -ignore-defaults
        ignore default excludes
  -init
        creates an initial configuration
  -no-color
        dont print colors
  -v    print debugging information
  -verbose
        print debugging information
  -version
        print the version number
```

If you run this tool from a repository root it will check all files which are added to the git repository and are text files. If the tool isn't able to determine a file type it will be added to be checked too.

If you run this tool from a normal directory it will check all files which are text files. If the tool isn't able to determine a file type it will be added to be checked too.

## Configuration

The configuration is done via arguments or an `.ecrc` file.

A sample `.ecrc` file can look like this and will be used from your current working directory if not specified via the `--config` argument:

```json
{
  "Verbose": false,
  "Debug": false,
  "IgnoreDefaults": false,
  "SpacesAftertabs": false,
  "NoColor": false,
  "Exclude": [],
  "AllowedContentTypes": [],
  "PassedFiles": [],
  "Disable": {
    // set these options to true to disable specific checks
    "EndOfLine": false,
    "Indentation": false,
    "IndentSize": false,
    "InsertFinalNewline": false,
    "TrimTrailingWhitespace": false,
    "MaxLineLength": false
  }
}
```

You could also specify command line arguments and they will get merged with the configuration file, the command line arguments have a higher precedence than the configuration.

You can create a configuration with the `init`-flag. If you specify an `config`-path it will be created there.

By default the allowed_content_types are `text/` and `application/octet-stream`(needed as a fallback when no content type could be determined). You can add additional accepted content types with the `allowed_content_types` key. But the default ones doesn't get removed.

## Excluding

### Excluding lines

You can exclude single lines inline. To do that you need a comment on that line that says: `editorconfig-checker-disable-line`.

```js
const myTemplateString = `
  first line
     wrongly indended line because it needs to be` // editorconfig-checker-disable-line
```

### Excluding blocks

To temporarily disable all checks, add a comment containing `editorconfig-checker-disable`. Re-enable with a comment containing `editorconfig-checker-enable`

```js
// editorconfig-checker-disable
const myTemplateString = `
  first line
     wrongly indended line because it needs to be
`
// editorconfig-checker-enable
```

### Excluding files

#### Inline

If you want to exclude a file inline you need a comment on the first line of the file that contains: `editorconfig-checker-disable-file`

```hs
-- editorconfig-checker-disable-file
add :: Int -> Int -> Int
add x y =
  let result = x + y -- falsy indentation would not report
  in result -- falsy indentation would not report
```

#### Default excludes

If you don't pass the `ignore-defaults` flag to the binary these files are excluded automatically:
```
"^\\.yarn/",
"^yarn\\.lock$",
"^package-lock\\.json$",
"^composer\\.lock$",
"^Cargo\\.lock$",
"^\\.pnp\\.cjs$",
"^\\.pnp\\.js$",
"^\\.pnp\\.loader\\.mjs$",
"\\.snap$",
"\\.otf$",
"\\.woff$",
"\\.woff2$",
"\\.eot$",
"\\.ttf$",
"\\.gif$",
"\\.png$",
"\\.jpg$",
"\\.jpeg$",
"\\.webp$",
"\\.avif",
"\\.mp4$",
"\\.wmv$",
"\\.svg$",
"\\.ico$",
"\\.bak$",
"\\.bin$",
"\\.pdf$",
"\\.zip$",
"\\.gz$",
"\\.tar$",
"\\.7z$",
"\\.bz2$",
"\\.log$",
"\\.patch$",
"\\.css\\.map$",
"\\.js\\.map$",
"min\\.css$",
"min\\.js$"
```

#### Manually excluding

##### via configuration

In your `.ecrc` file you can exclude files with the `"exclude"` key which takes an array of regular expressions.
This will get merged with the default excludes (if not ignored). You should remember to escape your regular expressions correctly. ;)

An `.ecrc` which would ignore all test files and all markdown files can look like this:

```
{
  "Verbose": false,
  "IgnoreDefaults": false,
  "Exclude": ["testfiles", "\\.md$"],
  "SpacesAfterTabs": false,
  "Disable": {
    "EndOfLine": false,
    "Indentation": false,
    "IndentSize": false,
    "InsertFinalNewline": false,
    "TrimTrailingWhitespace": false,
    "MaxLineLength": false
  }
}
```

##### via arguments

If you want to play around how the tool would behave you can also pass the `--exclude` argument to the binary. This will accept a regular expression as well. If you use this argument the default excludes as well as the excludes from the `.ecrc`-file will merged together.

For example: `ec --exclude node_modules`

##### Generally

Every exclude option is merged together.

If you want to see which files the tool would check without checking them you can pass the `--dry-run` flag.

Note that while `--dry-run` outputs absolute paths, a regular expression matches on relative paths from where the `ec` command is used.

## Docker

You are able to run this tool inside a Docker container.
To do this you need to have Docker installed and run this command in your repository root which you want to check:
`docker run --rm --volume=$PWD:/check mstruebing/editorconfig-checker`

Dockerhub: [mstruebing/editorconfig-checker](https://hub.docker.com/r/mstruebing/editorconfig-checker)

## Continuous Integration

### Mega-Linter

Instead of installing and configuring `editorconfig-checker` and all other linters in your project CI workflows (GitHub Actions & others), you can use [Mega-Linter](https://nvuillam.github.io/mega-linter/) which does all that for you with a single [assisted installation](https://nvuillam.github.io/mega-linter/installation/)

Mega-Linter embeds [editorconfig-checker](https://nvuillam.github.io/mega-linter/descriptors/editorconfig_editorconfig_checker/) by default in all its [flavors](https://nvuillam.github.io/mega-linter/flavors/), meaning that it will be run at each commit or Pull Request to detect any issue related to .`editorconfig`

If you want to use only `editorconfig-checker` and not the 70+ other linters, you can use the following `.mega-linter.yml` configuration file

```yaml
ENABLE:
- EDITORCONFIG
```

## Support
If you have any questions, suggestions, need a wrapper for a programming language or just want to chat join #editorconfig-checker on
freenode(IRC).
If you don't have an IRC-client set up you can use the
[freenode webchat](https://webchat.freenode.net/?channels=editorconfig-checker).
