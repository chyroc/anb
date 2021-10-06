# anb

[![codecov](https://codecov.io/gh/chyroc/anb/branch/master/graph/badge.svg?token=Z73T6YFF80)](https://codecov.io/gh/chyroc/anb)
[![go report card](https://goreportcard.com/badge/github.com/chyroc/anb "go report card")](https://goreportcard.com/report/github.com/chyroc/anb)
[![test status](https://github.com/chyroc/anb/actions/workflows/test.yml/badge.svg)](https://github.com/chyroc/anb/actions)
[![Apache-2.0 license](https://img.shields.io/badge/License-Apache%202.0-brightgreen.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/chyroc/anb)
[![Go project version](https://badge.fury.io/go/github.com%2Fchyroc%2Fanb.svg)](https://badge.fury.io/go/github.com%2Fchyroc%2Fanb)

![](./header.png)

## Install

By Brew:

```shell
brew install chyroc/tap/anb
```

By Go:

```shell
go get github.com/chyroc/anb
```

## Usage

### exec command

```yaml
server:
  user: root
  host: 1.2.3.4
tasks:
  - cmd: ls
  - name: exec commands
    cmd:
      - ls
      - ls -alh
```

### copy files from local to server

```yaml
server:
  user: root
  host: 1.2.3.4
tasks:
  - name: "copy file"
    copy:
      src: README.md
      dest: /tmp/README.md
  - name: "copy dir"
    copy:
      src: ./config/
      dest: /tmp/config/
```