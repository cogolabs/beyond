package main

import (
	"context"
	"fmt"
	"log"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	oidcConfig   *oauth2.Config
	oidcProvider *oidc.Provider
	oidcVerifier *oidc.IDTokenVerifier
)

func init() {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://yourcompany.onelogin.com/oidc")
	if err != nil {
		log.Fatal(err)
	}

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := &oauth2.Config{
		ClientID:     "f8b8b020-4ec2-0135-6452-027de1ec0c4e43491",
		ClientSecret: "cxLF74XOeRRFDJbKuJpZAOtL4pVPK1t2XGVrDbe5Rx0Uij1LS2e9k7opZI6jQzHC",
		RedirectURL:  "https://beyond.colofoo.net/oidc",

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	oidcVerifier = provider.Verifier(&oidc.Config{
		ClientID: oauth2Config.ClientID,
	})
	oidcProvider = provider
	oidcConfig = oauth2Config
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

	var claims struct {
		Email string `json:"email"`
	}
	err = tokenID.Claims(&claims)
	if err != nil {
		return "", err
	}

	return claims.Email, nil
}
