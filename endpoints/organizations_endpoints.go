package endpoints

import (
	"poc-fiber/commons"
	"poc-fiber/dtos"
	"poc-fiber/exceptions"
	"poc-fiber/services"
	"poc-fiber/validation"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

const OTEL_TRACER_NAME = "otel-collector"

func MakeOrgFindAll(orgSvc services.OrganizationService) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		c, span := otel.Tracer(OTEL_TRACER_NAME).Start(ctx.Context(), "API-ORG-LIST")
		defer span.End()
		tenantUuid := ctx.Params("tenantUuid")
		orgsList, errFindAll := orgSvc.FindAllOrganizations(tenantUuid, c)
		if errFindAll != nil {
			span.RecordError(errFindAll)
			var targetHttpStatus = commons.GuessHttpStatus(errFindAll)
			_ = ctx.SendStatus(targetHttpStatus)
			if commons.IsKnownFunctionalError(errFindAll) {
				apiErr := exceptions.ConvertToFunctionalError(errFindAll, targetHttpStatus)
				return ctx.JSON(apiErr)
			} else {
				apiErr := exceptions.ConvertToInternalError(errFindAll)
				return ctx.JSON(apiErr)
			}
		} else {
			_ = ctx.SendStatus(fiber.StatusOK)
			return ctx.JSON(orgsList)
		}
	}
}

func MakeOrgCreate(orgSvc services.OrganizationService) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {

		c, span := otel.Tracer(OTEL_TRACER_NAME).Start(ctx.Context(), "API-ORG-CREATE")
		defer span.End()

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

		// Create organization
		cid, errFindAll := orgSvc.CreateOrganization(tenantUuid, orgCreateReq, c)
		if errFindAll != nil {
			var targetHttpStatus = commons.GuessHttpStatus(errFindAll)
			_ = ctx.SendStatus(targetHttpStatus)
			if commons.IsKnownFunctionalError(errFindAll) {
				apiErr := exceptions.ConvertToFunctionalError(errFindAll, targetHttpStatus)
				return ctx.JSON(apiErr)
			} else {
				apiErr := exceptions.ConvertToInternalError(errFindAll)
				return ctx.JSON(apiErr)
			}
		} else {
			_ = ctx.SendStatus(fiber.StatusOK)
			var uuidResponse = dtos.UuidResponse{
				Uuid: cid.Uuid,
			}
			return ctx.JSON(uuidResponse)
		}
	}
}
