package endpoints

import (
	"poc-fiber/commons"
	"poc-fiber/dtos"
	"poc-fiber/exceptions"
	"poc-fiber/services"
	"poc-fiber/validation"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

func MakeUserCreate(userService services.UserService, logger zap.Logger) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tenantUuid := ctx.Params("tenantUuid")
		orgUuid := ctx.Params("organizationUuid")

		c, span := otel.Tracer(OTEL_TRACER_NAME).Start(ctx.Context(), "API-USER-CREATE")
		defer span.End()

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

		cid, errCreate := userService.CreateUser(tenantUuid, orgUuid, userReq, c)
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

func MakUsersList(userService services.UserService, logger zap.Logger) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tenantUuid := ctx.Params("tenantUuid")
		orgUuid := ctx.Params("organizationUuid")

		c, span := otel.Tracer(OTEL_TRACER_NAME).Start(ctx.Context(), "USER-LIST-API")
		defer span.End()

		userListResponse, errList := userService.FindAllUsers(tenantUuid, orgUuid, c)
		if errList != nil {
			var targetHttpStatus = commons.GuessHttpStatus(errList)
			_ = ctx.SendStatus(targetHttpStatus)
			if commons.IsKnownFunctionalError(errList) {
				apiErr := exceptions.ConvertToFunctionalError(errList, targetHttpStatus)
				return ctx.JSON(apiErr)
			} else {
				apiErr := exceptions.ConvertToInternalError(errList)
				return ctx.JSON(apiErr)
			}
		}
		ctx.SendStatus(fiber.StatusOK)
		return ctx.JSON(userListResponse)
	}
}
