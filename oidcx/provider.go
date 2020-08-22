package oidcx

import (
	"context"
	"net/url"

	"github.com/coreos/go-oidc"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	IssuerGoogle = "https://accounts.google.com"
)

type Provider struct {
	oidcConfig   *oidc.Config
	oidcProvider *oidc.Provider
}

func NewProvider(id string, issuer string) (*Provider, error) {
	p, err := oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Provider{
		oidcConfig: &oidc.Config{
			ClientID: id,
		},
		oidcProvider: p,
	}, nil
}

func (p *Provider) NewOAuth2Config(secret string, redirectUrl string, scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     p.oidcConfig.ClientID,
		ClientSecret: secret,
		Endpoint:     p.oidcProvider.Endpoint(),
		RedirectURL:  redirectUrl,
		Scopes:       scopes,
	}
}

func (p *Provider) NewClientCredentialsConfig(secret string, scopes []string, params map[string]string) *clientcredentials.Config {
	urlParams := url.Values{}

	for k, v := range params {
		urlParams.Set(k, v)
	}

	return &clientcredentials.Config{
		ClientID:       p.oidcConfig.ClientID,
		ClientSecret:   secret,
		TokenURL:       p.oidcProvider.Endpoint().TokenURL,
		Scopes:         scopes,
		EndpointParams: urlParams,
		AuthStyle:      oauth2.AuthStyleAutoDetect,
	}
}

func (p *Provider) Verifier() *oidc.IDTokenVerifier {
	return p.oidcProvider.Verifier(p.oidcConfig)
}

func (p *Provider) UserInfo(ctx context.Context, source oauth2.TokenSource) (*oidc.UserInfo, error) {
	return p.oidcProvider.UserInfo(ctx, source)
}
