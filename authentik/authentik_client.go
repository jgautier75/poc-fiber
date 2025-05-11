package authentik

import (
	"strings"

	"encoding/json"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

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

func FetchOAuthConfiguration(rootUrl string, logger zap.Logger) *OAuthEndpoints {
	var wellKnown = strings.TrimSuffix(rootUrl, "/") + "/.well-known/openid-configuration"
	client := resty.New()
	client.SetDebug(true)
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
