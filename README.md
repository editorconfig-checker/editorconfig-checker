# editorconfig-checker
[![Build Status](https://travis-ci.org/editorconfig-checker/editorconfig-checker.svg?branch=master)](https://travis-ci.org/editorconfig-checker/editorconfig-checker) 
[![codecov](https://codecov.io/gh/editorconfig-checker/editorconfig-checker/branch/master/graph/badge.svg)](https://codecov.io/gh/editorconfig-checker/editorconfig-checker)
[![Go Report Card](https://goreportcard.com/badge/github.com/editorconfig-checker/editorconfig-checker)](https://goreportcard.com/report/github.com/editorconfig-checker/editorconfig-checker)

![Logo](https://raw.githubusercontent.com/editorconfig-checker/editorconfig-checker/master/docs/logo.png "Logo")

1. [What](#what)  
2. [Installation](#installation)  
3. [Usage](#usage)  
4. [Excluding](#excluding)  
4.1 [Excluding Lines](#excluding-lines)  
4.2 [Excluding Files](#excluding-files)  
4.2.1 [Inline](#inline)  
4.2.2. [Default Excludes](#default-excludes)  
4.2.3. [Manually Excluding](#manually-excluding)  
4.2.4. [via ecrc](#via-ecrc)  
4.2.5. [via arguments](#via-arguments)  
4.2.6. [Generally](#generally)  
5. [Docker](#docker)
6. [Support](#support)


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

Unsupported features are:
* `charset`

## Installation

Grab a binary from the [release page](https://github.com/editorconfig-checker/editorconfig-checker/releases). 

If you have go installed you can run `go get github.com/editorconfig-checker/editorconfig-checker` and run `make build` inside the project folder. 
This will place a binary called `ec` into the `bin` directory.


## Usage

```
USAGE:
  -d    show which files would be checked
  -dry-run
        show which files would be checked
  -e string
        a regex which files should be excluded from checking - needs to be a valid regular expression
  -exclude string
        a regex which files should be excluded from checking - needs to be a valid regular expression
  -h    print the help
  -help
        print the help
  -i    ignore default excludes
  -ignore
        ignore default excludes
  -v    print debugging information
  -verbose
        print debugging information
  -version
        print the version number
```

If you run this tool from a repository root it will check all files which are added to the git repository and are text files. If the tool isn't able to determine a file type it will be added to be checked too.

If you run this tool from a normal directory it will check all files which are text files. If the tool isn't able to determine a file type it will be added to be checked too.

## Excluding

### Excluding lines

You can exclude single lines inline. To do that you need a comment on that line that says: `editorconfig-checker-disable-line`.

```js
const myTemplateString = `
  first line 
     wrongly indended line because it needs to be` // editorconfig-checker-disable-line
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

If you don't pass the `i` or `ignore` flag to the binary these files are excluded automatically:
```
"yarn\\.lock$",
"package-lock\\.json",
"composer\\.lock$",
"\\.snap$",
"\\.otf$",
"\\.woff$",
"\\.woff2$",
"\\.eot$",
"\\.ttf$",
"\\.gif$",
"\\.png$",
"\\.jpg$",
"\\.jepg$",
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
"\\.css\\.map$",
"\\.js\\.map$",
"min\\.css$",
"min\\.js$"
```

#### Manually excluding

##### via ecrc

You can create a file called `.ecrc` where you can put a regular expression on each line which files should be excluded. If you do this it will be merged with the default excludes.
Remember to escape your regular expressions correctly. :)

An `.ecrc` can look like this:

```
LICENSE$
slick-styles\.vanilla-css$
banner\.js$
react_crop\.vanilla-css$
vanilla
Resources/Public/Plugins
README\.md$
```

##### via arguments

If you want to play around how the tool would behave you can also pass the `--exclude|-e` argument to the binary. This will accept a regular expression as well. If you use this argument the default excludes as well as the excludes from the `.ecrc`-file will merged together.

For example: `ec --exclude node_modules`

##### Generally

Every exclude option is merged together.

If you want to see which files the tool would check without checking them you can pass the `--dry-run` or `-d` flag.

## Docker 

You are able to run this tool inside a Docker container.
To do this you need to have Docker installed and run this command in your repository root which you want to check:
`docker run --rm --volume=$PWD:/check mstruebing/editorconfig-checker`

Dockerhub: [mstruebing/editorconfig-checker](https://hub.docker.com/r/mstruebing/editorconfig-checker)

## Support
If you have any questions, suggestions, need a wrapper for a programming language or just want to chat join #editorconfig-checker on 
freenode(IRC).
If you don't have an IRC-client set up you can use the 
[freenode webchat](https://webchat.freenode.net/?channels=editorconfig-checker).
