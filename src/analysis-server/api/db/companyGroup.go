package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"financeMgr/src/analysis-server/model"
)

type CompanyGroupDao struct {
	//Logger log.ILog
}

var (
	companyGroupTN     = "companyGroup"
	companyGroupFields = []string{"company_group_id", "group_name", "group_status", "created_at", "updated_at"}
	scanCompanyGroup   = func(r DbScanner, st *model.CompanyGroup) error {
		return r.Scan(&st.CompanyGroupID, &st.GroupName, &st.GroupStatus, &st.CreatedAt, &st.UpdatedAt)
	}
)

func (dao *CompanyGroupDao) Get(ctx context.Context, do DbOperator, companyGroupId int) (*model.CompanyGroup, error) {
	strSql := "select " + strings.Join(companyGroupFields, ",") + " from " + companyGroupTN + " where company_group_id=?"
	gLogger.DebugContext(ctx, "[CompanyGroup/db/Get] [sql: %s ,values: %d]", strSql, companyGroupId)
	var compInfo = &model.CompanyGroup{}
	start := time.Now()
	defer func() {
		gLogger.InfoContext(ctx, "[CompanyGroup/db/Get] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanCompanyGroup(do.QueryRowContext(ctx, strSql, companyGroupId), compInfo); err {
	case nil:
		return compInfo, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		gLogger.ErrorContext(ctx, "[CompanyGroup/db/Get] [scanCompanyGroup: %s]", err.Error())
		return nil, err
	}
}

func (dao *CompanyGroupDao) Create(ctx context.Context, do DbOperator, st *model.CompanyGroup) error {
	strSql := "insert into " + companyGroupTN + " (" + strings.Join(companyGroupFields, ",") +
		") values (?, ?, ?, ?, ?)"
	values := []interface{}{st.CompanyGroupID, st.GroupName, st.GroupStatus, st.CreatedAt, st.UpdatedAt}
	gLogger.DebugContext(ctx, "[CompanyGroup/db/Create] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	_, err := do.ExecContext(ctx, strSql, values...)
	gLogger.InfoContext(ctx, "[CompanyGroup/db/Create] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		gLogger.ErrorContext(ctx, "[CompanyGroup/db/Create] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *CompanyGroupDao) Delete(ctx context.Context, do DbOperator, companyGroupId int) error {
	strSql := "delete from " + companyGroupTN + " where company_group_id = ?"

	gLogger.DebugContext(ctx, "[CompanyGroup/db/Delete] [sql: %s, id: %d]", strSql, companyGroupId)
	start := time.Now()
	defer func() {
		gLogger.InfoContext(ctx, "[CompanyGroup/db/Delete] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, companyGroupId); err != nil {
		gLogger.ErrorContext(ctx, "[CompanyGroup/db/Delete] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

// get the count of the table
func (dao *CompanyGroupDao) Count(ctx context.Context, do DbOperator) (int64, error) {
	var c int64
	strSql := "select count(1) from " + companyGroupTN
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql).Scan(&c)
	gLogger.InfoContext(ctx, "[CompanyGroup/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

// list count by filter
func (dao *CompanyGroupDao) CountByFilter(ctx context.Context, do DbOperator,
	filter map[string]interface{}) (int64, error) {
	var c int64
	start := time.Now()
	strSql, values := transferCountSql(companyGroupTN, filter)
	gLogger.DebugContext(ctx, "[CompanyGroup/db/CountByFilter] sql %s with values %v", strSql, values)
	err := do.QueryRowContext(ctx, strSql, values...).Scan(&c)
	gLogger.InfoContext(ctx, "[CompanyGroup/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

func (dao *CompanyGroupDao) List(ctx context.Context, do DbOperator, filter map[string]interface{}, limit int,
	offset int, order string, od int) ([]*model.CompanyGroup, error) {
	var companyGroupSlice []*model.CompanyGroup
	strSql, values := transferListSql(companyGroupTN, filter, companyGroupFields, limit, offset, order, od)
	gLogger.DebugContext(ctx, "[CompanyGroup/db/List] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		gLogger.InfoContext(ctx, "[CompanyGroup/db/List] [SqlElapsed: %v]", time.Since(start))
	}()
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		gLogger.ErrorContext(ctx, "[CompanyGroup/db/List] [do.Query: %s]", err.Error())
		return companyGroupSlice, err
	}
	defer result.Close()
	for result.Next() {
		companyInfo := new(model.CompanyGroup)
		err = scanCompanyGroup(result, companyInfo)
		if err != nil {
			gLogger.ErrorContext(ctx, "[CompanyGroup/db/List] [ScanSnapshot: %s]", err.Error())
			return companyGroupSlice, err
		}
		companyGroupSlice = append(companyGroupSlice, companyInfo)
	}
	return companyGroupSlice, nil
}

func (dao *CompanyGroupDao) Update(ctx context.Context, do DbOperator, companyGroupId int,
	params map[string]interface{}) error {

	filter := map[string]any{"company_group_id": companyGroupId}
	strSql, values := makeUpdateSqlWithMultiCondition(companyGroupTN, params, nil, filter, nil, nil)
	start := time.Now()
	gLogger.DebugContext(ctx, "[CompanyGroup/db/Update] [sql: %s, values: %v]", strSql, values)
	_, err := do.ExecContext(ctx, strSql, values...)
	gLogger.InfoContext(ctx, "[CompanyGroup/db/Update] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		gLogger.ErrorContext(ctx, "[CompanyGroup/db/Update] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

// get max companyGroupId
func (dao *CompanyGroupDao) GetMaxCompanyGroupId(ctx context.Context, do DbOperator) (int, error) {
	//use NullInt64 to avoid the problem that there is no data in the table and sql.ErrNoRows error will be returned,
	// so just use 0 as default value when there is no data in the table
	var maxCompanyGroupId sql.NullInt64
	strSql := "select max(company_group_id) from " + companyGroupTN
	gLogger.DebugContext(ctx, "[CompanyGroup/db/GetMaxCompanyGroupId] [sql: %s]", strSql)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql).Scan(&maxCompanyGroupId)
	gLogger.InfoContext(ctx, "[CompanyGroup/db/GetMaxCompanyGroupId] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		gLogger.ErrorContext(ctx, "[CompanyGroup/db/GetMaxCompanyGroupId] [do.QueryRowContext: %s]", err.Error())
		return 0, err
	}
	if maxCompanyGroupId.Valid {
		return int(maxCompanyGroupId.Int64), nil
	}
	return 0, nil
}
