package beyond

import (
	"context"
	"flag"
	"fmt"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	oidcIssuer       = flag.String("oidc-issuer", "https://yourcompany.onelogin.com/oidc", "issuer URL provided by IdP")
	oidcClientID     = flag.String("client-id", "f8b8b020-4ec2-0135-6452-027de1ec0c4e43491", "OIDC client ID")
	oidcClientSecret = flag.String("client-secret", "cxLF74XOeRRFDJbKuJpZAOtL4pVPK1t2XGVrDbe5Rx0Uij1LS2e9k7opZI6jQzHC", "OIDC client secret")

	oidcConfig   oidcConfigI
	oidcVerifier oidcVerifierI
)

type oidcClaims struct {
	Email string `json:"email"`
}

type oidcConfigI interface {
	AuthCodeURL(string, ...oauth2.AuthCodeOption) string
	Exchange(context.Context, string) (*oauth2.Token, error)
}

type oidcVerifierI interface {
	Verify(context.Context, string) (*oidc.IDToken, error)
}

func oidcSetup() error {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, *oidcIssuer)
	if err != nil {
		return err
	}

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := &oauth2.Config{
		ClientID:     *oidcClientID,
		ClientSecret: *oidcClientSecret,
		RedirectURL:  "https://" + *host + "/oidc",

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	oidcVerifier = provider.Verifier(&oidc.Config{
		ClientID: oauth2Config.ClientID,
	})
	oidcConfig = oauth2Config
	return nil
}

func oidcVerify(code string) (string, error) {
	ctx := context.Background()
	token, err := oidcConfig.Exchange(ctx, code)
	if err != nil {
		return "", err
	}

	rawID, ok := token.Extra("id_token").(string)
	if !ok {
		return "", fmt.Errorf("missing ID token")
	}

	tokenID, err := oidcVerifier.Verify(ctx, rawID)
	if err != nil {
		return "", err
	}

	claims := new(oidcClaims)
	err = tokenID.Claims(claims)
	if err != nil {
		return "", err
	}

	return claims.Email, nil
}
