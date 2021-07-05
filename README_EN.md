English | [简体中文](README.md)

`go-watcher` is a hot restart tool of Golang code, can replace the `go run` command to execute any `main` entry programs, including parameters

# Requirement
```shell
 Golang version >= 1.13
```

# Installation
```html
go get -u github.com/go-component/go-watcher
```

# Feature
*  Out of the box, perfect replacement for `go run` command
*  Just in time restart, the underlying atomic operation to avoid frequent restarts
*  Strong stability, the underlying to achieve the safety of multiple semaphore shutdown

# Tutorial

The use same as  `go run` command

```html
go-watcher {main.go}
```

For example, with command-line arguments

```html
go-watcher {main.go} -f etc/config.yaml
```

Please give a star if it's alright, thanks ^ ^