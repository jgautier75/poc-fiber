package model

type Organization struct {
	Id       int64  `db:"id"`
	TenantId int64  `db:"tenants_id"`
	Uuid     string `db:"uuid"`
	Code     string `db:"code"`
	Label    string `db:"label"`
	Type     string `db:"type"`
}
