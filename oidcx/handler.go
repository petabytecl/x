package oidcx

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/coreos/go-oidc"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
)

type OAuth2Handler struct {
	c        *oauth2.Config
	p        *Provider
	verifier *oidc.IDTokenVerifier
	state    string

	IDToken     *oidc.IDToken
	IDClaims    *json.RawMessage
	UserInfo    *oidc.UserInfo
	OAuth2Token *oauth2.Token
}

func (p *Provider) NewOAuth2Handler(c *oauth2.Config) *OAuth2Handler {
	return &OAuth2Handler{
		c:        c,
		p:        p,
		verifier: p.Verifier(),
		state:    randomState(),
	}
}

func (h *OAuth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != h.state {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("invalid state"))
		return
	}

	oauth2Token, err := h.c.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	h.OAuth2Token = oauth2Token

	for _, s := range h.c.Scopes {
		if s == "openid" {
			rawIDToken, ok := oauth2Token.Extra("id_token").(string)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println(err.Error())
				return
			}

			idToken, err := h.verifier.Verify(r.Context(), rawIDToken)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println(err.Error())
				return
			}

			h.IDToken = idToken

			var idTokenClaims json.RawMessage

			if err := idToken.Claims(&idTokenClaims); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println(err.Error())
				return
			}

			h.IDClaims = &idTokenClaims

			tokenSource := h.c.TokenSource(r.Context(), oauth2Token)
			userInfo, err := h.p.UserInfo(r.Context(), tokenSource)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println(err.Error())
				return
			}

			h.UserInfo = userInfo
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

func (h *OAuth2Handler) AuthCodeURL() string {
	return h.c.AuthCodeURL(
		h.state,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	)
}

func randomState() string {
	data := make([]byte, 16)
	_, _ = io.ReadFull(rand.Reader, data)

	return base64.StdEncoding.EncodeToString(data)
}
