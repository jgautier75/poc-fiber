package exceptions

import (
	"poc-fiber/commons"
	"poc-fiber/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

const (
	TENANT_NOT_FOUND = "tenant_not_found"
	ORG_NOT_FOUND    = "organization_not_found"
	SECTOR_NOT_FOUND = "sector_not_found"
)

func ConvertToInternalError(err error) commons.ApiError {
	errStack := errors.WithStack(err)
	return commons.ApiError{
		Code:         fiber.StatusInternalServerError,
		Kind:         string(commons.ErrorTypeTechnical),
		Message:      err.Error(),
		DebugMessage: errStack.Error(),
	}
}

func ConvertToFunctionalError(err error, targetStatus int) commons.ApiError {
	return commons.ApiError{
		Code:    targetStatus,
		Kind:    string(commons.ErrorTypeFunctional),
		Message: err.Error(),
	}
}

func ConvertValidationError(errors []validation.ErrorValidation) commons.ApiError {
	var details []commons.ApiErrorDetails
	for _, e := range errors {
		details = append(details, commons.ApiErrorDetails{Field: e.Field, Detail: e.Error})
	}
	return commons.ApiError{
		Code:    fiber.StatusBadRequest,
		Kind:    string(commons.ErrorTypeFunctional),
		Message: validation.GlobalValidationFailed,
		Details: details,
	}
}
