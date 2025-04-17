package dtos

type OrgLightResponse struct {
	Uuid  *string `json:"uuid"`
	Code  *string `json:"code"`
	Label *string `json:"label"`
	Type  *string `json:"type"`
}

type OrgLightReponseList struct {
	Orgs []OrgLightResponse `json:"orgs"`
}

type CreateOrgRequest struct {
	Code  *string `json:"code" validate:"required,max=50"`
	Label *string `json:"label" validate:"required,max=50"`
	Type  *string `json:"type" validate:"required,max=10"`
}
