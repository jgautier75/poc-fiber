package dao

import (
	"context"
	"poc-fiber/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

const CONFIG_SECTORS = "sql.sectors"

type SectorDao struct {
	DbPool *pgxpool.Pool
}

func NewSectorDao(pool *pgxpool.Pool) SectorDao {
	sectorDao := SectorDao{}
	sectorDao.DbPool = pool
	return sectorDao
}

func (sectorDao *SectorDao) CreateSector(sector model.Sector) (model.CompositeId, error) {
	insertStmt := viper.GetStringMapString(CONFIG_SECTORS)["create"]
	nuuid := uuid.New().String()
	var id int64
	errQuery := sectorDao.DbPool.QueryRow(context.Background(), insertStmt, sector.TenantId, sector.OrganizationId, nuuid, sector.Code, sector.Label, sector.ParentId, sector.HasParent, sector.Depth).Scan(&id)
	compId := model.CompositeId{
		Id:   id,
		Uuid: nuuid,
	}
	return compId, errQuery
}

func (s SectorDao) WithTxCreateSector(tx pgx.Tx, sector model.Sector) (model.CompositeId, error) {
	insertStmt := viper.GetStringMapString(CONFIG_SECTORS)["create"]
	nuuid := uuid.New().String()
	var id int64
	errQuery := tx.QueryRow(context.Background(), insertStmt, sector.TenantId, sector.OrganizationId, nuuid, sector.Code, sector.Label, sector.ParentId, sector.HasParent, sector.Depth).Scan(&id)
	compId := model.CompositeId{
		Id:   id,
		Uuid: nuuid,
	}
	return compId, errQuery
}

func (s SectorDao) FindByTenantandOrganization(tenantId int64, organizationId int64) ([]model.Sector, error) {
	selStmt := viper.GetString("sectors.findbytenantorg")
	rows, e := s.DbPool.Query(context.Background(), selStmt, tenantId, organizationId)
	if e != nil {
		return nil, e
	}
	defer rows.Close()

	sectors, errCollect := pgx.CollectRows(rows, pgx.RowToStructByName[model.Sector])
	if errCollect != nil {
		return nil, errCollect
	}
	return sectors, nil
}
