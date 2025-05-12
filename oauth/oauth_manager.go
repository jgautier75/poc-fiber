package oauth

import (
	"context"
	"poc-fiber/authentik"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type OAuthManager struct {
	Provider             *oidc.Provider
	OAuthEndpoints       *authentik.OAuthEndpoints
	OAuthConfig          *oauth2.Config
	Verifier             *oidc.IDTokenVerifier
	OAuthCallBackFullUrl string
	OAuthCallBackUri     string
}

func NewOAuthManager() OAuthManager {
	oauthMgr := OAuthManager{}
	return oauthMgr
}

func (oAuthManager OAuthManager) InitOAuthManager(ctx context.Context, logger zap.Logger) (OAuthManager, error) {
	appBase := viper.GetString("app.server.base")
	clientId := viper.GetString("oauth2.clientId")
	clientSecret := viper.GetString("oauth2.clientSecret")
	appContext := viper.GetString("app.server.context")
	appPort := viper.GetString("app.server.port")
	oauthCallBackUri := "/" + appContext + "/oauth2/callback"
	oauth2Issuer := viper.GetString("oauth2.issuer")

	var oauthCallBackFull string
	if appPort != "" {
		oauthCallBackFull = appBase + ":" + appPort + oauthCallBackUri
	} else {
		oauthCallBackFull = appBase + oauthCallBackUri
	}
	oAuthManager.OAuthCallBackFullUrl = oauthCallBackFull
	oAuthManager.OAuthCallBackUri = oauthCallBackUri

	// Setup OIDC - Fetch .well-known endpoint  asynchronously
	oidcprov, oidcError := oidc.NewProvider(context.Background(), oauth2Issuer)

	var nilManager OAuthManager
	if oidcError != nil {
		return nilManager, oidcError
	} else {
		oAuthManager.Provider = oidcprov
	}

	// Custom fetch (go-oidc does not fetch revoke token url)
	oAuthManager.OAuthEndpoints = authentik.FetchOAuthConfiguration(oauth2Issuer, logger)

	oAuthManager.OAuthConfig = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  oauthCallBackFull,
		// Discovery returns the OAuth2 endpoints.
		Endpoint: oAuthManager.Provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email", "offline_access"},
	}

	oAuthManager.Verifier = oAuthManager.Provider.Verifier(&oidc.Config{ClientID: clientId})
	return oAuthManager, nil
}
