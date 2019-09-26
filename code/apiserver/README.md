# Build Go RESTful API by gin

根据掘金小册[《基于 Go 语言构建企业级的 RESTful API 服务》](https://juejin.im/book/5b0778756fb9a07aa632301e)的实践。相应的笔记见 notes 目录。

## How to Run

clone repo，然后把 code/apiserver 整个目录拷贝到 `$GOPATH/src` 目录下，进入 `$GOAPTH/src/apiserver` 目录，执行 `go get` 和 `make` 命令即可。

```shell
$ git clone ...
$ cd $cloned_folder
$ cp -rf ./code/apiserver $GOPATH/src/
$ cd $GOPATH/src/apiserver
$ go get
$ make
$ ./apiserver
```

客户端可以用 curl 或 postman 进行测试。
