package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"io"
	"poc-fiber/commons"
	"poc-fiber/model"
	"strings"
)

type OAuthManager struct {
	Provider             *oidc.Provider
	OAuthEndpoints       *OAuthEndpoints
	OAuthConfig          *oauth2.Config
	Verifier             *oidc.IDTokenVerifier
	OAuthCallBackFullUrl string
	OAuthCallBackUri     string
}

type OAuthEndpoints struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	UserInfoEndpoint                  string   `json:"userinfo_endpoint"`
	EndSessionEndpoint                string   `json:"end_session_endpoint"`
	IntrospectionEndpoint             string   `json:"introspection_endpoint"`
	RevocationEndpoint                string   `json:"revocation_endpoint"`
	DeviceAuthorizationEndpoint       string   `json:"device_authorization_endpoint"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	ResponseModesSupported            []string `json:"response_modes_supported"`
	JwksUri                           string   `json:"jwks_uri"`
	GrantTypesSupported               []string `json:"grant_types_supported"`
	IdTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
	SubjectTypesSupported             []string `json:"subject_types_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	AcrValuesSupported                []string `json:"acr_values_supported"`
	ScopesSupported                   []string `json:"scopes_supported"`
	RequestParameterSupported         bool     `json:"request_parameter_supported"`
	ClaimsSupported                   []string `json:"claims_supported"`
	ClaimsPerameterSupported          bool     `json:"claims_parameter_supported"`
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported"`
}

func NewOAuthManager() OAuthManager {
	oauthMgr := OAuthManager{}
	return oauthMgr
}

func (oAuthManager OAuthManager) InitOAuthManager(logger zap.Logger, clientId string, clientSecret string) (OAuthManager, error) {
	appBase := viper.GetString("app.server.base")
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
	oAuthManager.OAuthEndpoints = FetchOAuthConfiguration(oauth2Issuer, logger)

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

func FetchOAuthConfiguration(rootUrl string, logger zap.Logger) *OAuthEndpoints {
	var wellKnown = strings.TrimSuffix(rootUrl, "/") + "/.well-known/openid-configuration"
	client := resty.New()
	client.SetDebug(false)
	client.SetCloseConnection(true)
	res, errGet := client.R().SetHeader("Cache-Control", "no-cache").Get(wellKnown)
	if errGet != nil {
		logger.Error("error retrieving .wellknown openid-configuration", zap.Error(errGet))
	}
	var oauthConfig *OAuthEndpoints
	errUnmarshal := json.Unmarshal(res.Body(), &oauthConfig)
	if errUnmarshal != nil {
		logger.Error("error unmarshalling .wellknown openid-configuration", zap.Error(errUnmarshal))
	}
	return oauthConfig
}

func VerifyAndStoreToken(token oauth2.Token, httpSession *session.Session, verifier *oidc.IDTokenVerifier) (model.Claims, error) {
	var claims model.Claims
	idToken, errVerify := verifier.Verify(context.Background(), token.AccessToken)
	if errVerify != nil {
		return claims, errVerify
	}

	errClaims := idToken.Claims(&claims)
	if errClaims != nil {
		return claims, errClaims
	}
	StoreToken(httpSession, token, claims)
	return claims, nil
}

func StoreToken(httpSession *session.Session, token oauth2.Token, claims model.Claims) {
	httpSession.Set(commons.SESSION_ATTR_TOKEN, token)
	httpSession.Set(commons.SESSION_ATTR_USERNAME, claims.PreferedUserName)
}

func GenerateState(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	state := base64.StdEncoding.EncodeToString(data)
	return state, nil
}
