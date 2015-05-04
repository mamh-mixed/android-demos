* 下载第三方测试包

```
go get github.com/smartystreets/goconvey/convey
```

* 下载代码覆盖率测试包(golang.org/x/tools/cmd/cover)

```
# 可以翻墙的直接下载
# 不能翻墙的利用github
git clone https://github.com/golang/tools $GOPATH/src/golang.org/x/tools
cd $GOPATH
go install golang.org/x/tools/cmd/cover
```

* 进入需要执行测试的目录

```
cd testCase
goconvey
```

打开 `http://localhost:8080/` 即可看到结果

* 如果执行时间过长(默认是5s)，导致测试被panic，启动时加上-timeout参数
```
goconvey -timeout=100s
```
