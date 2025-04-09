package endpoints

import (
	"errors"
	"poc-fiber/exceptions"
	"poc-fiber/services"

	"github.com/gofiber/fiber/v2"
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

	fConfig := fiber.Config{
		AppName:           appName,
		CaseSensitive:     true,
		StrictRouting:     true,
		EnablePrintRoutes: true,
		UnescapePath:      true,
		ErrorHandler:      defErrorHandler,
	}
	return fConfig
}
