package db

import (
	//"context"
	"database/sql"
	"strings"
	"time"

	"analysis-server/model"
	"common/log"
)

type IDInfoDao struct {
	// Logger *log.Logger
	Logger log.ILog
}

var (
	idInfoTN     = "idInfo"
	idInfoFields = []string{"company_id", "operator_id", "subject_id", "voucher_id", "voucher_record_id"}
	scanIdInfo   = func(r DbScanner, st *model.IDInfo) error {
		return r.Scan(&st.CompanyID, &st.OperatorID, &st.SubjectID, &st.VoucherID, &st.VoucherRecordID)
	}
)

func (dao *IDInfoDao) Get(do DbOperator) (*model.IDInfo, error) {
	strSql := "select " + strings.Join(idInfoFields, ",") + " from " + idInfoTN
	dao.Logger.Debug("[IDInfo/db/Get] [sql: %s ]", strSql)
	var idInfo = &model.IDInfo{}
	start := time.Now()
	defer func() {
		dao.Logger.Info("[IDInfo/db/Get] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanIdInfo(do.QueryRow(strSql), idInfo); err {
	case nil:
		return idInfo, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.Error("[IDInfo/db/Get] [scanIdInfo: %s]", err.Error())
		return nil, err
	}
}

func (dao *IDInfoDao) Create(do DbOperator, st *model.IDInfo) error {
	strSql := "insert into " + idInfoTN + " (" + strings.Join(idInfoFields, ",") +
		") values (?, ?, ?, ?, ?)"
	values := []interface{}{st.CompanyID, st.OperatorID, st.SubjectID, st.VoucherID, st.VoucherRecordID}
	dao.Logger.Debug("[IDInfo/db/Create] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	_, err := do.Exec(strSql, values...)
	dao.Logger.Info("[IDInfo/db/Create] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.Error("[IDInfo/db/Create] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *IDInfoDao) Delete(do DbOperator) error {
	strSql := "delete from " + idInfoTN

	dao.Logger.Debug("[IDInfo/db/Delete] [sql: %s]", strSql)
	start := time.Now()
	defer func() {
		dao.Logger.Info("[IDInfo/db/Delete] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.Exec(strSql, nil); err != nil {
		dao.Logger.Error("[IDInfo/db/Delete] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

//get the count of the table
func (dao *IDInfoDao) Count(do DbOperator) (int64, error) {
	var c int64
	strSql := "select count(1) from " + idInfoTN
	start := time.Now()
	err := do.QueryRow(strSql, nil).Scan(&c)
	dao.Logger.Info("[IDInfo/db/Count] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

func (dao *IDInfoDao) Update(do DbOperator, params map[string]interface{}) error {
	strSql := "update " + idInfoTN + " set "
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
	// strSql += " where subjectId = ?"
	// values = append(values, strSubID)
	start := time.Now()
	dao.Logger.Debug("[IDInfoDao/db/Update] [sql: %s, values: %v]", strSql, values)
	_, err := do.Exec(strSql, values...)
	dao.Logger.Info("[IDInfoDao/db/Update] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.Error("[IDInfoDao/db/Update] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}
