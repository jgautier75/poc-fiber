package endpoints

import (
	"context"

	"errors"
	"fmt"

	"net/url"
	"poc-fiber/exceptions"
	"poc-fiber/model"
	"poc-fiber/security"
	"poc-fiber/services"
	"runtime"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"github.com/gofiber/template/html/v2"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

const HEADER_STATE = "state"

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
	state, errState := security.GenerateState(28)
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
		httpSession.Set(HEADER_STATE, state)
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
		reqState := ctx.Query(HEADER_STATE)
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
		claims, errorVerify := security.VerifyAndStoreToken(ctx, *token, store, verifier)
		if errorVerify != nil {
			panic(fmt.Errorf("token verification error : [%w]", errorVerify))
		}

		sid := httpSession.ID()
		httpSession.Delete(HEADER_STATE)
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

func MakeSectorsFindAll(sectorsSvc services.SectorService, logger zap.Logger) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tenantUuid := ctx.Params("tenantUuid")
		orgUuid := ctx.Params("organizationUuid")
		sectorsList, errFindAll := sectorsSvc.FindSectorsByTenantAndOrganization(tenantUuid, orgUuid, logger)
		if errFindAll != nil {
			_ = ctx.SendStatus(fiber.StatusBadRequest)
			apiErr := exceptions.ConvertToFunctionalError(errFindAll, fiber.StatusBadRequest)
			return ctx.JSON(apiErr)
		} else {
			_ = ctx.SendStatus(fiber.StatusOK)
			return ctx.JSON(sectorsList)
		}
	}
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
