package security

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"poc-fiber/commons"
	"poc-fiber/model"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
)

func VerifyAndStoreToken(ctx *fiber.Ctx, token oauth2.Token, httpSession *session.Session, verifier *oidc.IDTokenVerifier) (model.Claims, error) {
	var claims model.Claims
	idToken, errVerify := verifier.Verify(context.Background(), token.AccessToken)
	if errVerify != nil {
		return claims, errVerify
	}

	errClaims := idToken.Claims(&claims)
	if errClaims != nil {
		return claims, errClaims
	}
	StoreToken(httpSession, ctx, token, claims)
	return claims, nil
}

func StoreToken(httpSession *session.Session, ctx *fiber.Ctx, token oauth2.Token, claims model.Claims) {
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
