# Contributing

Thank you for your interest in BaDaaS! This document provides the guidelines for how to contribute to the project through issues and pull-requests. Contributions can also come in additional ways such as joining the [DitRit discord server](https://discord.gg/zkKfj9gj2C), commenting on issues or pull requests and more.

## Issues

### Issue types

There are 3 types of issues:

- Bug report: You've found a bug with the code, and want to report it, or create an issue to track the bug.
- Discussion: You have something on your mind, which requires input form others in a discussion, before it eventually manifests as a proposal.
- Feature request: Used for items that propose a new idea or functionality. This allows feedback from others before code is written.

To ask questions and troubleshoot, please join the [DitRit discord server](https://discord.gg/zkKfj9gj2C) (use the BADAAS channel).

### Before submitting

Before you submit an issue, make sure you’ve checked the following:

1. Check for existing issues
   - Before you create a new issue, please do a search in [open issues](https://github.com/ditrit/badaas/issues) to see if the issue or feature request has already been filed.
   - If you find your issue already exists, make relevant comments and add your reaction.
2. For bugs
   - It’s not an environment issue.
   - You have as much data as possible. This usually comes in the form of logs and/or stacktrace.
3. You are assigned to the issue, a branch is created from the issue and the `wip` tag is added if you are also planning to develop the solution.

## Pull Requests

All contributions come through pull requests. To submit a proposed change, follow this workflow:

1. Make sure there's an issue (bug report or feature request) opened, which sets the expectations for the contribution you are about to make
2. Assign yourself to the issue and add the `wip` tag
3. Fork the [repo](https://github.com/ditrit/badaas) and create a new [branch](#branch-naming-policy) from the issue
4. Install the necessary [development environment](developing.md#environment)
5. Create your change and the corresponding [tests](developing.md#tests)
6. Update relevant documentation for the change in `docs/`
7. If changes are necessary in [BaDaaS example](https://github.com/ditrit/badaas-example) and [badaas-orm example](https://github.com/ditrit/badaas-orm-example), follow the same workflow there
8. Open a PR (and add links to the example repos' PR if they exist)
9. Wait for the CI process to finish and make sure all checks are green
10. A maintainer of the project will be assigned

### Use work-in-progress PRs for early feedback

A good way to communicate before investing too much time is to create a "Work-in-progress" PR and share it with your reviewers. The standard way of doing this is to add a "[WIP]" prefix in your PR’s title and assign the do-not-merge label. This will let people looking at your PR know that it is not well baked yet.

### Branch naming policy

`[BRANCH_TYPE]/[BRANCH_NAME]`

- `BRANCH_TYPE` is a prefix to describe the purpose of the branch.
  Accepted prefixes are:
  - `feature`, used for feature development
  - `bugfix`, used for bug fix
  - `improvement`, used for refactor
  - `library`, used for updating library
  - `prerelease`, used for preparing the branch for the release
  - `release`, used for releasing project
  - `hotfix`, used for applying a hotfix on main
  - `poc`, used for proof of concept
- `BRANCH_NAME` is managed by this regex: `[a-z0-9._-]` (`_` is used as space character).

## Code of Conduct

This project has adopted the [Contributor Covenant Code of Conduct](https://github.com/ditrit/badaas/blob/main/CODE_OF_CONDUCT.md)
