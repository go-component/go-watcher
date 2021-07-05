简体中文 | [English](README_EN.md)

`go-watcher` 是一款实现 Golang 源码热重启的工具，可替换 `go run` 命令执行任意 `main` 入口程序，包括参数

# 限制
```shell
 Golang 版本 >= 1.13
```

# 安装
```html
go get -u github.com/go-component/go-watcher
```

# 特点
*   开箱即用，完美替代 `go run` 命令
*   即时重启，底层基于原子操作避免频繁重启
*   稳定性强，底层实现多信号量安全关闭

# 使用

与 `go run` 命令使用一致 

```html
go-watcher {main.go}
```

例如携带命令行参数

```html
go-watcher {main.go} -f etc/config.yaml
```

路过的小伙伴们，可以的话，还烦请点个 star 哦，谢谢 ^ ^