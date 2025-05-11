package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"poc-fiber/commons"
	"poc-fiber/exceptions"
	"poc-fiber/oauth"
	"poc-fiber/security"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func InitOidcMiddleware(oauthmgr oauth.OAuthManager, apiBaseUri string, renewRedirectUri string, store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		p := c.Path()
		if strings.HasPrefix(p, apiBaseUri) {
			auth := c.GetReqHeaders()[commons.HEADER_AUTHORIZATION]
			sid := c.Cookies(commons.HEADER_SESSION_ID)
			if auth != nil {
				errAuth := checkAuthorizationHeader(c, oauthmgr.Verifier)
				if errAuth != nil {
					return c.Status(fiber.StatusUnauthorized).JSON(exceptions.ConvertToFunctionalError(errAuth, fiber.StatusUnauthorized))
				}
			} else if sid != "" {
				httpSession, errSession := store.Get(c)
				defer httpSession.Save()
				if errSession != nil {
					c.Status(fiber.StatusUnauthorized).JSON(exceptions.ConvertToFunctionalError(errors.New("no session found"), fiber.StatusUnauthorized))
				}
				refreshToken, errSession := checkRefreshTokenInSession(httpSession)
				if errSession != nil {
					return c.Status(fiber.StatusUnauthorized).JSON(exceptions.ConvertToFunctionalError(errSession, fiber.StatusUnauthorized))
				}
				tokenData, errFetch := fetchNewToken(oauthmgr.Provider, refreshToken, renewRedirectUri, viper.GetString("oauth2.clientId"), "oauth2.clientSecret")
				if errFetch != nil {
					return c.Status(fiber.StatusUnauthorized).JSON(exceptions.ConvertToInternalError(errFetch))
				}
				_, errStore := security.VerifyAndStoreToken(c, tokenData, httpSession, oauthmgr.Verifier)
				if errStore != nil {
					return c.Status(fiber.StatusUnauthorized).JSON(exceptions.ConvertToFunctionalError(errStore, fiber.StatusUnauthorized))
				}
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(errors.New("neither authorization header nor session cookie found"))
			}
		}
		return c.Next()
	}
}

func checkRefreshTokenInSession(httpSession *session.Session) (string, error) {
	var nilString string
	tkn := httpSession.Get(commons.SESSION_ATTR_TOKEN)
	var nilTkn interface{}
	if tkn == nilTkn {
		return nilString, errors.New("no token in session")
	} else {
		return tkn.(oauth2.Token).RefreshToken, nil
	}
}

func checkAuthorizationHeader(c *fiber.Ctx, verifier *oidc.IDTokenVerifier) error {
	auth := c.GetReqHeaders()[commons.HEADER_AUTHORIZATION]
	if auth != nil {
		if !strings.HasPrefix(auth[0], commons.HEADER_BEARER) {
			return errors.New("bearer expected")
		}
		reqToken := strings.Split(auth[0], " ")[1]
		_, errDecode := verifier.Verify(context.Background(), reqToken)
		if errDecode != nil {
			if errors.Is(errDecode, &oidc.TokenExpiredError{}) {
				return errors.New("expired token")
			} else {
				return errDecode
			}
		}
	}
	return nil
}

/* Issue refresh token request */
func fetchNewToken(provider *oidc.Provider, refreshToken string, redirectUri string, clientId string, clientSecret string) (oauth2.Token, error) {
	client := resty.New()
	client.SetDebug(true)
	client.SetCloseConnection(true)
	res, errPost := client.SetBasicAuth(clientId, clientSecret).R().SetFormData(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"redirect_uri":  redirectUri,
		"scope":         "offline_access openid profile email",
	}).
		SetHeader("Cache-Control", "no-cache").
		Post(provider.Endpoint().TokenURL)
	var newToken oauth2.Token
	if errPost != nil {
		return newToken, errPost
	}
	errUnmarshal := json.Unmarshal(res.Body(), &newToken)
	if errUnmarshal != nil {
		return newToken, errUnmarshal
	}
	return newToken, nil
}
