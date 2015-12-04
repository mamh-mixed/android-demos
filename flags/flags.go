package flags

import (
	"flag"
	"fmt"
	"os"
)

func init() {
	var version bool

	flag.BoolVar(&version, "version", false, "print application version and build info")

	flag.Parse()

	if version {
		fmt.Printf(`quickpay v1.3.9 linux amd64
Build info:
  go version go1.5.1 darwin/amd64
  current user 芝锐 陈(501)
  build time 2015-11-19 11:43:16.399053133 +0800 CST
  system info Darwin rui-n1.local 15.0.0 Darwin Kernel Version 15.0.0: Wed Aug 26 16:57:32 PDT 2015; root:xnu-3247.1.106~1/RELEASE_X86_64 x86_64
Source version:
  git version 2.4.9 (Apple Git-60)
  current branch develop
  last commit 663a196da5f87a888b26385c8cc638fddbc8fcd9
`)

		os.Exit(0)
	}
}