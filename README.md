收银宝
=========

# 名字解释

merchant 商户
channel 渠道
CFCA 中金

# 数据

1. 商户
2. 渠道商户
3. 路由
4. 卡属性
5. 商户绑定Id
6. 渠道绑定Id，对应关系
7. 交易，渠道应答
8. 清算
9. 交易对账

# 开发环境

* Editor

    Github Atom https://atom.io/

* plugins

    + go-plus Golang 插件

        - Autocomplete using gocode (you must have the autocomplete-plus package installed for this to work)
            `go get -u -v github.com/nsf/gocode`

        - Formatting and managing imports using `goimports`, `goreturns`, or `gofmt`
            `go get golang.org/x/tools/cmd/goimports`

        - Code quality inspection using `go vet`
            `sudo GOPATH=$GOPATH go get golang.org/x/tools/cmd/vet`

        - Linting using `golint`
            `go get github.com/golang/lint/golint`

        - Formatting source using `gofmt`
        - Syntax checking using `go build` and `go test`
        - Display of test coverage using `go test -coverprofile`

        - oracle `go get golang.org/x/tools/cmd/oracle`

    + go-def 跳转

        - go get -v code.google.com/p/rog-go/exp/cmd/godef


# 常用命令

```shell

# 查看 jks 中证书列表
keytool -list -keypass cfca1234 -keystore temp/trust.jks

# 从 jks 中导出证书
keytool -exportcert -alias cfca_ev_oca -keypass cfca1234 -keystore temp/trust.jks -rfc -file cfca_ev_oca_crt.pem

```
