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
*  Out of the box, it‘s replace the `go run` command well.
*  To restart quick, using atomic operation make it to avoid to restart frequently.
*  Stability, using multiple semaphore make it to stop safely.

# go-watcher commandline args

```shell
 -w string
        path of watch, default exec path
```


# Tutorial

Use the `go run` command to start.

```html
go-watcher {main.go}
```

For example

```html
go-watcher {main.go} -f etc/config.yaml
```

Please give a star if it's alright, thanks!
