package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"financeMgr/src/analysis-server/model"
)

type VoucherInfoDao struct {
	// Logger *log.Logger
	//Logger log.ILog
}

var (
	voucherInfoTN     = "voucherInfo"
	voucherInfoFields = []string{"voucher_id", "company_id", "voucher_month", "num_of_month",
		"voucher_filler", "voucher_auditor", "voucher_date", "bill_count", "status", "created_at", "updated_at"}
	scanVoucherInfo = func(r DbScanner, st *model.VoucherInfo) error {
		return r.Scan(&st.VoucherID, &st.CompanyID, &st.VoucherMonth, &st.NumOfMonth, &st.VoucherFiller,
			&st.VoucherAuditor, &st.VoucherDate, &st.BillCount, &st.Status, &st.CreatedAt, &st.UpdatedAt)
	}
)

func (dao *VoucherInfoDao) Get(ctx context.Context, do DbOperator, voucherId, iYear int) (*model.VoucherInfo, error) {
	strSql := "select " + strings.Join(voucherInfoFields, ",") + " from " +
		GenTableName(iYear, voucherInfoTN) + " where voucher_id=?"
	gLogger.DebugContext(ctx, "[VoucherInfo/db/Get] [sql: %s ,values: %d]", strSql, voucherId)
	var voucherInfo = &model.VoucherInfo{}
	start := time.Now()
	defer func() {
		gLogger.InfoContext(ctx, "[VoucherInfo/db/Get] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanVoucherInfo(do.QueryRowContext(ctx, strSql, voucherId), voucherInfo); err {
	case nil:
		return voucherInfo, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		gLogger.ErrorContext(ctx, "[VoucherInfo/db/Get] [scanVoucherInfo: %s]", err.Error())
		return nil, err
	}
}

// get the count of the table
func (dao *VoucherInfoDao) Count(ctx context.Context, do DbOperator, iYear int) (int64, error) {
	var c int64
	strSql := "select count(1) from " + GenTableName(iYear, voucherInfoTN)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql).Scan(&c)
	gLogger.InfoContext(ctx, "[VoucherInfo/db/Count] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

// list count by filter ...
func (dao *VoucherInfoDao) CountByFilter(ctx context.Context, do DbOperator, iYear int,
	filter map[string]interface{}) (int64, error) {
	var c int64
	strSql, values := makeCountSqlWithMultiCondition(GenTableName(iYear, voucherInfoTN), nil, filter, nil, nil)
	gLogger.DebugContext(ctx, "[VoucherInfo/db/CountByFilter] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, values...).Scan(&c)
	gLogger.InfoContext(ctx, "[voucherInfo/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

// 通过vouchrId 和 voucherMonth 获取该月份目前最大的凭证号
func (dao *VoucherInfoDao) GetMaxNumByIdAndMonth(ctx context.Context, do DbOperator,
	iYear, iVoucherMonth, iVoucherID int) (int64, error) {
	tableName := GenTableName(iYear, voucherInfoTN)
	strSql := "select count(*) from " + tableName +
		" where voucher_month=? and company_id in (select company_id from " + tableName + " where voucher_id=? )"
	gLogger.DebugContext(ctx, "[VoucherInfo/db/GetMaxNumByIdAndMonth] [sql: %s, values: %d-%d]",
		strSql, iVoucherMonth, iVoucherID)
	var c int64
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, iVoucherMonth, iVoucherID).Scan(&c)
	gLogger.InfoContext(ctx, "[voucherInfo/db/GetMaxNumByIdAndMonth] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

func (dao *VoucherInfoDao) Create(ctx context.Context, do DbOperator, st *model.VoucherInfo) error {
	tableName := GenTableName(st.VoucherDate/10000, voucherInfoTN)
	strSql := "insert into " + tableName + " (" + strings.Join(voucherInfoFields, ",") +
		") values (?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?)"
	values := []interface{}{st.VoucherID, st.CompanyID, st.VoucherMonth, st.NumOfMonth, st.VoucherFiller,
		st.VoucherAuditor, st.VoucherDate, st.BillCount, st.Status, st.CreatedAt, st.UpdatedAt}
	gLogger.DebugContext(ctx, "[VoucherInfo/db/Create] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	_, err := do.ExecContext(ctx, strSql, values...)
	gLogger.InfoContext(ctx, "[VoucherInfo/db/Create] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		gLogger.ErrorContext(ctx, "[VoucherInfo/db/Create] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}
func (dao *VoucherInfoDao) Delete(ctx context.Context, do DbOperator, voucherId, iYear int) error {
	strSql := "delete from " + GenTableName(iYear, voucherInfoTN) + " where voucher_id=?"

	gLogger.DebugContext(ctx, "[VoucherInfo/db/Delete] [sql: %s, id: %d]", strSql, voucherId)
	start := time.Now()
	defer func() {
		gLogger.InfoContext(ctx, "[VoucherInfo/db/Delete] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, voucherId); err != nil {
		gLogger.ErrorContext(ctx, "[VoucherInfo/db/Delete] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *VoucherInfoDao) BatchDelete(ctx context.Context, do DbOperator, iYear int, voucherIds []int) error {
	handleArrFilter := func(arr []int, s *string) (fv []interface{}) {
		for i, ki := range arr {
			if i == 0 {
				*s += "?"
			} else {
				*s += ", ?"
			}
			fv = append(fv, ki)
		}
		return
	}
	strSql := "delete from " + GenTableName(iYear, voucherInfoTN) + " where voucher_id IN ("
	fv := handleArrFilter(voucherIds, &strSql)
	strSql += ")"
	gLogger.DebugContext(ctx, "[VoucherInfo/db/BatchDelete] [sql: %s, ids: %v]", strSql, voucherIds)
	start := time.Now()
	defer func() {
		gLogger.InfoContext(ctx, "[VoucherInfo/db/BatchDelete] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, fv...); err != nil {
		gLogger.ErrorContext(ctx, "[VoucherInfo/db/BatchDelete] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

// 没有复杂的匹配条件
func (dao *VoucherInfoDao) SimpleList(ctx context.Context, do DbOperator, filter map[string]interface{},
	iYear, limit, offset, od int, order string) ([]*model.VoucherInfo, error) {
	var voucherInfoSlice []*model.VoucherInfo
	strSql, values := transferListSql(GenTableName(iYear, voucherInfoTN), filter, voucherInfoFields,
		limit, offset, order, od)
	gLogger.DebugContext(ctx, "[VoucherInfo/db/SimpleList] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		gLogger.InfoContext(ctx, "[VoucherInfo/db/SimpleList] [SqlElapsed: %v]", time.Since(start))
	}()
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		gLogger.ErrorContext(ctx, "[VoucherInfo/db/SimpleList] [do.Query: %s]", err.Error())
		return voucherInfoSlice, err
	}
	defer result.Close()
	for result.Next() {
		voucherInfo := new(model.VoucherInfo)
		err = scanVoucherInfo(result, voucherInfo)
		if err != nil {
			gLogger.ErrorContext(ctx, "[VoucherInfo/db/SimpleList] [scanVoucherInfo: %s]", err.Error())
			return voucherInfoSlice, err
		}
		voucherInfoSlice = append(voucherInfoSlice, voucherInfo)
	}
	return voucherInfoSlice, nil
}

// 1、在where条件里增加between ... and
// 2、该函数没有增加like，因为暂时没有用，所以没有增加该参数
// 3、增加了多列排序。
func (dao *VoucherInfoDao) List(ctx context.Context, do DbOperator, filterNo map[string]interface{},
	filter map[string]interface{}, intervalFilter map[string]interface{}, fuzzyMatchFilter map[string]string,
	orderFiler []*model.OrderItem, iYear, limit, offset int) ([]*model.VoucherInfo, error) {

	strSql, values := makeSelSqlWithMultiCondition(GenTableName(iYear, voucherInfoTN), voucherInfoFields,
		filterNo, filter, intervalFilter, fuzzyMatchFilter, orderFiler, limit, offset)
	gLogger.DebugContext(ctx, "[VoucherInfo/db/List] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		gLogger.InfoContext(ctx, "[VoucherInfo/db/List] [SqlElapsed: %v]", time.Since(start))
	}()
	var voucherInfoSlice []*model.VoucherInfo
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		gLogger.ErrorContext(ctx, "[VoucherInfo/db/List] [do.Query: %s]", err.Error())
		return voucherInfoSlice, err
	}
	defer result.Close()
	for result.Next() {
		voucherInfo := new(model.VoucherInfo)
		err = scanVoucherInfo(result, voucherInfo)
		if err != nil {
			gLogger.ErrorContext(ctx, "[VoucherInfo/db/List] [ScanSnapshot: %s]", err.Error())
			return voucherInfoSlice, err
		}
		voucherInfoSlice = append(voucherInfoSlice, voucherInfo)
	}
	return voucherInfoSlice, nil
}

// 通过公司ID以及根据登录日期获取最新的凭证信息
func (dao *VoucherInfoDao) GetLatestVoucherInfo(ctx context.Context, do DbOperator,
	iMonth, iYear, iCompanyID int) ([]*model.VoucherInfo, error) {
	tableName := GenTableName(iYear, voucherInfoTN)
	var voucherInfoSlice []*model.VoucherInfo
	strSql := "select " + strings.Join(voucherInfoFields, ",") + " from " + tableName +
		" where company_id = ? and voucher_month = (select  max(voucher_month) from " +
		tableName + " where company_id = ? and voucher_month <= ?) order by num_of_month "
	gLogger.DebugContext(ctx, "[VoucherInfo/db/GetLatestVoucherInfoByCompanyID] sql %s with values %v", strSql, iCompanyID, iMonth)
	start := time.Now()
	defer func() {
		gLogger.InfoContext(ctx, "[VoucherInfo/db/GetLatestVoucherInfoByCompanyID] [SqlElapsed: %v]", time.Since(start))
	}()
	result, err := do.QueryContext(ctx, strSql, iCompanyID, iCompanyID, iMonth)
	if err != nil {
		gLogger.ErrorContext(ctx, "[VoucherInfo/db/GetLatestVoucherInfoByCompanyID] [do.Query: %s]", err.Error())
		return voucherInfoSlice, err
	}
	defer result.Close()
	for result.Next() {
		voucherInfo := new(model.VoucherInfo)
		err = scanVoucherInfo(result, voucherInfo)
		if err != nil {
			gLogger.ErrorContext(ctx, "[VoucherInfo/db/GetLatestVoucherInfoByCompanyID] [scanVoucherInfo: %s]", err.Error())
			return voucherInfoSlice, err
		}
		voucherInfoSlice = append(voucherInfoSlice, voucherInfo)
	}
	return voucherInfoSlice, nil
}

func (dao *VoucherInfoDao) Update(ctx context.Context, do DbOperator, voucherId, iYear int,
	params map[string]interface{}) error {
	strSql := "update " + GenTableName(iYear, voucherInfoTN) + " set "
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
		//voucher_date 字段的类型为int,不需要转换成time.Time类型
		// if dbKey == "voucher_date" {
		// 	iDate, bOk := value.(int)
		// 	if !bOk {
		// 		continue
		// 	}
		// 	iYear := iDate / 10000
		// 	iMonth := (iDate - iYear*10000) / 100
		// 	iDay := iDate % 100
		// 	t := time.Date(iYear, time.Month(iMonth), iDay, 0, 0, 0, 0, time.Local)
		// 	value = t
		// }
		values = append(values, value)
	}
	if first {
		return nil
	}
	strSql += " where voucher_id=?"
	values = append(values, voucherId)
	start := time.Now()
	gLogger.DebugContext(ctx, "[VoucherInfo/db/Update] [sql: %s, values: %v]", strSql, values)
	_, err := do.ExecContext(ctx, strSql, values...)
	gLogger.InfoContext(ctx, "[VoucherInfo/db/Update] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		gLogger.ErrorContext(ctx, "[VoucherInfo/db/Update] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

// 用于批量审核/取消凭证。
func (dao *VoucherInfoDao) BatchUpdateVoucherInfoStatus(ctx context.Context, do DbOperator, iYear, iStatus int,
	strVoucherAuditor string, voucherIds []int) error {
	updateField := map[string]any{"status": iStatus, "voucher_auditor": strVoucherAuditor, "updated_at": time.Now()}
	filter := map[string]any{"voucher_id": voucherIds}
	strSql, values := makeUpdateSqlWithMultiCondition(GenTableName(iYear, voucherInfoTN), updateField, nil, filter, nil, nil)
	gLogger.DebugContext(ctx, "[VoucherInfo/db/BatchUpdate] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	defer func() {
		gLogger.InfoContext(ctx, "[VoucherInfo/db/BatchUpdate] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, values...); err != nil {
		gLogger.ErrorContext(ctx, "[VoucherInfo/db/BatchUpdate] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *VoucherInfoDao) GetMaxVoucherIdInLatestVoucherInfoTable(ctx context.Context, do DbOperator,
	latestVoucherInfoTable string) (int, error) {
	strSql := "select max(voucher_id) from " + latestVoucherInfoTable
	gLogger.DebugContext(ctx, "[VoucherInfo/db/GetMaxVoucherIdInLatestVoucherInfoTable] sql %s", strSql)
	var maxVoucherId sql.NullInt64
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql).Scan(&maxVoucherId)
	gLogger.InfoContext(ctx, "[VoucherInfo/db/GetMaxVoucherIdInLatestVoucherInfoTable] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		gLogger.ErrorContext(ctx, "[VoucherInfo/db/GetMaxVoucherIdInLatestVoucherInfoTable] [do.QueryRowContext: %s]", err.Error())
		return 0, err
	}
	if maxVoucherId.Valid {
		return int(maxVoucherId.Int64), nil
	} else {
		gLogger.WarnContext(ctx, "[VoucherInfo/db/GetMaxVoucherIdInLatestVoucherInfoTable] No voucher records found in table %s", latestVoucherInfoTable)
		gLogger.DebugContext(ctx, "[VoucherInfo/db/GetMaxVoucherIdInLatestVoucherInfoTable] [maxVoucherId: %v]", maxVoucherId)
		return 0, nil
	}
}
