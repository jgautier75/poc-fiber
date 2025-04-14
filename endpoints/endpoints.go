package endpoints

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"poc-fiber/exceptions"
	"poc-fiber/model"
	"poc-fiber/services"
	"runtime"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"github.com/gofiber/template/html/v2"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func MakeOrgFindAll(orgSvc services.OrganizationService) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tenantUuid := ctx.Params("tenantUuid")
		orgsList, errFindAll := orgSvc.FindAllOrganizations(tenantUuid)
		if errFindAll != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errFindAll)
			return ctx.JSON(apiErr)
		} else {
			_ = ctx.SendStatus(fiber.StatusOK)
			return ctx.JSON(orgsList)
		}
	}
}

func MakeIndex(oauthCfg oauth2.Config, store *session.Store) func(ctx *fiber.Ctx) error {
	state, errState := generateState(28)
	if errState != nil {
		panic(fmt.Errorf("error reading config : [%w]", errState))
	}
	dState, _ := url.QueryUnescape(state)
	url := oauthCfg.AuthCodeURL(dState)
	return func(c *fiber.Ctx) error {
		httpSession, errSession := store.Get(c)
		if errSession != nil {
			panic(fmt.Errorf("error instantiating session : [%w]", errState))
		}
		httpSession.Set("state", state)
		errSessionSave := httpSession.Save()
		if errSessionSave != nil {
			panic(fmt.Errorf("error session save [%s]", errSessionSave))
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

func MakeOAuthCallback(oauthCfg oauth2.Config, store *session.Store, verifier *oidc.IDTokenVerifier) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		code := ctx.Query("code")
		reqState := ctx.Query("state")
		decState, errDecode := url.QueryUnescape(reqState)
		if errDecode != nil {
			panic(fmt.Errorf("error decoding state: [%w]", errDecode))
		}
		if decState != reqState {
			panic("state does not match")
		}
		httpSession, errSession := store.Get(ctx)
		if errSession != nil {
			panic(fmt.Errorf("error instantiating session : [%w]", errSession))
		}
		token, err := oauthCfg.Exchange(context.Background(), code)
		if err != nil {
			return ctx.Render("error", fiber.Map{
				"ErrorMsg": fmt.Errorf("failed to exchange token : [%w]", err),
			})
		}

		idToken, errVerify := verifier.Verify(context.Background(), token.AccessToken)
		if errVerify != nil {
			return ctx.Render("error", fiber.Map{
				"ErrorMsg": fmt.Errorf("token verification failed : [%w]", errVerify),
			})
		}

		var claims model.Claims
		idToken.Claims(&claims)

		sid := httpSession.ID()
		httpSession.Delete("state")
		httpSession.Set("token", token)
		httpSession.Set("userName", claims.PreferedUserName)
		errSessionSave := httpSession.Save()
		if errSessionSave != nil {
			fmt.Printf("error session save [%s]", errSessionSave.Error())
			return errSessionSave
		}
		return ctx.Render("welcome", fiber.Map{
			"UserName":     claims.PreferedUserName,
			"AccessToken":  token.AccessToken,
			"RefreshToken": token.RefreshToken,
			"SessionId":    sid,
		})
	}
}

func generateState(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	state := base64.StdEncoding.EncodeToString(data)
	return state, nil
}

func BuildFiberConfig(appName string) fiber.Config {
	var defErrorHandler = func(c *fiber.Ctx, err error) error {
		var e *fiber.Error
		code := fiber.StatusInternalServerError
		if errors.As(err, &e) {
			code = e.Code
			if code >= fiber.StatusBadRequest && code < fiber.StatusInternalServerError {
				apiError := exceptions.ConvertToFunctionalError(err, code)
				return c.Status(code).JSON(apiError)
			} else {
				apiError := exceptions.ConvertToInternalError(err)
				return c.Status(code).JSON(apiError)
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(exceptions.ConvertToInternalError(err))
	}

	// load only the contents of the subfolder www
	engine := html.New("./www", ".html")
	engine.Delims("{{", "}}") // define delimiters to use in the templates

	fConfig := fiber.Config{
		AppName:           appName,
		CaseSensitive:     true,
		StrictRouting:     true,
		EnablePrintRoutes: true,
		UnescapePath:      true,
		ErrorHandler:      defErrorHandler,
		Views:             engine,
	}
	return fConfig
}

func ConfigureRedisStorage(redisHost string, redisPort int) *redis.Storage {
	return redis.New(redis.Config{
		Host:      redisHost,
		Port:      redisPort,
		Username:  "",
		Password:  "",
		URL:       "",
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	},
	)
}

func RenewToken(oauthCfg oauth2.Config, store *session.Store, verifier *oidc.IDTokenVerifier,
	provider *oidc.Provider, logger zap.Logger, clientId string, clientSecret string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		httpSession, errSession := store.Get(ctx)
		if errSession != nil {
			panic(errors.New("session not found"))
		}
		tkn := httpSession.Get("token").(oauth2.Token)
		var nilToken oauth2.Token
		if tkn == nilToken {
			panic(errors.New("token not found in session"))
		}

		client := resty.New()
		client.SetDebug(true)
		client.SetCloseConnection(true)
		res, errPost := client.R().
			SetFormData(map[string]string{
				"grant_type":    "refresh_token",
				"refresh_token": tkn.RefreshToken,
				"client_id":     clientId,
				"client_secret": clientSecret,
			}).
			SetHeader("Cache-Control", "no-cache").
			Post(provider.Endpoint().TokenURL)
		if errPost != nil {
			logger.Error("error refresh token", zap.Error(errPost))
		}
		bodyString := string(res.Body())
		logger.Info(bodyString)

		var newToken oauth2.Token
		json.Unmarshal(ctx.Body(), &newToken)

		idToken, errVerify := verifier.Verify(context.Background(), newToken.AccessToken)
		if errVerify != nil {
			logger.Error("token verification failed", zap.Error(errVerify))
		}
		var claims model.Claims
		idToken.Claims(&claims)

		sid := httpSession.ID()
		httpSession.Delete("state")
		httpSession.Set("token", newToken)
		httpSession.Set("userName", claims.PreferedUserName)
		errSessionSave := httpSession.Save()
		if errSessionSave != nil {
			fmt.Printf("error session save [%s]", errSessionSave.Error())
			return errSessionSave
		}
		return ctx.Render("welcome", fiber.Map{
			"UserName":     claims.PreferedUserName,
			"AccessToken":  newToken.AccessToken,
			"RefreshToken": newToken.RefreshToken,
			"SessionId":    sid,
		})
	}
}
