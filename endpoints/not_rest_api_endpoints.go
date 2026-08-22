package endpoints

import (
	"context"
	"errors"
	"net/url"
	"poc-fiber/commons"
	"poc-fiber/exceptions"
	"poc-fiber/model"
	"poc-fiber/oauth"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-playground/validator"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

const HEADER_STATE = "state"
const PKCE_VERIFIER = "pkceAuthCodeOption"

var validate = validator.New()

func MakeIndex(oauthCfg *oauth2.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		httpSession := session.FromContext(c)

		state, errState := oauth.GenerateState(28)
		if errState != nil {
			apiError := exceptions.ConvertToInternalError(errState)
			return c.Status(fiber.StatusInternalServerError).JSON(apiError)
		}

		dState, err := url.QueryUnescape(state)
		if err != nil {
			apiError := exceptions.ConvertToInternalError(err)
			return c.Status(fiber.StatusInternalServerError).JSON(apiError)
		}

		pkceVerifier := oauth2.GenerateVerifier()
		pkceAuthCodeOption := oauth2.S256ChallengeOption(pkceVerifier)

		authURL := oauthCfg.AuthCodeURL(dState, pkceAuthCodeOption)

		httpSession.Set(HEADER_STATE, state)
		httpSession.Set(PKCE_VERIFIER, pkceVerifier)

		return c.Render("index", fiber.Map{
			"AuthUrl": authURL,
		})
	}
}

func MakeVersions(appVersion string) func(ctx fiber.Ctx) error {
	return func(ctx fiber.Ctx) error {
		v := model.VersionResponse{
			Version: appVersion,
		}
		return ctx.Status(fiber.StatusOK).JSON(v)
	}
}

func MakeOAuthCallback(oauthCfg *oauth2.Config, verifier *oidc.IDTokenVerifier) func(ctx fiber.Ctx) error {
	return func(ctx fiber.Ctx) error {
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
		httpSession := session.FromContext(ctx)
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

func DeleteSession(clientId string, clientSecret string, oauthCfg *oauth.OAuthEndpoints, logger zap.Logger) func(ctx fiber.Ctx) error {
	return func(ctx fiber.Ctx) error {
		httpSession := session.FromContext(ctx)
		tkn := httpSession.Get(commons.SESSION_ATTR_TOKEN)
		if tkn != nil {
			client := resty.New()
			client.SetDebug(viper.GetBool("app.debug"))
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
		logger.Error("error deleting token", zap.Error(errDelete))
	}
	logger.Info("response status", zap.String("response status", response.Status()))
}
