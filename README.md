# x

[![Go Report Card](https://goreportcard.com/badge/github.com/petabytecl/x)](https://goreportcard.com/report/github.com/petabytecl/x)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/petabytecl/x)](https://pkg.go.dev/github.com/petabytecl/x)

## examples

```go
package main

import (
    "github.com/petabytecl/x/errorsx"
    "github.com/petabytecl/x/fmtx"
)

func main() {
    fmtx.JSONPrettyPrint(errorsx.ErrInternalServerError)
}
```
