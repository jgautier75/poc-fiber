package security

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"poc-fiber/model"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
)

func VerifyAndStoreToken(ctx *fiber.Ctx, token oauth2.Token, store *session.Store, verifier *oidc.IDTokenVerifier) (model.Claims, error) {
	var claims model.Claims
	idToken, errVerify := verifier.Verify(context.Background(), token.AccessToken)
	if errVerify != nil {
		return claims, errVerify
	}

	idToken.Claims(&claims)
	errorStorage := StoreToken(store, ctx, token, claims)
	return claims, errorStorage
}

func StoreToken(store *session.Store, ctx *fiber.Ctx, token oauth2.Token, claims model.Claims) error {
	httpSession, errSession := store.Get(ctx)
	if errSession != nil {
		return errSession
	}
	httpSession.Set("token", token)
	httpSession.Set("userName", claims.PreferedUserName)
	httpSession.Save()
	return nil
}

func GenerateState(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	state := base64.StdEncoding.EncodeToString(data)
	return state, nil
}
