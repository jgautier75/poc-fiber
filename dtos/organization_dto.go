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
