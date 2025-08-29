package endpoints

import (
	"poc-fiber/antlr/parser"
	"poc-fiber/commons"
	"poc-fiber/dtos"
	"poc-fiber/exceptions"
	"poc-fiber/model"
	"poc-fiber/services"
	"poc-fiber/validation"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

func MakeUserCreate(userService services.UserService) func(ctx *fiber.Ctx) error {
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

func MakeUserDelete(userService services.UserService) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tenantUuid := ctx.Params("tenantUuid")
		orgUuid := ctx.Params("organizationUuid")
		userUuid := ctx.Params("userUuid")

		c, span := otel.Tracer(OTEL_TRACER_NAME).Start(ctx.Context(), "API-USER-DELETE")
		defer span.End()

		deleted, errDelete := userService.DeleteUser(tenantUuid, orgUuid, userUuid, c)
		if deleted == false && errDelete != nil {
			apiErr := exceptions.ConvertToFunctionalError(errDelete, fiber.StatusBadRequest)
			return ctx.JSON(apiErr)
		}
		ctx.SendStatus(fiber.StatusNoContent)
		return nil
	}
}

func MakeUsersList(userService services.UserService) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tenantUuid := ctx.Params("tenantUuid")
		orgUuid := ctx.Params("organizationUuid")

		c, span := otel.Tracer(OTEL_TRACER_NAME).Start(ctx.Context(), "API-USER-LIST")
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

func MakeUsersDelete(userService services.UserService) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tenantUuid := ctx.Params("tenantUuid")
		orgUuid := ctx.Params("organizationUuid")
		userUuid := ctx.Params("userUuid")

		c, span := otel.Tracer(OTEL_TRACER_NAME).Start(ctx.Context(), "API-USER-DELETE")
		defer span.End()

		userExists, errExists := userService.DeleteUser(tenantUuid, orgUuid, userUuid, c)
		if errExists != nil {
			var targetHttpStatus = commons.GuessHttpStatus(errExists)
			_ = ctx.SendStatus(targetHttpStatus)
			if commons.IsKnownFunctionalError(errExists) {
				apiErr := exceptions.ConvertToFunctionalError(errExists, targetHttpStatus)
				return ctx.JSON(apiErr)
			} else {
				apiErr := exceptions.ConvertToInternalError(errExists)
				return ctx.JSON(apiErr)
			}
		}
		if !userExists {
			apiErr := commons.ApiError{
				Code:         fiber.StatusNotFound,
				Kind:         string(commons.ErrorTypeFunctional),
				DebugMessage: commons.UserNotFound,
			}
			ctx.SendStatus(fiber.StatusNotFound)
			return ctx.JSON(apiErr)
		}

		ctx.SendStatus(fiber.StatusNoContent)
		return nil
	}
}

func MakeUsersFilter(userService services.UserService) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tenantUuid := ctx.Params("tenantUuid")
		orgUuid := ctx.Params("organizationUuid")

		c, span := otel.Tracer(OTEL_TRACER_NAME).Start(ctx.Context(), "API-USER-FILTER")
		defer span.End()

		searchFilter := ctx.Query("filter")
		page := ctx.Query("page")
		rowsPerPage := ctx.Query("rowsPerPage")
		sorting := ctx.Query("sort")

		searchExpressions, errorNodes, errListener := parser.FromInputString(searchFilter)
		if errorNodes != nil {
			apiErr := parser.ConvertErrorNodes(fiber.StatusBadRequest, errorNodes)
			ctx.SendStatus(fiber.StatusBadRequest)
			return ctx.JSON(apiErr)
		}
		if len(errListener.Errors) > 0 {
			apiErr := commons.ApiError{
				Code:         fiber.StatusBadRequest,
				Kind:         string(commons.ErrorTypeFunctional),
				Message:      commons.SearchFilter,
				DebugMessage: strings.Join(errListener.Errors, " | "),
			}
			ctx.SendStatus(fiber.StatusBadRequest)
			return ctx.JSON(apiErr)
		}

		pagination, errPagination := model.FromParameters(rowsPerPage, page, sorting)
		if errPagination != nil {
			apiErr := exceptions.ConvertToFunctionalError(errPagination, fiber.StatusBadRequest)
			return ctx.JSON(apiErr)
		}

		userListResponse, errList := userService.FilterUsers(tenantUuid, orgUuid, searchExpressions, pagination, c)
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
