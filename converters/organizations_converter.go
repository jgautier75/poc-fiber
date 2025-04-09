package converters

import (
	"poc-fiber/dtos"
	"poc-fiber/model"
)

func ConvertOrgEntityToOrgLight(org model.Organization) dtos.OrgLightResponse {
	var resp = dtos.OrgLightResponse{
		Uuid:  &org.Uuid,
		Code:  &org.Code,
		Label: &org.Label,
		Type:  &org.Type,
	}
	return resp
}

func ConvertOrgEntityListToOrgLightList(orgs []model.Organization) dtos.OrgLightReponseList {
	var dtosList = make([]dtos.OrgLightResponse, len(orgs))
	for inc, org := range orgs {
		orgResponse := ConvertOrgEntityToOrgLight(org)
		dtosList[inc] = orgResponse
	}
	var responseList = dtos.OrgLightReponseList{
		Orgs: dtosList,
	}
	return responseList
}
