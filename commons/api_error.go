package commons

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	TenantNotFound         = "func_tenant_not_found"
	OrgAlreadyExistsByCode = "func_org_code_already_used"
	OrgDoesNotExistByCode  = "func_org_does_not_exist"
	OrgNotFound            = "func_org_not_found"
	SectorAlreadyExist     = "func_sector_already_exists"
	SectorRootNotFound     = "func_sector_root_not_found"
	SectorNotFound         = "func_sector_not_found"
	UserNotFound           = "func_user_not_found"
	UserLoginAlreadyInUse  = "func_user_login_already_in_use"
	UserEmailAlreadyInUse  = "func_user_email_already_in_use"
	OAuthStateMismatch     = "func_oauth_state_mismatch"
)

type ApiErrorType string

const (
	ErrorTypeFunctional ApiErrorType = "functional"
	ErrorTypeTechnical  ApiErrorType = "technical"
)

type ApiError struct {
	Code         int               `json:"code"`
	Kind         string            `json:"kind"`
	Message      string            `json:"message"`
	DebugMessage string            `json:"debugMessage"`
	Details      []ApiErrorDetails `json:"details,omitempty"`
}

type ApiErrorDetails struct {
	Field  string `json:"field"`
	Detail string `json:"detail"`
}

func IsKnownFunctionalError(err error) bool {
	return strings.HasPrefix(err.Error(), "func_")
}

func GuessHttpStatus(err error) int {
	if IsKnownFunctionalError(err) {
		if IsNotFound(err) {
			return fiber.StatusNotFound
		} else {
			return fiber.StatusConflict
		}
	} else {
		return fiber.StatusInternalServerError
	}
}

func IsNotFound(err error) bool {
	return strings.HasSuffix(err.Error(), "not_found")
}
