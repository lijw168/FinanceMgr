package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"common/log"
	"analysis-server/model"
)

type CompanyDao struct {
	// Logger *log.Logger
	Logger log.ILog
}

var (
	companyInfoTN     = "companyInfo"
	companyInfoFields = []string{"companyId", "companyName", "abbreviationName", "corporator", "phone", 
								"e_mail", "companyAddr", "backup"}
	scanCompanyInfo   = func(r DbScanner, st *model.CompanyInfo) error {
		return r.Scan(&st.CompanyID, &st.CompanyName, &st.AbbrevName, &st.Corporator, &st.Phone,
					 &st.Summary, &st.Email, &st.CompanyAddr, &st.Backup)
	}
)

func (dao *CompanyDao) Get(ctx context.Context, do DbOperator, companyId int) (*model.CompanyInfo, error) {
	strSql := "select " + strings.Join(companyInfoFields, ",") + " from " + companyInfoTN + " where companyId=?"
	dao.Logger.DebugContext(ctx, "[CompanyInfo/db/Get] [sql: %s ,values: %s]", strSql, companyId)
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

//list count by filter
// func (d *CompanyDao) CountByFilter(ctx context.Context, do DbOperator, filter map[string]interface{}) (int64, error) {
// 	var c int64
// 	strSql, values := transferCountSql(companyInfoTN, filter)
// 	start := time.Now()
// 	err := do.QueryRowContext(ctx, strSql, values...).Scan(&c)
// 	d.Logger.InfoContext(ctx, "[CompanyInfo/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
// 	return c, err
// }

func (dao *CompanyDao) Create(ctx context.Context, do DbOperator, st *model.CompanyInfo) error {
	strSql := "insert into " + companyInfoTN + " (" + strings.Join(companyInfoFields, ",") + 
				") values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	values := []interface{}{st.CompanyID, st.CompanyName, st.AbbrevName, st.Corporator, st.Phone,
		                    st.Summary, st.Email,st.CompanyAddr,st.Backup}
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
	strSql := "delete from " + companyInfoTN + " where companyId = ?"

	dao.Logger.DebugContext(ctx, "[CompanyInfo/db/Delete] [sql: %s, id: %s]", strSql, companyId)
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
	var keyMap = map[string]string{"CompanyID": "companyId","CompanyName": "companyName", "AbbreviationName": "abbreviationName",
							       "Corporator": "corporator", "Phone": "phone", "E_mail": "e_mail",
							       "CompanyAddr": "companyAddr","Backup": "backup"}
	strSql := "update " + companyInfoTN + " set "
	var values []interface{}
	var first bool = true
	for key, value := range params {
		if dbKey, ok := keyMap[key]; ok {
			if first {
				strSql += dbKey + "=?"
				first = false
			} else {
				strSql += "," + dbKey + "=?"
			}
			values = append(values, value)
		}
	}
	if first {
		return nil
	}
	strSql += " where companyId = ?"
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
