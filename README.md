# x

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
