# editorconfig-checker
[![Build Status](https://travis-ci.org/editorconfig-checker/editorconfig-checker.go.svg?branch=master)](https://travis-ci.org/editorconfig-checker/editorconfig-checker.go) 
[![codecov](https://codecov.io/gh/editorconfig-checker/editorconfig-checker.go/branch/master/graph/badge.svg)](https://codecov.io/gh/editorconfig-checker/editorconfig-checker.go)
[![Go Report Card](https://goreportcard.com/badge/github.com/editorconfig-checker/editorconfig-checker.go)](https://goreportcard.com/report/github.com/editorconfig-checker/editorconfig-checker.go)

![Logo](https://raw.githubusercontent.com/editorconfig-checker/editorconfig-checker.go/master/docs/logo.png "Logo")

1. [What](#what)
1. [Installation](#installation)
1. [Usage](#usage)
1. [Excluding Files](#excluding-files)
2. [Default Excludes](#default-excludes)
2. [Manually Excluding](#manually-excluding)
3. [via ecrc](#via-ecrc)
3. [Generally](#generally)
1. [Support](#support)

## What?

This is a tool to check if your files consider your `.editorconfig`-rules. 
Most tools - like linters for example - only test one filetype and need an extra configuration. 
This tool only needs your `.editorconfig` to check all files.

If you don't know about editorconfig already you can read about it here: [editorconfig.org](https://editorconfig.org/).


## Installation

Grab a binary from the release page. 

If you have go installed you can run `go get github.com/editorconfig-checker/editorconfig-checker.go` and run `make build` inside the project folder. 
This will place a binary called `ec` into the `bin` directory.


## Usage

```
USAGE:
  -e string
        a regex which files should be excluded from checking - needs to be a valid regular expression
  -exclude string
        a regex which files should be excluded from checking - needs to be a valid regular expression
  -h    print the help
  -help
        print the help
  -v    print debugging information
  -verbose
        print debugging information
  -version
        print the version number
```

If you run this tool from a repository root it will check all files which are added to the git repository and are text files. If the tool isn't able to determine a file type it will be added to be checked too.

If you run this tool from a normal directory it will check all files which are text files. If the tool isn't able to determine a file type it will be added to be checked too.


## Excluding files

### Default excludes

If you don't manually exclude files these files are currently excluded automatic: `composer.lock`, `yarn.lock`, `*.min.css`, `*.min.js` and `package-lock.json`.

### Manually excluding

#### via ecrc

You can create a file called `.ecrc` where you can put a regular expression on each line which files should be excluded. If you do this the default excludes will *not* be active anymore.
Remember to escape your regular expressions correctly. :)

An `.ecrc` can look like this:

```
\.spec\.js\.snap$
yarn\.lock$
\.spec\.tsx\.snap$
LICENSE$
slick-styles\.vanilla-css$
banner\.js$
react_crop\.vanilla-css$
vanilla
\.svg$
Resources/Public/Plugins
README\.md$
```

#### via arguments

If you want to play around how the tool would behave you can also pass the `--exclude|-e` argument to the binary. This will accept a regular expression as well. If you use this argument the default excludes as well as the excludes from the `.ecrc`-file will *not* be active anymore.

For example: `ec --exclude node_modules`

#### Generally

If nothing is set the default excludes are considered.
If there is an `.ecrc`-file that will be considered and the default excludes will be ignored.
If there are arguments passed directly to the binary it will ignore the default excludes as well as the `.ecrc`-file.


## Support
If you have any questions, suggestions, need a wrapper for a programming language or just want to chat join #editorconfig-checker on 
freenode(IRC).
If you don't have an IRC-client set up you can use the 
[freenode webchat](https://webchat.freenode.net/?channels=editorconfig-checker).
