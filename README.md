orb-update [![GoDoc](https://godoc.org/github.com/sawadashota/orb-update?status.svg)](https://godoc.org/github.com/sawadashota/orb-update) [![CircleCI](https://circleci.com/gh/sawadashota/orb-update.svg?style=svg)](https://circleci.com/gh/sawadashota/orb-update) [![codecov](https://codecov.io/gh/sawadashota/orb-update/branch/master/graph/badge.svg)](https://codecov.io/gh/sawadashota/orb-update) [![Go Report Card](https://goreportcard.com/badge/github.com/sawadashota/orb-update)](https://goreportcard.com/report/github.com/sawadashota/orb-update) [![GolangCI](https://golangci.com/badges/github.com/sawadashota/orb-update.svg)](https://golangci.com/r/github.com/sawadashota/orb-update) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
===

Update CircleCI Orbs versions

Usage
---

Create CircleCI config yaml

```yaml
orbs:
  slack: circleci/slack@1.0.0
  hello-build: circleci/hello-build@0.0.13
```

Execute command

```
$ orb-update
```

Then orb's versions are updated

```yaml
orbs:
  slack: circleci/slack@3.4.1
  hello-build: circleci/hello-build@0.0.14
```

### Pull Request Creation Option

orb-update accept `--repo` (shorthand is `-r`) option and `--pull-request` (shorthand is `-p`).

```
$ orb-update -r sawadashota/orb-update --pull-request
```

This option requires following environment variables.

* `GIT_AUTHOR_NAME`: commit's author name
* `GIT_AUTHOR_EMAIL`: commit's author email
* `GITHUB_USERNAME`: GitHub token's user
* `GITHUB_TOKEN`: [GitHub access token](https://github.com/settings/tokens/new?scopes=repo&description=Octotree%20browser%20extension)
* `BASE_BRANCH`: Base branch of Pull Request (default is `master`)
* `FILESYSTEM_STRATEGY`: Filesystem to operate. `os` or `memory` (default is `os`)

Installation
---

```
$ go get -u github.com/sawadashota/orb-update
```

or 

```
$ brew tap sawadashota/homebrew-cheers
$ brew install orb-update
```

Using CircleCI Orb
---

It's easy to check and update orb version every night. Here is an example.

```yaml
version: 2.1

workflows:
orb-update:
  jobs:
    - orb-update/orb-update:
        repository: owner/repository-name
  triggers:
    - schedule:
        cron: "0 19 * * *"
        filters:
          branches:
            only:
              - master
```

https://circleci.com/orbs/registry/orb/sawadashota/orb-update

Using Docker Image
---

```
$ docker run --rm -v $(pwd):/repo sawadashota/orb-update orb-update
```

https://hub.docker.com/r/sawadashota/orb-update