package oidcx

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"

	"github.com/petabytecl/x/randx"
)

type OAuth2Handler struct {
	c        *oauth2.Config
	p        *Provider
	verifier *oidc.IDTokenVerifier
	state    string
	nonce    string

	IDToken        *oidc.IDToken
	IDClaims       *json.RawMessage
	UserInfo       *oidc.UserInfo
	UserInfoClaims *json.RawMessage
	OAuth2Token    *oauth2.Token
}

func (p *Provider) NewOAuth2Handler(c *oauth2.Config) *OAuth2Handler {
	return &OAuth2Handler{
		c:        c,
		p:        p,
		verifier: p.Verifier(),
		state:    randx.MustString(24, randx.AlphaNum),
		nonce:    randx.MustString(24, randx.AlphaNum),
	}
}

func (h *OAuth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query().Get("error")) > 0 {
		fmt.Printf("auth request error: %s, desc: %s, hint: %s, debug: %s\n",
			r.URL.Query().Get("error"),
			r.URL.Query().Get("error_description"),
			r.URL.Query().Get("error_hint"),
			r.URL.Query().Get("error_debug"),
		)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("auth request error"))
		return
	}

	if r.URL.Query().Get("state") != h.state {
		fmt.Printf("states do not match. expected %s, got %s\n",
			h.state,
			r.URL.Query().Get("state"),
		)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("invalid state"))
		return
	}

	oauth2Token, err := h.c.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		fmt.Println(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.OAuth2Token = oauth2Token

	for _, s := range h.c.Scopes {
		if s == "openid" {
			rawIDToken, ok := oauth2Token.Extra("id_token").(string)
			if !ok {
				fmt.Println("unable to read id_token")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			idToken, err := h.verifier.Verify(r.Context(), rawIDToken)
			if err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if idToken.Nonce != h.nonce {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("invalid nonce"))
				return
			}

			h.IDToken = idToken

			var idTokenClaims json.RawMessage
			if err := idToken.Claims(&idTokenClaims); err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			h.IDClaims = &idTokenClaims

			tokenSource := h.c.TokenSource(r.Context(), oauth2Token)
			userInfo, err := h.p.UserInfo(r.Context(), tokenSource)
			if err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			h.UserInfo = userInfo

			var userInfoClaims json.RawMessage
			if err := userInfo.Claims(&userInfoClaims); err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			h.UserInfoClaims = &userInfoClaims
		}
	}

	json, err := jsoniter.Marshal(&h)
	if err != nil {
		fmt.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(json)
}

// The response_mode parameter addresses a different concern. If the value of response_mode
// is set to query, the response parameters are sent back to the client as query parameters
// appended to the redirect_uri; and if the value is set to fragment, then the response
// parameters are appended to the redirect_uri as a URI fragment.
// Possible values are:
//  - query
//  - fragment
//
// oauth2.SetAuthURLParam("response_mode", responseMode),

// Mitigates replay attacks. The authorization server must reject any request if it finds
// two request with the same nonce value. If a nonce is present in the authorization grant
// request, then the authorization server must include the same value in the ID token. The
// client application must validate the value of the nonce once it receives the ID token from
// the authorization server.
//
// oauth2.SetAuthURLParam("nonce", string(nonce)),

// Indicates how the client application expects the authorization server to display
// the login page and the user consent page.
// Possible values are:
//  - page
//  - popup
//  - touch
//  - wap
//
// oauth2.SetAuthURLParam("display", display),

// Indicates whether to display the login or the user consent page at the authorization
// server. If the value is none, then neither the login page nor the user consent page should be
// presented to the user. In other words, it expects the user to have an authenticated session
// at the authorization server and a preconfigured user consent. If the value is login, the
// authorization server must reauthenticate the user. If the value is consent, the authorization
// server must display the user consent page to the end user. The select_account option can be
// used if the user has multiple accounts on the authorization server. The authorization server
// must then give the user an option to select from which account he or she requires attributes.
// Possible values are:
//  - none
//  - login
//  - consent
//  - select_account
//
// oauth2.SetAuthURLParam("prompt", strings.Join(prompt, " ")),

// The parameters ask the authorization server to compare the value with max_age. If it's
// less than (current time - max_age), the authorization server must reauthenticate the user.
// When the client includes the max_age parameter in the request, the authorization server
// must include the auth_time parameter in the ID token.
//
// oauth2.SetAuthURLParam("max_age", strconv.Itoa(maxAge)),

// Expresses the end user's preferred language for the user interface.
//
// oauth2.SetAuthURLParam("ui_locales", uiLocales),

// And ID token itself. This could be an ID token previously obtained by the
// client application. If the token is encrypted, it has to be decrypted first and then encrypted
// back by the public key of the authorization server and then placed into the authentication
// request. If the value of the parameter prompt is set to none, then the id_token_hint should
// be present in the request, but it isn't a requirement.
//
// oauth2.SetAuthURLParam("id_token_hint", ""),

// This is an indication of the login identifier that the end user may use at the
// authorization server. For example, if the client application already knows the e-mail
// address of phone number of the end user, this could be set as the value of the login_hint.
// This helps provide a better user experience.
//
// oauth2.SetAuthURLParam("login_hint", loginHint),

// Stands for "Authentication Context Reference Values". It includes
// a space-separated set of values that indicates the level of authentication required at the
// authorization server. The authorization server may or may not respect these values.
//
// oauth2.SetAuthURLParam("acr_values", strings.Join(acrValues, "+")),

func (h *OAuth2Handler) AuthCodeURL(opts ...oauth2.AuthCodeOption) string {
	return h.c.AuthCodeURL(h.state, opts...)
}
