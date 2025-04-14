package security

import (
	"context"
	"encoding/json"
	"poc-fiber/model"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
)

func UnmarshalAndSaveToken(ctx *fiber.Ctx, tokenData []byte, store *session.Store, verifier *oidc.IDTokenVerifier) error {
	var newToken oauth2.Token
	json.Unmarshal(tokenData, &newToken)
	idToken, errVerify := verifier.Verify(context.Background(), newToken.AccessToken)
	if errVerify != nil {
		return errVerify
	}
	var claims model.Claims
	idToken.Claims(&claims)
	errorStorage := StoreToken(store, ctx, newToken, claims)
	return errorStorage
}

func StoreToken(store *session.Store, ctx *fiber.Ctx, token oauth2.Token, claims model.Claims) error {
	httpSession, errSession := store.Get(ctx)
	if errSession != nil {
		return errSession
	}
	httpSession.Set("token", token)
	httpSession.Set("userName", claims.PreferedUserName)
	return nil
}
