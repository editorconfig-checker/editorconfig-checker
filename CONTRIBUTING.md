# Contributing

Thanks a lot for your interest in contributing to **editorconfig-checker**!

## Types of contributions

Contributions to **editorconfig-checker** _include_, but are _not limited to_:

- Reporting a bug.
- Suggesting a new feature.
- Correcting spelling errors, or additions to documentation files (README, CONTRIBUTING...).
- Improving structure/format/performance/refactor/tests of the code.

All contributions are welcome from anyone willing to work in good faith with other contributors and the community. No contribution is too small, and all contributions are valued.

## Open Development

All work on **editorconfig-checker** happens directly on [GitHub](https://github.com/editorconfig-checker/editorconfig-checker). Both core team members and external contributors send pull requests, which go through the same review process.

## Pull Requests

Pull Requests are the way concrete changes are made to the code and documentation.

We encourage discussing significant changes through an [issue](https://github.com/editorconfig-checker/editorconfig-checker/issues) before submitting a PR. This helps ensure your work aligns with the project's direction and it might save you time. For smaller changes or fixes, feel free to submit a PR directly.

If you're adding new features to **editorconfig-checker**, please include tests.

The Pull Request must **pass all CI checks**, including tests and build.

The commit messages should follow the [Conventional Commits](https://www.conventionalcommits.org/) format.

Maintainers are responsible for reviewing and merging PRs, and they follow the [Maintainers](MAINTAINERS.md) guidelines.

## Software development

Our [Makefile](Makefile) provide a lot of targets that should prove helpful.
- `make run` runs a typical use case of editorconfig-checker - in this case against this very repo
- `make build` is a convenient shortcut to build the binary for your current platform
- `make test` runs all the tests, govet and asks you to gofmt if you did not already

If you want to run most of the checks the Continous Integration in GitHub does, you can also enable [pre-commit](https://pre-commit.com/). We provided a config that helps you notice problems with Conventional Commit messages or failing tests before you commit.
