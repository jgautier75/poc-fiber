package endpoints

import (
	"poc-fiber/commons"
	"poc-fiber/converters"
	"poc-fiber/dtos"
	"poc-fiber/exceptions"
	"poc-fiber/services"
	"poc-fiber/validation"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

func MakeSectorCreate(sectorsSvc services.SectorService) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tenantUuid := ctx.Params("tenantUuid")
		orgUuid := ctx.Params("organizationUuid")

		c, span := otel.Tracer(OTEL_TRACER_NAME).Start(ctx.Context(), "API-SECTOR-CREATE")
		defer span.End()

		// Deserialize request
		sectorReq := dtos.SectorCreateRequest{}
		if err := ctx.BodyParser(&sectorReq); err != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(err)
			return ctx.JSON(apiErr)
		}

		// Validate payload
		errValid := validate.Struct(sectorReq)
		if errValid != nil {
			_ = ctx.SendStatus(fiber.StatusBadRequest)
			apiError := exceptions.ConvertValidationError(validation.ConvertValidationErrors(errValid))
			return ctx.JSON(apiError)
		}

		// Create sector
		cid, errCreate := sectorsSvc.CreateSector(tenantUuid, orgUuid, sectorReq, c)
		if errCreate != nil {
			var targetHttpStatus = commons.GuessHttpStatus(errCreate)
			_ = ctx.SendStatus(targetHttpStatus)
			if commons.IsKnownFunctionalError(errCreate) {
				apiErr := exceptions.ConvertToFunctionalError(errCreate, targetHttpStatus)
				return ctx.JSON(apiErr)
			} else {
				apiErr := exceptions.ConvertToInternalError(errCreate)
				return ctx.JSON(apiErr)
			}
		}
		uuidResponse := dtos.UuidResponse{
			Uuid: cid.Uuid,
		}
		ctx.SendStatus(fiber.StatusOK)
		return ctx.JSON(uuidResponse)
	}
}

func MakeSectorsFindAll(sectorsSvc services.SectorService) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tenantUuid := ctx.Params("tenantUuid")
		orgUuid := ctx.Params("organizationUuid")

		c, span := otel.Tracer(OTEL_TRACER_NAME).Start(ctx.Context(), "API-SECTOR-LIST")
		defer span.End()

		sectorsList, errFindAll := sectorsSvc.FindSectorsByTenantAndOrganization(tenantUuid, orgUuid, c)
		sectorLightResponse := converters.BuildSectorsLightHierarchy(sectorsList)
		var sectorLightResponseList = dtos.SectorLightResponseList{
			Sectors: sectorLightResponse,
		}

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
			return ctx.JSON(sectorLightResponseList)
		}
	}
}
