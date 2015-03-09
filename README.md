收银宝
=========

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
            `go get code.google.com/p/go.tools/cmd/vet`

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
