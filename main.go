package main

import (
	"github.com/petabytecl/x/pkg/errorsx"
	"github.com/petabytecl/x/pkg/fmtx"
)

func main() {
	fmtx.JSONPrettyPrint(errorsx.ErrInternalServerError)
}
