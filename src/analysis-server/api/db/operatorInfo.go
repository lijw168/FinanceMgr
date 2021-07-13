package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"analysis-server/model"
	"common/log"
)

type OperatorInfoDao struct {
	// Logger *log.Logger
	Logger log.ILog
}

var (
	operatorInfoTN     = "operatorInfo"
	operatorInfoFields = []string{"operator_id", "company_id", "name", "password", "job",
		"department", "status", "role", "created_at", "updated_at"}
	scanOperatorInfo = func(r DbScanner, st *model.OperatorInfo) error {
		return r.Scan(&st.OperatorID, &st.CompanyID, &st.Name, &st.Password,
			&st.Job, &st.Department, &st.Status, &st.Role, &st.CreatedAt, &st.UpdatedAt)
	}
)

func (dao *OperatorInfoDao) GetOptInfoByName(ctx context.Context, do DbOperator,
	strName string, iCompanyID int) (*model.OperatorInfo, error) {
	strSql := "select " + strings.Join(operatorInfoFields, ",") + " from " +
		operatorInfoTN + " where name=? and companyId=?"
	dao.Logger.DebugContext(ctx, "[OperatorInfo/db/GetOptInfoByName] [sql: %s ,name: %s,icompanyId:%d]", strSql, strName, iCompanyID)
	var optInfo = &model.OperatorInfo{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[OperatorInfo/db/GetOptInfoByName] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanOperatorInfo(do.QueryRowContext(ctx, strSql, strName, iCompanyID), optInfo); err {
	case nil:
		return optInfo, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[OperatorInfo/db/GetOptInfoByName] [scanOperatorInfo: %s]", err.Error())
		return nil, err
	}
}

func (dao *OperatorInfoDao) GetOptInfoById(ctx context.Context, do DbOperator,
	optID int) (*model.OperatorInfo, error) {
	strSql := "select " + strings.Join(operatorInfoFields, ",") + " from " + operatorInfoTN + " where operator_id=?"
	dao.Logger.DebugContext(ctx, "[OperatorInfo/db/GetOptInfoById] [sql: %s ,values: %d]", strSql, optID)
	var optInfo = &model.OperatorInfo{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[OperatorInfo/db/GetOptInfoById] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanOperatorInfo(do.QueryRowContext(ctx, strSql, optID), optInfo); err {
	case nil:
		return optInfo, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[OperatorInfo/db/GetOptInfoById] [scanOperatorInfo: %s]", err.Error())
		return nil, err
	}
}

//list count by filter
func (dao *OperatorInfoDao) CountByFilter(ctx context.Context, do DbOperator, filter map[string]interface{}) (int64, error) {
	var c int64
	strSql, values := transferCountSql(operatorInfoTN, filter)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, values...).Scan(&c)
	dao.Logger.InfoContext(ctx, "[OperatorInfo/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

func (dao *OperatorInfoDao) Create(ctx context.Context, do DbOperator, st *model.OperatorInfo) error {
	strSql := "insert into " + operatorInfoTN + " (" + strings.Join(operatorInfoFields, ",") +
		") values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	values := []interface{}{st.OperatorID, st.CompanyID, st.Name, st.Password, st.Job,
		st.Department, st.Status, st.Role, st.CreatedAt, st.UpdatedAt}
	dao.Logger.DebugContext(ctx, "[OperatorInfo/db/Create] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[OperatorInfo/db/Create] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[OperatorInfo/db/Create] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *OperatorInfoDao) Delete(ctx context.Context, do DbOperator, optID int) error {
	strSql := "delete from " + operatorInfoTN + " where operator_id = ?"

	dao.Logger.DebugContext(ctx, "[OperatorInfo/db/Delete] [sql: %s, id: %d]", strSql, optID)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[OperatorInfo/db/Delete] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, optID); err != nil {
		dao.Logger.ErrorContext(ctx, "[OperatorInfo/db/Delete] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *OperatorInfoDao) List(ctx context.Context, do DbOperator, filter map[string]interface{}, limit int,
	offset int, order string, od int) ([]*model.OperatorInfo, error) {
	var OperatorInfoSlice []*model.OperatorInfo
	strSql, values := transferListSql(operatorInfoTN, filter, operatorInfoFields, limit, offset, order, od)
	dao.Logger.DebugContext(ctx, "[OperatorInfo/db/List] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[OperatorInfo/db/List] [SqlElapsed: %v]", time.Since(start))
	}()
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[OperatorInfo/db/List] [do.Query: %s]", err.Error())
		return OperatorInfoSlice, err
	}
	defer result.Close()
	for result.Next() {
		OperatorInfo := new(model.OperatorInfo)
		err = scanOperatorInfo(result, OperatorInfo)
		if err != nil {
			dao.Logger.ErrorContext(ctx, "[OperatorInfo/db/List] [ScanSnapshot: %s]", err.Error())
			return OperatorInfoSlice, err
		}
		OperatorInfoSlice = append(OperatorInfoSlice, OperatorInfo)
	}
	return OperatorInfoSlice, nil
}

// func (dao *OperatorInfoDao) ListWithFilterNo(ctx context.Context, do DbOperator, filter map[string]interface{},
// 	filterNo map[string]interface{}, limit int, offset int, order string, od int) ([]*model.OperatorInfo, error) {
// 	var OperatorInfoSlice []*model.OperatorInfo
// 	strSql, values := transferListSqlWithNo(operatorInfoTN, filter, filterNo, operatorInfoFields, limit, offset, order, od)
// 	dao.Logger.DebugContext(ctx, "[OperatorInfo/db/ListWithFilterNo] sql %s with values %v", strSql, values)
// 	start := time.Now()
// 	defer func() {
// 		dao.Logger.InfoContext(ctx, "[OperatorInfo/db/ListWithFilterNo] [SqlElapsed: %v]", time.Since(start))
// 	}()
// 	result, err := do.QueryContext(ctx, strSql, values...)
// 	if err != nil {
// 		dao.Logger.ErrorContext(ctx, "[OperatorInfo/db/ListWithFilterNo] [do.Query: %s]", err.Error())
// 		return OperatorInfoSlice, err
// 	}
// 	defer result.Close()
// 	for result.Next() {
// 		OperatorInfo := new(model.OperatorInfo)
// 		err = scanOperatorInfo(result, OperatorInfo)
// 		if err != nil {
// 			dao.Logger.ErrorContext(ctx, "[OperatorInfo/db/ListWithFilterNo] [ScanSnapshot: %s]", err.Error())
// 			return OperatorInfoSlice, err
// 		}
// 		OperatorInfoSlice = append(OperatorInfoSlice, OperatorInfo)
// 	}
// 	return OperatorInfoSlice, nil
// }

func (dao *OperatorInfoDao) Update(ctx context.Context, do DbOperator, optID int,
	params map[string]interface{}) error {
	strSql := "update " + operatorInfoTN + " set "
	var values []interface{}
	var first bool = true
	for key, value := range params {
		dbKey := camelToUnix(key)
		if first {
			strSql += dbKey + "=?"
			first = false
		} else {
			strSql += "," + dbKey + "=?"
		}
		values = append(values, value)
	}
	if first {
		return nil
	}
	strSql += " where operator_id = ?"
	values = append(values, optID)
	start := time.Now()
	dao.Logger.DebugContext(ctx, "[OperatorInfo/db/Update] [sql: %s, values: %v]", strSql, values)
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[OperatorInfo/db/Update] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[OperatorInfo/db/Update] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}
