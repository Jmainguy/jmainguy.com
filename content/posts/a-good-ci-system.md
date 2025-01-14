---
title: "A Good CI System"
date: 2024-09-17T11:46:51-04:00
draft: false
---

What makes a good Continuous Integration system? I believe a good CI system should

* Require minimal to no effort on the part of the Application Developers to maintain.
* handle common tasks like linting, security scans, building artifacts, pushing artifacts to repositories, and running unit tests. 
* be flexible enough to handle simple codebases which are fine with all the defaults, and complicated codebases which wish to run additional checks / tests of their choosing.
* enforce organizational compliance standards.
* support [Shift Left](https://en.wikipedia.org/wiki/Shift-left_testing) mentality of test early and often.
* provide all information produced from testing / checks / scans inside the Pull Request / Merge Request.

I believe in the theoretical Application Developer who knows nothing, other than their specific language of choice. They have no desire to learn DevOps, SRE, System Administration, or Platform Engineering. They want to make great applications quickly, without red tape, and without being blocked by other teams.

![Know Nothing](/images/903.jpg)

The system needs to support this Application Developer persona, while acknowledging they know more about their codebase than we ever will. So we trust them to let us know when the CI system needs to add new features that all teams using that language will benefit from. We also trust them when they require more from their CI than the organizational standards require.

To this point, this week I added an easily maintained reusable workflow for use by my application developer (myself), in all their projects. I use golang in almost all my codebases, so that was the language I supported in the CI system. [https://github.com/Jmainguy/golang-workflows](https://github.com/Jmainguy/golang-workflows). 

The magic is in [https://github.com/Jmainguy/golang-workflows/blob/main/.github/workflows/golang-ci.yml](https://github.com/Jmainguy/golang-workflows/blob/main/.github/workflows/golang-ci.yml) which runs 
* go test
* golangci-lint run
* staticcheck
* gofmt
* go vet
  
This gets me unit testing, linting, formatting, and static code analysis, which points out issues in my code base like using `ioutil/io` which has been out of favor for some time now.

There is also magic in

[https://github.com/Jmainguy/golang-workflows/blob/main/.github/workflows/release-please.yml](https://github.com/Jmainguy/golang-workflows/blob/main/.github/workflows/release-please.yml) which runs

* release-please
* goreleaser

This creates github releases based on [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) and compiles the code and pushes to those releases / my docker registry when a release is cut.

I will add some security scanning to this soon, which brings up the next point, how do my developers get the latest version of CI system with minimal effort?

[https://github.com/Jmainguy/golang-workflows/blob/main/.github/workflows/release-please.yml](https://github.com/Jmainguy/golang-workflows/blob/main/.github/workflows/release-please.yml) cuts release versions of the workflows, as well as a major and minor tag pointing to the latest release in that version. The application developer trusts the Devops team to a point, and points directly to the version they trust the Devops team with. In my case, that is the major version. 

```bash
name: Call Reusable Golang Release-Please Workflow

on:
  push:
    branches:
      - main

jobs:
  release-please:
    uses: Jmainguy/golang-workflows/.github/workflows/golang-release.yml@v1
    secrets:
      token: ${{ secrets.GITHUB_TOKEN }}
```

So they always get the latest major version one without any changes on their end. If a new major version is cut, [Renovate](https://github.com/renovatebot/renovate) is configured to look at all my repositories and open a Pull Request asking the developer to upgrade to the latest release at their convenience. 

Before I added the reusable-workflow approach this week, each repo had its own workflows which looked similar but were mildly different and out of date with each other. This was causing my CI runs to fail on almost all my repositories for lack of effort on my Application Developer persona to keep them up to date. After migrating all my repos are passing their testing, and cutting releases. This solves 50% of my CI/CD needs, CD is next =).

I was quite happy to see my github notifications go from scary to empty after merging all the PR's.

 ![Clean Notifications](/images/clean-notifications.png)
