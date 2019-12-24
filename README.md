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

orb-update can update orb and create pull request.

Here is minimum sample.

```yaml
version: 2.1

orbs:
  orb-update: sawadashota/orb-update@volatile

workflows:
  orb-update:
    jobs:
      - orb-update/orb-update
```

And following environment variables are required.

* `GITHUB_USERNAME`: GitHub token's user
* `GITHUB_TOKEN`: [GitHub access token](https://github.com/settings/tokens/new?scopes=repo,user:email&description=CircleCI%20for%20orb-update)

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

orbs:
  orb-update: sawadashota/orb-update@volatile

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

Configuration
---

Define configuration `.orb-update.yml` or CLI argument `--config`.  

```yaml
# target config file path
# default is `.circleci/config.yml`
target_files:
  - .circleci/config.yml

repository:
  # name of this repository
  name: sawadashota/orb-update

git:
  # author of commit
  # require when Pull Request Creation
  # if empty, fetch from GitHub
  author:
    name: sawadashota
    email: example@example.com

github:
  # Pull Request creation option
  # default is false
  pull_request: true

  # these should be configured by environment variable because of credentials
  #
  # `GITHUB_USERNAME`
  #username: sawadashota
  # `GITHUB_TOKEN`
  #token: github_token

# base branch
# default is `master`
base_branch: master

filesystem:
  # filesystem strategy supports `os` and `memory`
  # default is `os` for easy to use in local
  # but in CI, `memory` is recommended
  strategy: memory

ignore:
  - circleci/orb-tools
```

Using Docker Image
---

```
$ docker run --rm -v $(pwd):/repo sawadashota/orb-update orb-update
```

https://hub.docker.com/r/sawadashota/orb-update