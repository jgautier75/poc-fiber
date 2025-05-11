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

type FetchOidcResult struct {
	Provider *oidc.Provider
	Error    error
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
	var oauthCallBackFull string
	if appPort != "" {
		oauthCallBackFull = appBase + ":" + appPort + oauthCallBackUri
	} else {
		oauthCallBackFull = appBase + oauthCallBackUri
	}
	oAuthManager.OAuthCallBackFullUrl = oauthCallBackFull
	oAuthManager.OAuthCallBackUri = oauthCallBackUri
	oauth2Issuer := viper.GetString("oauth2.issuer")

	// Setup OIDC - Fetch .well-known endpoint  asynchronously
	var nilProvider oidc.Provider
	var fetchResult FetchOidcResult
	oidcprov, oidcError := oidc.NewProvider(context.Background(), oauth2Issuer)
	if oidcError != nil {
		fetchResult = FetchOidcResult{
			Provider: &nilProvider,
			Error:    oidcError,
		}
	} else {
		fetchResult = FetchOidcResult{
			Provider: oidcprov,
			Error:    nil,
		}
	}

	var nilManager OAuthManager
	if fetchResult.Error != nil {
		return nilManager, fetchResult.Error
	} else {
		oAuthManager.Provider = fetchResult.Provider
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
