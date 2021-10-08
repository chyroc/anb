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

### Task Command Args

- `id`: unique id for task, can be used with set-output command
- `if`
  - `if` args support `exist` function
  - and `!`, `&&`, `||` operator
- `dir`: support cmd and local_cmd task

#### run task if path not exist

```yaml
server:
  user: root
  host: 1.2.3.4
tasks:
  - name: "clone app"
    if: |
      !exist("/app-path")
    local_cmd:
      - git clone https://github.com/user/repo app-path
```

#### run task if path exist && run command in dir

```yaml
server:
  user: root
  host: 1.2.3.4
tasks:
  - name: "pull app"
    if: exist("/app-path")
    dir: app-path
    local_cmd:
      - git pull
```

### Support Multi Task

- cmd
- local_cmd
- upload
- download

#### exec server command

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

#### exec local command

```yaml
server:
  user: root
  host: 1.2.3.4
tasks:
  - name: exec local command
    local_cmd: go build -o /tmp/bin-file main.go
  - name: exec server commands
    cmd:
      - ls
```

#### upload files from local to server

```yaml
server:
  user: root
  host: 1.2.3.4
tasks:
  - name: "upload file"
    upload:
      src: README.md
      dest: /tmp/README.md
  - name: "upload dir"
    upload:
      src: ./config/
      dest: /tmp/config/
```

#### download files from server to local

```yaml
server:
  user: root
  host: 1.2.3.4
tasks:
  - name: "download file"
    upload:
      src: /tmp/server-README.md
      dest: /tmp/local-README.md
  - name: "upload dir"
    upload:
      src: /tmp/server-config/
      dest: /tmp/local-config/
```