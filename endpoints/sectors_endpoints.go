package endpoints

import (
	"poc-fiber/converters"
	"poc-fiber/dtos"
	"poc-fiber/exceptions"
	"poc-fiber/services"
	"poc-fiber/validation"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func MakeSectorCreate(sectorsSvc services.SectorService, logger zap.Logger) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tenantUuid := ctx.Params("tenantUuid")
		orgUuid := ctx.Params("organizationUuid")

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

		cid, errCreate := sectorsSvc.CreateSector(tenantUuid, orgUuid, sectorReq, logger)
		if errCreate != nil {
			_ = ctx.SendStatus(fiber.StatusBadRequest)
			apiErr := exceptions.ConvertToFunctionalError(errCreate, fiber.StatusBadRequest)
			return ctx.JSON(apiErr)
		}
		uuidResponse := dtos.UuidResponse{
			Uuid: cid.Uuid,
		}
		ctx.SendStatus(fiber.StatusOK)
		return ctx.JSON(uuidResponse)
	}
}

func MakeSectorsFindAll(sectorsSvc services.SectorService, logger zap.Logger) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tenantUuid := ctx.Params("tenantUuid")
		orgUuid := ctx.Params("organizationUuid")
		sectorsList, errFindAll := sectorsSvc.FindSectorsByTenantAndOrganization(tenantUuid, orgUuid, logger)
		sectorLightResponse := converters.BuildSectorsLightHierarchy(sectorsList)
		var sectorLightResponseList = dtos.SectorLightResponseList{
			Sectors: sectorLightResponse,
		}

		if errFindAll != nil {
			_ = ctx.SendStatus(fiber.StatusBadRequest)
			apiErr := exceptions.ConvertToFunctionalError(errFindAll, fiber.StatusBadRequest)
			return ctx.JSON(apiErr)
		} else {
			_ = ctx.SendStatus(fiber.StatusOK)
			return ctx.JSON(sectorLightResponseList)
		}
	}
}
