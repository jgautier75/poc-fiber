package commons

const (
	OrgAlreadyExistsByCode  = "org_already_exists"
	OrgAlreadyExistsByLabel = "org_already_label"
	OrgDoesNotExistByCode   = "org_does_not_exist"
	OrgNotFound             = "org_not_found"
	SectorAlreadyExist      = "sector_already_exists"
	SectorRootNotFound      = "sector_root_not_found"
	SectorNotFound          = "sector_not_found"
	UserNotFound            = "user_not_found"
	UserLoginAlreadyInUse   = "user_login_already_in_use"
	UserEmailAlreadyInUse   = "user_email_already_in_use"
	OAuthStateMismatch      = "oauth_state_mismatch"
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
