English | [简体中文](README.md)

The `go-watcher` is a  tool of the Go programing language to restart automatically.
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

# go-watcher commandline args

```shell
 -w string
        path of watch, default exec path
```


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
