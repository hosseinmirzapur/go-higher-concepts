package openid

import (
	"context"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var provider *oidc.Provider

func GoogleConfig() *oauth2.Config {
	provider, err := oidc.NewProvider(context.Background(), "https://account.google.com")
	if err != nil {
		panic(err)
	}

	oauth2config := oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_CALLBACK_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"}, // a mechanism to limit application access to user data
		Endpoint:     provider.Endpoint(),
	}
	return &oauth2config
}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, GoogleConfig().AuthCodeURL("randomstate"), http.StatusFound)
}

func HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	// Verify state and errors.
	oauth2Token, err := GoogleConfig().Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		panic(err)
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		panic("token extraction failed")
	}

	// Parse and verify ID Token payload.
	verifier := provider.Verifier(&oidc.Config{ClientID: GoogleConfig().ClientID})
	idToken, err := verifier.Verify(context.TODO(), rawIDToken)
	if err != nil {
		panic("ID token payload verification failed")
	}

	// Extract custom claims
	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		panic(err)
	}
}
