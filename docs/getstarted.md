
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
