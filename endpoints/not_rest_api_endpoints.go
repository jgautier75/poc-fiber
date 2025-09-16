package endpoints

import (
	"context"
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-playground/validator"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"net/url"
	"poc-fiber/commons"
	"poc-fiber/exceptions"
	"poc-fiber/model"
	"poc-fiber/oauth"
)

const HEADER_STATE = "state"
const PKCE_VERIFIER = "pkceAuthCodeOption"

var validate = validator.New()

func MakeIndex(oauthCfg *oauth2.Config, store *session.Store) func(ctx *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		httpSession, errSession := store.Get(c)
		if errSession != nil {
			apiError := exceptions.ConvertToInternalError(errSession)
			c.SendStatus(fiber.StatusUnauthorized)
			return c.JSON(apiError)
		}

		state, errState := oauth.GenerateState(28)
		if errState != nil {
			apiError := exceptions.ConvertToInternalError(errSession)
			c.SendStatus(fiber.StatusUnauthorized)
			return c.JSON(apiError)
		}
		dState, _ := url.QueryUnescape(state)

		pkceVerifier := oauth2.GenerateVerifier()
		pkceAuthCodeOption := oauth2.S256ChallengeOption(pkceVerifier)

		url := oauthCfg.AuthCodeURL(dState, pkceAuthCodeOption)

		httpSession.Set(HEADER_STATE, state)
		httpSession.Set(PKCE_VERIFIER, pkceVerifier)
		errSave := httpSession.Save()
		if errSave != nil {
			apiError := exceptions.ConvertToInternalError(errSave)
			c.SendStatus(fiber.StatusUnauthorized)
			return c.JSON(apiError)
		}

		return c.Render("index", fiber.Map{
			"AuthUrl": url,
		})
	}
}

func MakeVersions(appVersion string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		v := model.VersionResponse{
			Version: appVersion,
		}
		return ctx.Status(fiber.StatusOK).JSON(v)
	}
}

func MakeOAuthCallback(oauthCfg *oauth2.Config, store *session.Store, verifier *oidc.IDTokenVerifier) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		code := ctx.Query("code")
		reqState := ctx.Query(HEADER_STATE)
		decState, errDecode := url.QueryUnescape(reqState)
		if errDecode != nil {
			apiError := exceptions.ConvertToInternalError(errDecode)
			ctx.SendStatus(fiber.StatusInternalServerError)
			return ctx.JSON(apiError)
		}
		if decState != reqState {
			apiError := errors.New("state does not match")
			ctx.SendStatus(fiber.StatusUnauthorized)
			return ctx.JSON(apiError)
		}
		httpSession, errSession := store.Get(ctx)
		if errSession != nil {
			apiError := exceptions.ConvertToInternalError(errSession)
			ctx.SendStatus(fiber.StatusUnauthorized)
			return ctx.JSON(apiError)
		}
		defer httpSession.Save()

		pkceVerififer := httpSession.Get(PKCE_VERIFIER)
		token, err := oauthCfg.Exchange(context.Background(), code, oauth2.VerifierOption(pkceVerififer.(string)))
		if err != nil {
			apiError := exceptions.ConvertToInternalError(err)
			ctx.SendStatus(fiber.StatusUnauthorized)
			return ctx.JSON(apiError)
		}
		claims, errorVerify := oauth.VerifyAndStoreToken(*token, httpSession, verifier)
		if errorVerify != nil {
			apiError := exceptions.ConvertToInternalError(errorVerify)
			ctx.SendStatus(fiber.StatusUnauthorized)
			return ctx.JSON(apiError)
		}

		sid := httpSession.ID()
		httpSession.Delete(HEADER_STATE)
		return ctx.Render("welcome", fiber.Map{
			"UserName":     claims.PreferedUserName,
			"AccessToken":  token.AccessToken,
			"RefreshToken": token.RefreshToken,
			"SessionId":    sid,
		})
	}
}

func DeleteSession(clientId string, clientSecret string, store *session.Store, oauthCfg *oauth.OAuthEndpoints, logger zap.Logger) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		httpSession, errSession := store.Get(ctx)
		if errSession != nil {
			apiError := exceptions.ConvertToFunctionalError(errSession, fiber.StatusBadRequest)
			return ctx.Status(fiber.StatusBadRequest).JSON(apiError)
		}
		tkn := httpSession.Get(commons.SESSION_ATTR_TOKEN)
		if tkn != nil {
			client := resty.New()
			client.SetDebug(true)
			client.SetCloseConnection(true)
			// https://datatracker.ietf.org/doc/html/rfc7009

			// Delete access token
			resAccess, errPostAccess := client.SetBasicAuth(clientId, clientSecret).R().SetFormData(map[string]string{
				"token":           tkn.(oauth2.Token).AccessToken,
				"token_type_hint": "access_token",
			}).
				SetHeader("Cache-Control", "no-cache").
				Post(oauthCfg.RevocationEndpoint)
			if errPostAccess != nil {
				apiError := exceptions.ConvertToInternalError(errPostAccess)
				return ctx.Status(fiber.StatusInternalServerError).JSON(apiError)
			}

			logDeleteToken(resAccess, errPostAccess, logger)

			// Delete refresh token
			resRefresh, errRefresh := client.SetBasicAuth(clientId, clientSecret).R().SetFormData(map[string]string{
				"token":           tkn.(oauth2.Token).RefreshToken,
				"token_type_hint": "refresh_token",
			}).
				SetHeader("Cache-Control", "no-cache").
				Post(oauthCfg.RevocationEndpoint)
			logDeleteToken(resRefresh, errRefresh, logger)
		}
		return httpSession.Destroy()
	}
}

func logDeleteToken(response *resty.Response, errDelete error, logger zap.Logger) {
	if errDelete != nil {
		logger.Error("error deleting access token", zap.Error(errDelete))
	}
	logger.Info("response status", zap.String("response status", response.Status()))
}
