package endpoints

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"poc-fiber/exceptions"
	"poc-fiber/services"
	"runtime"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"github.com/gofiber/template/html/v2"
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

		var claims struct {
			Email            string `json:"email"`
			Verified         bool   `json:"email_verified"`
			Name             string `json:"name"`
			GivenName        string `json:"given_name"`
			PreferedUserName string `json:"preferred_username"`
		}
		idToken.Claims(&claims)

		httpSession.Delete("state")
		httpSession.Set("access_token", token.AccessToken)
		httpSession.Set("refresh_token", token.RefreshToken)
		httpSession.Set("token_type", token.TokenType)
		httpSession.Set("token_expiresin", token.ExpiresIn)
		httpSession.Set("token_expiry", token.Expiry)
		errSessionSave := httpSession.Save()
		if errSessionSave != nil {
			fmt.Printf("error session save [%s]", errSessionSave.Error())
			return errSessionSave
		}
		return ctx.Render("welcome", fiber.Map{
			"UserName": claims.PreferedUserName,
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
