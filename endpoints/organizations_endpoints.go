package endpoints

import (
	"poc-fiber/dtos"
	"poc-fiber/exceptions"
	"poc-fiber/services"
	"poc-fiber/validation"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

const OTEL_TRACER_NAME = "go.opentelemetry.io/contrib/examples/otel-collector"

func MakeOrgFindAll(orgSvc services.OrganizationService, logger zap.Logger) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		c, span := otel.Tracer(OTEL_TRACER_NAME).Start(ctx.Context(), "ORG-LIST-API")
		defer span.End()
		tenantUuid := ctx.Params("tenantUuid")
		orgsList, errFindAll := orgSvc.FindAllOrganizations(tenantUuid, logger, c)
		if errFindAll != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errFindAll)
			span.RecordError(errFindAll)
			return ctx.JSON(apiErr)
		} else {
			_ = ctx.SendStatus(fiber.StatusOK)
			return ctx.JSON(orgsList)
		}
	}
}

func MakeOrgCreate(orgSvc services.OrganizationService, logger zap.Logger) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tenantUuid := ctx.Params("tenantUuid")
		var orgCreateReq = dtos.CreateOrgRequest{}
		if err := ctx.BodyParser(&orgCreateReq); err != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(err)
			return ctx.JSON(apiErr)
		}

		// Validate payload
		errValid := validate.Struct(orgCreateReq)
		if errValid != nil {
			_ = ctx.SendStatus(fiber.StatusBadRequest)
			apiError := exceptions.ConvertValidationError(validation.ConvertValidationErrors(errValid))
			return ctx.JSON(apiError)
		}

		cid, errFindAll := orgSvc.CreateOrganization(tenantUuid, orgCreateReq)
		if errFindAll != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errFindAll)
			return ctx.JSON(apiErr)
		} else {
			_ = ctx.SendStatus(fiber.StatusOK)
			var uuidResponse = dtos.UuidResponse{
				Uuid: cid.Uuid,
			}
			return ctx.JSON(uuidResponse)
		}
	}
}
