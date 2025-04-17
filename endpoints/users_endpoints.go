package endpoints

import (
	"poc-fiber/dtos"
	"poc-fiber/exceptions"
	"poc-fiber/services"
	"poc-fiber/validation"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func MakeUserCreate(userService services.UserService, logger zap.Logger) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tenantUuid := ctx.Params("tenantUuid")
		orgUuid := ctx.Params("organizationUuid")

		// Deserialize request
		userReq := dtos.CreateUserRequest{}
		if err := ctx.BodyParser(&userReq); err != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(err)
			return ctx.JSON(apiErr)
		}

		// Validate payload
		errValid := validate.Struct(userReq)
		if errValid != nil {
			_ = ctx.SendStatus(fiber.StatusBadRequest)
			apiError := exceptions.ConvertValidationError(validation.ConvertValidationErrors(errValid))
			return ctx.JSON(apiError)
		}

		cid, errCreate := userService.CreateUser(tenantUuid, orgUuid, userReq, logger)
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

func MakUsersList(userService services.UserService, logger zap.Logger) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tenantUuid := ctx.Params("tenantUuid")
		orgUuid := ctx.Params("organizationUuid")
		userListResponse, errList := userService.FindAllUsers(tenantUuid, orgUuid, logger)
		if errList != nil {
			_ = ctx.SendStatus(fiber.StatusBadRequest)
			apiErr := exceptions.ConvertToFunctionalError(errList, fiber.StatusBadRequest)
			return ctx.JSON(apiErr)
		}
		ctx.SendStatus(fiber.StatusOK)
		return ctx.JSON(userListResponse)
	}
}
