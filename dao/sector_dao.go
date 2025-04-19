package dao

import (
	"context"
	"poc-fiber/logger"
	"poc-fiber/model"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
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

func (sectorDao *SectorDao) CreateSector(sector model.Sector, parentContext context.Context) (model.CompositeId, error) {
	_, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "FUNC-SECTOR-CREATE")
	defer span.End()

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

func (s SectorDao) WithTxCreateSector(tx pgx.Tx, sector model.Sector, parentContext context.Context) (model.CompositeId, error) {
	_, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "FUNC-SECTOR-CREATE_WITH_TX")
	defer span.End()

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

func (s SectorDao) FindAllByTenantAndOrganization(tenantId int64, organizationId int64, parentContext context.Context) ([]model.Sector, error) {
	c, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "FUNC-SECTOR-FIND_ALL_FOR_ORG")
	defer span.End()

	selStmt := viper.GetStringMapString(CONFIG_SECTORS)["findbytenantorg"]
	rows, e := s.DbPool.Query(context.Background(), selStmt, tenantId, organizationId)
	if e != nil {
		return nil, e
	}
	defer rows.Close()

	sectors, errCollect := pgx.CollectRows(rows, pgx.RowToStructByName[model.Sector])
	if errCollect != nil {
		return nil, errCollect
	}
	logger.LogRecord(c, LOGGER_NAME, "nb of results ["+strconv.Itoa(len(sectors))+"]")
	return sectors, nil
}

func (s SectorDao) FindByUuid(uuid string, parentContext context.Context) (model.Sector, error) {
	var nilSector model.Sector

	_, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "FUNC-SECTOR-FIND_BY_UUID")
	defer span.End()

	selStmt := viper.GetStringMapString(CONFIG_SECTORS)["findbyuuid"]
	rows, e := s.DbPool.Query(context.Background(), selStmt, uuid)
	if e != nil {
		return nilSector, e
	}
	defer rows.Close()
	sector, errCollect := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Sector])
	if errCollect != nil {
		return nilSector, errCollect
	}
	return sector, nil
}
