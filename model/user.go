package model

type User struct {
	Id             int64  `db:"id"`
	TenantId       int64  `db:"tenants_id"`
	OrganizationId int64  `db:"organizations_id"`
	Uuid           string `db:"uuid"`
	LastName       string `db:"last_name"`
	FirstName      string `db:"first_name"`
	Login          string `db:"login"`
	Email          string `db:"email"`
}
