package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"analysis-server/model"
	"common/log"
)

type VoucherInfoDao struct {
	// Logger *log.Logger
	Logger log.ILog
}

var (
	voucherInfoTN     = "voucherInfo"
	voucherInfoFields = []string{"voucher_id", "company_id", "voucher_month", "num_of_month",
		"voucher_filler", "voucher_auditor", "voucher_date", "created_at", "updated_at"}
	scanVoucherInfo = func(r DbScanner, st *model.VoucherInfo) error {
		return r.Scan(&st.VoucherID, &st.CompanyID, &st.VoucherMonth, &st.NumOfMonth, &st.VoucherFiller,
			&st.VoucherAuditor, &st.VoucherDate, &st.CreatedAt, &st.UpdatedAt)
	}
)

func (dao *VoucherInfoDao) Get(ctx context.Context, do DbOperator, voucherId int) (*model.VoucherInfo, error) {
	strSql := "select " + strings.Join(voucherInfoFields, ",") + " from " + voucherInfoTN + " where voucher_id=?"
	dao.Logger.DebugContext(ctx, "[VoucherInfo/db/Get] [sql: %s ,values: %d]", strSql, voucherId)
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
	dao.Logger.InfoContext(ctx, "[VoucherInfo/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
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

//通过vouchrId 和 voucherMonth 获取该月份目前最大的凭证号
func (dao *VoucherInfoDao) GetMaxNumByIdAndMonth(ctx context.Context, do DbOperator,
	iVoucherMonth, iVoucherID int) (int64, error) {
	strSql := "select count(*) from " + voucherInfoTN +
		" where voucher_month=? and company_id in (select company_id from " + voucherInfoTN + " where voucher_id=? )"
	var c int64
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, iVoucherMonth, iVoucherID).Scan(&c)
	dao.Logger.InfoContext(ctx, "[voucherInfo/db/GetMaxNumByIdAndMonth] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

func (dao *VoucherInfoDao) Create(ctx context.Context, do DbOperator, st *model.VoucherInfo) error {
	strSql := "insert into " + voucherInfoTN + " (" + strings.Join(voucherInfoFields, ",") +
		") values (?, ?, ?, ?, ? ,? ,?, ?, ?)"
	values := []interface{}{st.VoucherID, st.CompanyID, st.VoucherMonth, st.NumOfMonth, st.VoucherFiller,
		st.VoucherAuditor, st.VoucherDate, st.CreatedAt, st.UpdatedAt}
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
	strSql := "delete from " + voucherInfoTN + " where voucher_id=?"

	dao.Logger.DebugContext(ctx, "[VoucherInfo/db/Delete] [sql: %s, id: %d]", strSql, voucherId)
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
		err = scanVoucherInfo(result, voucherInfo)
		if err != nil {
			dao.Logger.ErrorContext(ctx, "[VoucherInfo/db/List] [ScanSnapshot: %s]", err.Error())
			return voucherInfoSlice, err
		}
		voucherInfoSlice = append(voucherInfoSlice, voucherInfo)
	}
	return voucherInfoSlice, nil
}

func (dao *VoucherInfoDao) GetLatestVoucherInfoByCompanyID(ctx context.Context, do DbOperator,
	iCompanyID int) ([]*model.VoucherInfo, error) {
	var voucherInfoSlice []*model.VoucherInfo
	strSql := "select " + strings.Join(voucherInfoFields, ",") + " from " + voucherInfoTN +
		" where voucher_month in (select  max(voucher_month) from " +
		voucherInfoTN + " where company_id = ?) order by num_of_month "
	dao.Logger.DebugContext(ctx, "[VoucherInfo/db/GetLatestVoucherInfoByCompanyID] sql %s with values %v", strSql, iCompanyID)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[VoucherInfo/db/GetLatestVoucherInfoByCompanyID] [SqlElapsed: %v]", time.Since(start))
	}()
	result, err := do.QueryContext(ctx, strSql, iCompanyID)
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherInfo/db/GetLatestVoucherInfoByCompanyID] [do.Query: %s]", err.Error())
		return voucherInfoSlice, err
	}
	defer result.Close()
	for result.Next() {
		voucherInfo := new(model.VoucherInfo)
		err = scanVoucherInfo(result, voucherInfo)
		if err != nil {
			dao.Logger.ErrorContext(ctx, "[VoucherInfo/db/GetLatestVoucherInfoByCompanyID] [scanVoucherInfo: %s]", err.Error())
			return voucherInfoSlice, err
		}
		voucherInfoSlice = append(voucherInfoSlice, voucherInfo)
	}
	return voucherInfoSlice, nil
}

func (dao *VoucherInfoDao) Update(ctx context.Context, do DbOperator, voucherId int,
	params map[string]interface{}) error {
	// var keyMap = map[string]string{"VoucherID": "voucher_id", "CompanyID": "company_id", "VoucherDate": "voucher_month",
	// 	"NumOfMonth": "num_of_month", "VoucherMonth": "voucher_date", "CreatedAt": "create_at", "UpdatedAt": "update_at"}
	strSql := "update " + voucherInfoTN + " set "
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
		if dbKey == "voucher_date" {
			iDate, bOk := value.(int)
			if !bOk {
				continue
			}
			iYear := iDate / 10000
			iMonth := (iDate - iYear*10000) / 100
			iDay := iDate % 100
			t := time.Date(iYear, time.Month(iMonth), iDay, 0, 0, 0, 0, time.Local)
			value = t
		}
		values = append(values, value)
	}
	if first {
		return nil
	}
	strSql += " where voucher_id=?"
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
