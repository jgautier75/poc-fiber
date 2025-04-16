package dtos

import "database/sql"

type SectorResponse struct {
	Id       sql.NullInt64    `json:"id,omitempty"`
	Uuid     *string          `json:"uuid"`
	Code     *string          `json:"code"`
	Label    *string          `json:"label"`
	Depth    int              `json:"depth"`
	Children []SectorResponse `json:"children,omitempty"`
	ParentId sql.NullInt64    `json:"parent_id,omitempty"`
}

type SectorLightResponse struct {
	Uuid     *string               `json:"uuid"`
	Code     *string               `json:"code"`
	Label    *string               `json:"label"`
	Depth    int                   `json:"depth"`
	Children []SectorLightResponse `json:"children,omitempty"`
}

type SectorResponseList struct {
	Sectors SectorResponse `json:"sectors"`
}

type SectorLightResponseList struct {
	Sectors SectorLightResponse `json:"sectors"`
}

type SectorCreateRequest struct {
	Code       *string `json:"code"`
	Label      *string `json:"label"`
	ParentUuid *string `json:"parent_uuid"`
}
