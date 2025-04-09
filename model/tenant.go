package model

type Tenant struct {
	Id    int64  `db:"id"`
	Uuid  string `db:"uuid"`
	Code  string `db:"code"`
	Label string `db:"label"`
}
