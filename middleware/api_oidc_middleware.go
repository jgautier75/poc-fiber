package middleware

import (
	"context"
	"errors"
	"poc-fiber/exceptions"
	"poc-fiber/security"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
)

func NewApiOidcHandler(apiBaseUri string, renewRedirectUri string, provider *oidc.Provider,
	verifier *oidc.IDTokenVerifier, store *session.Store, clientId string, clientSecret string) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		p := c.Path()
		if strings.HasPrefix(p, apiBaseUri) {
			auth := c.GetReqHeaders()["Authorization"]
			sid := c.Cookies("session_id")
			if auth != nil {
				errAuth := checkAuthorization(c, verifier)
				if errAuth != nil {
					return c.Status(fiber.StatusUnauthorized).JSON(exceptions.ConvertToFunctionalError(errAuth, fiber.StatusUnauthorized))
				}
			} else if sid != "" {
				refreshToken, errSession := checkSession(c, store)
				if errSession != nil {
					return c.Status(fiber.StatusUnauthorized).JSON(exceptions.ConvertToFunctionalError(errSession, fiber.StatusUnauthorized))
				}
				tokenData, errFetch := fetchNewToken(provider, refreshToken, renewRedirectUri, clientId, clientSecret)
				if errFetch != nil {
					return c.Status(fiber.StatusUnauthorized).JSON(exceptions.ConvertToFunctionalError(errFetch, fiber.StatusUnauthorized))
				}
				errStore := security.UnmarshalAndSaveToken(c, tokenData, store, verifier)
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

func checkSession(c *fiber.Ctx, store *session.Store) (string, error) {
	sid := c.Cookies("session_id")
	var nilString string
	if sid != "" {
		httpSession, errSession := store.Get(c)
		if errSession != nil {
			return nilString, errors.New("invalid session")
		}
		tkn := httpSession.Get("token").(oauth2.Token)
		var nilToken oauth2.Token
		if tkn == nilToken {
			return nilString, errors.New("no token in session")
		} else {
			return tkn.RefreshToken, nil
		}
	} else {
		return nilString, errors.New("no session provided")
	}
}

func checkAuthorization(c *fiber.Ctx, verifier *oidc.IDTokenVerifier) error {
	auth := c.GetReqHeaders()["Authorization"]
	if auth != nil {
		if !strings.HasPrefix(auth[0], "Bearer") {
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

func fetchNewToken(provider *oidc.Provider, refreshToken string, redirectUri string, clientId string, clientSecret string) ([]byte, error) {
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
	if errPost != nil {
		return nil, errPost
	}
	return res.Body(), nil
}
