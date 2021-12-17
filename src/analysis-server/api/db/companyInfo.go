package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"analysis-server/model"
	"common/log"
)

type CompanyDao struct {
	// Logger *log.Logger
	Logger log.ILog
}

var (
	companyInfoTN     = "companyInfo"
	companyInfoFields = []string{"company_id", "company_name", "abbre_name", "corporator", "phone",
		"e_mail", "company_addr", "backup", "start_account_period", "latest_account_year", "created_at",
		"updated_at", "company_group_id"}
	scanCompanyInfo = func(r DbScanner, st *model.CompanyInfo) error {
		return r.Scan(&st.CompanyID, &st.CompanyName, &st.AbbrevName, &st.Corporator, &st.Phone,
			&st.Email, &st.CompanyAddr, &st.Backup, &st.StartAccountPeriod, &st.LatestAccountYear,
			&st.CreatedAt, &st.UpdatedAt, &st.CompanyGroupID)
	}
)

func (dao *CompanyDao) Get(ctx context.Context, do DbOperator, companyId int) (*model.CompanyInfo, error) {
	strSql := "select " + strings.Join(companyInfoFields, ",") + " from " + companyInfoTN + " where company_id=?"
	dao.Logger.DebugContext(ctx, "[CompanyInfo/db/Get] [sql: %s ,values: %d]", strSql, companyId)
	var compInfo = &model.CompanyInfo{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[CompanyInfo/db/Get] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanCompanyInfo(do.QueryRowContext(ctx, strSql, companyId), compInfo); err {
	case nil:
		return compInfo, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[CompanyInfo/db/Get] [scanCompanyInfo: %s]", err.Error())
		return nil, err
	}
}

func (dao *CompanyDao) GetCompanyByOperatorId(ctx context.Context, do DbOperator,
	operatorId int) (*model.CompanyInfo, error) {
	strSql := "select b." + strings.Join(companyInfoFields, ",b.") +
		" from operatorInfo as a, companyInfo as b where a.operator_id =? and a.company_id = b.company_id"
	dao.Logger.DebugContext(ctx, "[CompanyInfo/db/GetCompanyByOperatorId] [sql: %s ,values: %d]",
		strSql, operatorId)
	var compInfo = &model.CompanyInfo{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[CompanyInfo/db/GetCompanyByOperatorId] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanCompanyInfo(do.QueryRowContext(ctx, strSql, operatorId), compInfo); err {
	case nil:
		return compInfo, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[CompanyInfo/db/Get] [scanCompanyInfo: %s]", err.Error())
		return nil, err
	}
}

func (dao *CompanyDao) GetCompanyByAccSubId(ctx context.Context, do DbOperator,
	subjectId int) (*model.CompanyInfo, error) {
	strSql := "select b." + strings.Join(companyInfoFields, ",b.") +
		" from accountSubject as a, companyInfo as b where a.subject_id =? and a.company_id = b.company_id"
	dao.Logger.DebugContext(ctx, "[CompanyInfo/db/GetCompanyByAccSubId] [sql: %s ,values: %d]",
		strSql, subjectId)
	var compInfo = &model.CompanyInfo{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[CompanyInfo/db/GetCompanyByAccSubId] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanCompanyInfo(do.QueryRowContext(ctx, strSql, subjectId), compInfo); err {
	case nil:
		return compInfo, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[CompanyInfo/db/GetCompanyByAccSubId] [scanCompanyInfo: %s]", err.Error())
		return nil, err
	}
}

func (dao *CompanyDao) Create(ctx context.Context, do DbOperator, st *model.CompanyInfo) error {
	strSql := "insert into " + companyInfoTN + " (" + strings.Join(companyInfoFields, ",") +
		") values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	values := []interface{}{st.CompanyID, st.CompanyName, st.AbbrevName, st.Corporator, st.Phone,
		st.Email, st.CompanyAddr, st.Backup, st.StartAccountPeriod, st.LatestAccountYear,
		st.CreatedAt, st.UpdatedAt, st.CompanyGroupID}
	dao.Logger.DebugContext(ctx, "[CompanyInfo/db/Create] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[CompanyInfo/db/Create] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[CompanyInfo/db/Create] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *CompanyDao) Delete(ctx context.Context, do DbOperator, companyId int) error {
	strSql := "delete from " + companyInfoTN + " where company_id = ?"

	dao.Logger.DebugContext(ctx, "[CompanyInfo/db/Delete] [sql: %s, id: %d]", strSql, companyId)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[CompanyInfo/db/Delete] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, companyId); err != nil {
		dao.Logger.ErrorContext(ctx, "[CompanyInfo/db/Delete] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

//get the count of the table
func (dao *CompanyDao) Count(ctx context.Context, do DbOperator) (int64, error) {
	var c int64
	strSql := "select count(1) from " + companyInfoTN
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, nil).Scan(&c)
	dao.Logger.InfoContext(ctx, "[CompanyInfo/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

//list count by filter
func (dao *CompanyDao) CountByFilter(ctx context.Context, do DbOperator, filter map[string]interface{}) (int64, error) {
	var c int64
	strSql, values := transferCountSql(companyInfoTN, filter)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, values...).Scan(&c)
	dao.Logger.InfoContext(ctx, "[CompanyInfo/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

func (dao *CompanyDao) List(ctx context.Context, do DbOperator, filter map[string]interface{}, limit int,
	offset int, order string, od int) ([]*model.CompanyInfo, error) {
	var companyInfoSlice []*model.CompanyInfo
	strSql, values := transferListSql(companyInfoTN, filter, companyInfoFields, limit, offset, order, od)
	dao.Logger.DebugContext(ctx, "[CompanyInfo/db/List] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[CompanyInfo/db/List] [SqlElapsed: %v]", time.Since(start))
	}()
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[CompanyInfo/db/List] [do.Query: %s]", err.Error())
		return companyInfoSlice, err
	}
	defer result.Close()
	for result.Next() {
		companyInfo := new(model.CompanyInfo)
		err = scanCompanyInfo(result, companyInfo)
		if err != nil {
			dao.Logger.ErrorContext(ctx, "[CompanyInfo/db/List] [ScanSnapshot: %s]", err.Error())
			return companyInfoSlice, err
		}
		companyInfoSlice = append(companyInfoSlice, companyInfo)
	}
	return companyInfoSlice, nil
}

func (dao *CompanyDao) Update(ctx context.Context, do DbOperator, companyId int,
	params map[string]interface{}) error {
	strSql := "update " + companyInfoTN + " set "
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
	strSql += " where company_id = ?"
	values = append(values, companyId)
	start := time.Now()
	dao.Logger.DebugContext(ctx, "[CompanyInfo/db/Update] [sql: %s, values: %v]", strSql, values)
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[CompanyInfo/db/Update] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[CompanyInfo/db/Update] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

//create voucherInfo/voucherRecordInfo ...
func (dao *CompanyDao) CreateNewTable(ctx context.Context, do DbOperator, oldTableName, newTableName string) error {
	dao.Logger.DebugContext(ctx, "[CompanyInfo/db/createNewTable] [oldTableName: %s, newTableName: %s]", oldTableName, newTableName)
	//judge ,is not exist
	var c int64
	strSql := "select count(1) from information_schema.TABLES where table_name = ?"
	err := do.QueryRowContext(ctx, strSql, newTableName).Scan(&c)
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[CompanyInfo/db/CreateNewTable] [do.QueryRowContext: %s]", err.Error())
		return err
	}
	if c == 0 {
		//create new table
		strCreateTableSql := fmt.Sprintf("create table %s like %s", newTableName, oldTableName)
		if _, err = do.ExecContext(ctx, strCreateTableSql); err != nil {
			dao.Logger.ErrorContext(ctx, "[CompanyInfo/db/CreateNewTable] [do.Exec: %s]", err.Error())
		}
	}
	return err
}
