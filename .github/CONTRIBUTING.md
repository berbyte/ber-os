# Thank you for visiting and welcome to BER's HOWTO!
## For new contributors
You can check our [main README.md](README.md) for an overview and general technical information about BER.

## Guideline for version control
BER code and collaboration is hosted on Github. Contributing process is simple and follows [the Github flow standard](https://docs.github.com/en/get-started/using-github/github-flow).

In the following we describe the standard procedure for BER development. After having read the contributor's guide you should feel comfortable with creating discussions, raising issues, and submitting changes to the repository (by opening and merging pull requests).

In case you have questions about this document please [open a discussion](./discussions/new/choose).

### Issues
If no issue exists for the feature request, bug, changeset you would submit please open a new issue and select the preferred issue template:
 - `Feature`
 - `Bug`
 - `Feedback`
 - `Blank`

#### Connecting with branches
Consider including the issue ID in the branch name you choose to work on. Github will automagically connect the two.

#### Labelling
Try to add a maximum of 5 labels.

### Pull requests
#### Local tasks
 - [ ] `git fetch -a`
 - [ ] `git rebase origin main` [on your new branch]
 - [ ] `git rebase -i HEAD~X` [where X is an integer your number of commits]
       * This command opens the interactive TODO-EDIT text file for git
       * Review how your work is version controlled before `push`
       * Minuscule, misc. commits should be `s`quashed. Commits can be `#`commented in the next step
 - [ ] `git push` [might be `git push -u origin <BRANCH_NAME>` if it is the first push of the branch]

#### Repository tasks
 - [ ] Fill out the necessary information when opening your request
 - [ ] Set to merge your branch against `main` branch
 - [ ] Assign reviewers
 - [ ] (Optional) Set labels, milestones
 - [ ] Manage comments, change requests from reviewers
 - [ ] Finally merge and delete your branch

### Versioning
Versioning of BER uses Semver method.

It is done by using git tags.

Each version must have a corresponding release. Releases are [hosted with the repository](https://github.com/berbyte/ber/releases).

#### When to apply?
Code or contribution is tagged when and only in the case of **changes landing on the `main` branch**.

#### How to understand and apply increments?
> [!NOTE]
> In case of manual tagging make sure to use `git tag -a`.

The table shows the meaning of vMAJOR.MINOR.PATCH(-COMMIT_HASH).

| \                 | MAJOR             | MINOR             | PATCH                               | (COMMIT_HASH) |
|:-----------------:|:-----------------:|:-----------------:|:-----------------------------------:|:-------------:|
| Deployment        | Production always | Production always | Staging always, production optional | Staging only  |
| Changelog         | Yes               | Yes               | Optional                            | No            |
| Product release   | Yes               | No                | No                                  | No            |
| Milestone release | No                | Yes               | No                                  | No            |
| Issue release     | No                | No                | Yes                                 | Yes              |

In plain English this means, the aim of a PR determines the version bump happening.

Examples:
 - Fixing a typo in a file through a Pr is a type of issue release (`v1.2.3-abcdef64` => `v1.2.3-ghijkl32`)
 - Renaming and refactoring, correcting something is also a type of issue release (`v1.2.3` => `v1.2.4`)
 - Shipping a new option, releasing a new abstraction (refactoring of a different kind), producing new documentation is a type of milestone release (`v1.2.3` => `v1.3.0`)
 - Creation of a new tool, package, Api, integration, Sdk etc. is a type of product release (`v1.2.3.` => `v2.0.0`)

The frequency of releases is probably in this order. The impact is in reverse order. These two facts are considered for deployment and changelog tasks.

### Discussions
Free-for-all, for now. There are self-explanatory categories any is encouraged to be used. Standards are not yet introduced, use to your best knowledge.
