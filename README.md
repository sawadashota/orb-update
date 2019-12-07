orb-update [![GoDoc](https://godoc.org/github.com/sawadashota/orb-update?status.svg)](https://godoc.org/github.com/sawadashota/orb-update)  [![codecov](https://codecov.io/gh/sawadashota/orb-update/branch/master/graph/badge.svg)](https://codecov.io/gh/sawadashota/orb-update) [![Go Report Card](https://goreportcard.com/badge/github.com/sawadashota/orb-update)](https://goreportcard.com/report/github.com/sawadashota/orb-update) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
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

Using Docker Image
---

```
$ docker run --rm -v $(pwd):/repo sawadashota/orb-update orb-update
```

https://hub.docker.com/r/sawadashota/orb-update