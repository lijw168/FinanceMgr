package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"common/log"
	"analysis-server/model"
)

type VoucherInfoDao struct {
	// Logger *log.Logger
	Logger log.ILog
}

var (
	voucherInfoTN     = "voucherInfo_2020"
	voucherInfoFields = []string{"voucherId", "companyId", "voucherMonth", "numOfMonth", "voucherDate"}
	scanVoucherInfo   = func(r DbScanner, st *model.VoucherInfo) error {
		return r.Scan(&st.VoucherID, &st.CompanyID, &st.VoucherMonth, &st.NumOfMonth, &st.VoucherDate)
	}
)

func (dao *VoucherInfoDao) Get(ctx context.Context, do DbOperator, voucherId int) (*model.VoucherInfo, error) {
	strSql := "select " + strings.Join(voucherInfoFields, ",") + " from " + voucherInfoTN + " where voucherId=?"
	dao.Logger.DebugContext(ctx, "[VoucherInfo/db/Get] [sql: %s ,values: %s]", strSql, voucherId)
	var voucherInfo = &model.VoucherInfo{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[VoucherInfo/db/Get] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanVoucherInfo(do.QueryRowContext(ctx, strSql, voucherId), voucherInfo); err {
	case nil:
		return voucherInfo, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[VoucherInfo/db/Get] [scanVoucherInfo: %s]", err.Error())
		return nil, err
	}
}

//get the count of the table
func (dao *VoucherInfoDao) Count(ctx context.Context, do DbOperator) (int64, error) {
	var c int64
	strSql := "select count(1) from " + voucherInfoTN
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, nil).Scan(&c)
	dao.Logger.InfoContext(ctx, "[accountSubject/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

//list count by filter
func (dao *VoucherInfoDao) CountByFilter(ctx context.Context, do DbOperator, filter map[string]interface{}) (int64, error) {
	var c int64
	strSql, values := transferCountSql(voucherInfoTN, filter)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, values...).Scan(&c)
	dao.Logger.InfoContext(ctx, "[voucherInfo/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

//voucherID由数据库中的自增ID产生，所以此时就不用再对该字段进行赋值了。
func (dao *VoucherInfoDao) Create(ctx context.Context, do DbOperator, st *model.VoucherInfo) error {
	strSql := "insert into " + voucherInfoTN + " (" + strings.Join(voucherInfoFields, ",") + ") 
			   values (?, ?, ?, ?)"
	//values := []interface{}{st.VoucherID, st.CompanyID, st.VoucherMonth, st.NumOfMonth, st.VoucherDate}
	values := []interface{}{st.CompanyID, st.VoucherMonth, st.NumOfMonth, st.VoucherDate}
	dao.Logger.DebugContext(ctx, "[VoucherInfo/db/Create] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[VoucherInfo/db/Create] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherInfo/db/Create] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}
func (dao *VoucherInfoDao) Delete(ctx context.Context, do DbOperator, voucherId int) error {
	strSql := "delete from " + voucherInfoTN + " where voucherId=?"

	dao.Logger.DebugContext(ctx, "[VoucherInfo/db/Delete] [sql: %s, id: %s]", strSql, voucherId)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[VoucherInfo/db/Delete] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, voucherId); err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherInfo/db/Delete] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}
func (dao *VoucherInfoDao) List(ctx context.Context, do DbOperator, filter map[string]interface{}, limit int,
								offset int, order string, od int) ([]*model.VoucherInfo, error) {
	var voucherInfoSlice []*model.VoucherInfo
	strSql, values := transferListSql(voucherInfoTN, filter, voucherInfoFields, limit, offset, order, od)
	dao.Logger.DebugContext(ctx, "[VoucherInfo/db/List] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[VoucherInfo/db/List] [SqlElapsed: %v]", time.Since(start))
	}()
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherInfo/db/List] [do.Query: %s]", err.Error())
		return voucherInfoSlice, err
	}
	defer result.Close()
	for result.Next() {
		voucherInfo := new(model.VoucherInfo)
		err = scanVoucherInfo(result, VoucherInfo)
		if err != nil {
			dao.Logger.ErrorContext(ctx, "[VoucherInfo/db/List] [ScanSnapshot: %s]", err.Error())
			return voucherInfoSlice, err
		}
		voucherInfoSlice = append(voucherInfoSlice, voucherInfo)
	}
	return voucherInfoSlice, nil
}

func (dao *VoucherInfoDao) Update(ctx context.Context, do DbOperator, voucherId int,
	                               params map[string]interface{}) error {
	var keyMap = map[string]string{"VoucherID": "voucherId","CompanyID": "companyId", "VoucherDate": "voucherDate",
									"NumOfMonth": "numOfMonth", "VoucherMonth": "voucherMonth"}
	strSql := "update " + voucherInfoTN + " set "
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
	strSql += " where voucherId=?"
	values = append(values, voucherId)
	start := time.Now()
	dao.Logger.DebugContext(ctx, "[VoucherInfo/db/Update] [sql: %s, values: %v]", strSql, values)
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[VoucherInfo/db/Update] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherInfo/db/Update] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}
