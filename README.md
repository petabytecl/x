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
    "github.com/petabytecl/x/errorsx"
    "github.com/petabytecl/x/fmtx"
)

func main() {
    fmtx.JSONPrettyPrint(errorsx.ErrInternalServerError)
}
```

```go
package main

import (
  "fmt"
  "log"
  "net/http"

  "github.com/julienschmidt/httprouter"
  "github.com/petabytecl/x/oidcx"
)

const (
  oauthClientID     = "google-client-id"
  oauthClientSecret = "google-client-secret"
)

func main() {
  p, err := oidcx.NewProvider(oauthClientID, oidcx.IssuerGoogle)
  if err != nil {
    fmt.Printf("%s\n", err.Error())
  }

  handler := p.NewOAuth2Handler(
    p.NewOAuth2Config(
      oauthClientSecret,
      "http://localhost:5000/auth/google/callback",
      []string{
        "openid",
        "https://www.googleapis.com/auth/userinfo.profile",
        "https://www.googleapis.com/auth/userinfo.email",
      },
    ),
  )

  authURLParams := []oauth2.AuthCodeOption{

  }
  fmt.Printf("%s\n", handler.AuthCodeURL(authURLParams))

  r := httprouter.New()
  r.Handler("GET", "/auth/google/callback", handler)

  log.Fatal(http.ListenAndServe(":5000", r))
}
```
