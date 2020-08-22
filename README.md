# x

![GitHub issues](https://img.shields.io/github/issues-raw/petabytecl/x)
[![Go Report Card](https://goreportcard.com/badge/github.com/petabytecl/x)](https://goreportcard.com/report/github.com/petabytecl/x)
![GitHub last commit](https://img.shields.io/github/last-commit/petabytecl/x)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/petabytecl/x)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/petabytecl/x)](https://pkg.go.dev/github.com/petabytecl/x)
[![GitHub](https://img.shields.io/github/license/petabytecl/x?color=%23007D9C)](https://raw.githubusercontent.com/petabytecl/x/master/LICENSE)

## examples

```go
package main

import (
    "github.com/petabytecl/x/pkg/errorsx"
    "github.com/petabytecl/x/pkg/fmtx"
)

func main() {
    fmtx.JSONPrettyPrint(errorsx.ErrInternalServerError)
}
```
