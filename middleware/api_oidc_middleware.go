package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"poc-fiber/commons"
	"poc-fiber/exceptions"
	"poc-fiber/oauth"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func InitOidcMiddleware(oauthmgr oauth.OAuthManager, apiBaseUri string, renewRedirectUri string, store *session.Store, clientId string, clientSecret string) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		p := c.Path()
		if strings.HasPrefix(p, apiBaseUri) {
			forceRefreshToken, accessDenied, errCheck := checkHeaderAndSession(c, store, oauthmgr)
			if accessDenied || errCheck != nil {
				if errCheck != nil {
					return c.Status(fiber.StatusUnauthorized).JSON(errCheck)
				} else {
					return c.Status(fiber.StatusUnauthorized).JSON(errors.New("access denied"))
				}
			} else if forceRefreshToken {
				httpSession, _ := store.Get(c)
				defer httpSession.Save()
				tkn := httpSession.Get(commons.SESSION_ATTR_TOKEN)
				if tkn != nil {
					tokenData, errFetch := fetchNewToken(oauthmgr.Provider, tkn.(oauth2.Token).RefreshToken, renewRedirectUri, clientId, clientSecret)
					if errFetch != nil {
						return c.Status(fiber.StatusUnauthorized).JSON(exceptions.ConvertToInternalError(errFetch))
					}
					_, errStore := oauth.VerifyAndStoreToken(c, tokenData, httpSession, oauthmgr.Verifier)
					if errStore != nil {
						return c.Status(fiber.StatusUnauthorized).JSON(exceptions.ConvertToFunctionalError(errStore, fiber.StatusUnauthorized))
					}
				} else {
					return c.Status(fiber.StatusUnauthorized).JSON(exceptions.ConvertToFunctionalError(errors.New("token not found"), fiber.StatusUnauthorized))
				}
			} else {
				return c.Next()
			}
		}
		return c.Next()
	}
}

func checkHeaderAndSession(c *fiber.Ctx, store *session.Store, oauthmgr oauth.OAuthManager) (bool, bool, error) {
	hasAuth, _, errAuth := hasAuthorizationBearer(c, oauthmgr.Verifier)
	sid := c.Cookies(commons.HEADER_SESSION_ID)
	httpSession, errSession := store.Get(c)

	var forceRefreshToken = false
	var accessDenied = false
	if hasAuth {
		if errAuth != nil {
			if isExpiredToken(errAuth) {
				forceRefreshToken = true
			} else {
				accessDenied = true
			}
		}
	} else if sid != "" {

		if errSession != nil {
			accessDenied = true
		}

		// Check if token in session exists and is valid
		_, errSession := checkRefreshTokenInSession(httpSession, oauthmgr.Verifier)
		if errSession != nil {
			if isExpiredToken(errSession) {
				forceRefreshToken = true
			} else {
				accessDenied = true
			}
		}
	}
	return forceRefreshToken, accessDenied, errSession
}

func hasAuthorizationBearer(c *fiber.Ctx, verifier *oidc.IDTokenVerifier) (bool, time.Time, error) {
	auth := c.GetReqHeaders()[commons.HEADER_AUTHORIZATION]
	var nilTime time.Time
	if auth != nil {
		expiryTime, errAuth := checkAuthorizationHeader(c, verifier)
		return true, expiryTime, errAuth
	}
	return false, nilTime, nil
}

func checkRefreshTokenInSession(httpSession *session.Session, verifier *oidc.IDTokenVerifier) (oauth2.Token, error) {
	var nilToken oauth2.Token
	tkn := httpSession.Get(commons.SESSION_ATTR_TOKEN)
	var nilTkn interface{}
	if tkn == nilTkn {
		return nilToken, errors.New("no token in session")
	} else {
		var token = tkn.(oauth2.Token)
		_, errVerify := verifier.Verify(context.Background(), token.AccessToken)
		if errVerify != nil {
			return nilToken, errVerify
		}
		return token, nil
	}
}

func checkAuthorizationHeader(c *fiber.Ctx, verifier *oidc.IDTokenVerifier) (time.Time, error) {
	var nilTime time.Time
	auth := c.GetReqHeaders()[commons.HEADER_AUTHORIZATION]
	if auth != nil {
		if !strings.HasPrefix(auth[0], commons.HEADER_BEARER) {
			return nilTime, errors.New("bearer expected")
		}
		reqToken := strings.Split(auth[0], " ")[1]
		idToken, errDecode := verifier.Verify(context.Background(), reqToken)
		if errDecode != nil {
			if errors.Is(errDecode, &oidc.TokenExpiredError{}) {
				return nilTime, errors.New("expired token")
			} else {
				return nilTime, errDecode
			}
		}
		expiryTime := idToken.Expiry
		return expiryTime, nil
	}
	return nilTime, nil
}

/* Issue refresh token request */
func fetchNewToken(provider *oidc.Provider, refreshToken string, redirectUri string, clientId string, clientSecret string) (oauth2.Token, error) {
	appBase := viper.GetString("app.server.base")
	appPort := viper.GetString("app.server.port")
	var strBuffer strings.Builder
	strBuffer.Write([]byte(appBase))
	if appPort != "" {
		strBuffer.Write([]byte(":"))
		strBuffer.Write([]byte(appPort))
	}
	strBuffer.Write([]byte(redirectUri))

	client := resty.New()
	client.SetDebug(false)
	client.SetCloseConnection(true)

	fullRedirectUri := strBuffer.String()
	res, errPost := client.SetBasicAuth(clientId, clientSecret).R().SetFormData(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"redirect_uri":  fullRedirectUri,
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

func isExpiredToken(errAuth error) bool {
	fmtError := fmt.Errorf("%w", errAuth)
	return strings.Contains(fmtError.Error(), "token is expired")
}
