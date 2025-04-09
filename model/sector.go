package model

import "database/sql"

type Sector struct {
	Id             int64         `db:"id"`
	TenantId       int64         `db:"tenants_id"`
	OrganizationId int64         `db:"organizations_id"`
	Uuid           string        `db:"uuid"`
	Code           string        `db:"code"`
	Label          string        `db:"label"`
	ParentId       sql.NullInt64 `db:"parent_id"`
	HasParent      bool          `db:"has_parent"`
	Depth          int           `db:"depth"`
}
