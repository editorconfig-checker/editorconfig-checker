# Maintainers

Maintainers are responsible for the overall health of the project. They are responsible for reviewing issues, pull requests, and ensuring that the project is moving in the right direction.

## Pull Requests - Merging Requirements

To ensure everyone in the community has the opportunity to share their feedback, ideas, and suggestions, a Pull Request must go through a review process, and must meet the following requirements to be merged:

- **Approval:** At least **2 core team members** must approve the Pull Request.
- **Open Duration:** The Pull Request must **remain open for at least 72 hours** after the **last commit** to allow for proper review. If any significant changes are made during this time, the 72-hour clock resets to ensure everyone has the chance to review the latest updates.
- **CI checks:** The Pull Request must **pass all CI checks**, including tests and build.

**Note:** If a critical bug fix is required, and the Pull Request was approved by at least **3 core team members**, the Pull Request **can be merged before the 72 hours** open duration.

The delay serves to:

- Give all maintainers time to review.
- Encourage discussions, so the best solution can emerge.
- Prevent issues that could arise from changes being merged too quickly.

This collaborative process ensures that all perspectives are considered, and helps improve the project quality as a whole.

## Releasing new versions

With goreleaser the only thing required is to create or push a new tag to GitHub. A GitHub action then will build the binaries, archives and container images, and upload them to the right places. But since we now also have release-please, all you need to do is merge the Pull Request created and maintained by release-please when you are ready to release the next version. It will then have updated all the files that need version updates (if one is missed check the `release-please-config.json` if `extra-files` contains the file, and that the concerning lines are marked with either the block or inline version of the release-please markers)
