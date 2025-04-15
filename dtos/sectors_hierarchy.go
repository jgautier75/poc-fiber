package dtos

type SectorResponse struct {
	Id       *int64
	Uuid     *string          `json:"uuid"`
	Code     *string          `json:"code"`
	Label    *string          `json:"label"`
	Depth    int              `json:"depth"`
	Children []SectorResponse `json:"children,omitempty"`
	ParentId *int64
}

type SectorResponseList struct {
	Sectors SectorResponse `json:"sectors"`
}
