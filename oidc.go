package main

import (
	"context"
	"log"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	oidcProvider *oidc.Provider
	oidcO2Config *oauth2.Config
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

	oidcProvider = provider
	oidcO2Config = oauth2Config
}
