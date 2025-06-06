package db

import (
	"context"
	"database/sql"
	"financeMgr/src/common/log"
	"strings"
	"time"

	"financeMgr/src/analysis-server/model"
)

type VoucherRecordDao struct {
	// Logger *log.Logger
	Logger log.ILog
}

var (
	voucherRecordTN     = "voucherRecordInfo"
	voucherRecordFields = []string{"record_id", "voucher_id", "subject_name", "debit_money", "credit_money",
		"summary", "sub_id1", "sub_id2", "sub_id3", "sub_id4", "created_at", "updated_at"}
	scanVoucherRecord = func(r DbScanner, st *model.VoucherRecord) error {
		return r.Scan(&st.RecordID, &st.VoucherID, &st.SubjectName, &st.DebitMoney, &st.CreditMoney,
			&st.Summary, &st.SubID1, &st.SubID2, &st.SubID3, &st.SubID4, &st.CreatedAt, &st.UpdatedAt)
	}
)

func (dao *VoucherRecordDao) Get(ctx context.Context, do DbOperator, recordId,
	iYear int) (*model.VoucherRecord, error) {
	strSql := "select " + strings.Join(voucherRecordFields, ",") + " from " +
		GenTableName(iYear, voucherRecordTN) + " where record_id=?"
	dao.Logger.DebugContext(ctx, "[VoucherRecord/db/Get] [sql: %s ,values: %d]", strSql, recordId)
	var recInfo = &model.VoucherRecord{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[VoucherRecord/db/Get] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanVoucherRecord(do.QueryRowContext(ctx, strSql, recordId), recInfo); err {
	case nil:
		return recInfo, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[VoucherRecord/db/Get] [scanVoucherRecord: %s]", err.Error())
		return nil, err
	}
}

// get the count of the table
func (dao *VoucherRecordDao) Count(ctx context.Context, do DbOperator, iYear int) (int64, error) {
	var c int64
	strSql := "select count(1) from " + GenTableName(iYear, voucherRecordTN)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql).Scan(&c)
	dao.Logger.InfoContext(ctx, "[VoucherRecord/db/Count] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

// list count by filter
func (d *VoucherRecordDao) CountByFilter(ctx context.Context, do DbOperator, iYear int,
	filter map[string]interface{}) (int64, error) {
	var c int64
	strSql, values := transferCountSql(GenTableName(iYear, voucherRecordTN), filter)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, values...).Scan(&c)
	d.Logger.InfoContext(ctx, "[VoucherRecord/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

func (dao *VoucherRecordDao) Create(ctx context.Context, do DbOperator, iYear int, st *model.VoucherRecord) error {
	strSql := "insert into " + GenTableName(iYear, voucherRecordTN) + " (" + strings.Join(voucherRecordFields, ",") +
		") values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	values := []interface{}{st.RecordID, st.VoucherID, st.SubjectName, st.DebitMoney, st.CreditMoney,
		st.Summary, st.SubID1, st.SubID2, st.SubID3, st.SubID4, st.CreatedAt, st.UpdatedAt}
	dao.Logger.DebugContext(ctx, "[VoucherRecord/db/Create] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[VoucherRecord/db/Create] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherRecord/db/Create] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *VoucherRecordDao) Delete(ctx context.Context, do DbOperator, recordId, iYear int) error {
	strSql := "delete from " + GenTableName(iYear, voucherRecordTN) + " where record_id = ?"
	dao.Logger.DebugContext(ctx, "[VoucherRecord/db/Delete] [sql: %s, id: %s]", strSql, recordId)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[VoucherRecord/db/Delete] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, recordId); err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherRecord/db/Delete] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *VoucherRecordDao) DeleteByVoucherId(ctx context.Context, do DbOperator, voucherId, iYear int) error {
	strSql := "delete from " + GenTableName(iYear, voucherRecordTN) + " where voucher_id = ?"
	dao.Logger.DebugContext(ctx, "[VoucherRecord/db/DeleteByVoucherId] [sql: %s, id: %s]", strSql, voucherId)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[VoucherRecord/db/DeleteByVoucherId] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, voucherId); err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherRecord/db/DeleteByVoucherId] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *VoucherRecordDao) DeleteByMultiCondition(ctx context.Context, do DbOperator,
	iYear int, filter map[string]interface{}) error {
	strSql, values := transferDeleteSql(GenTableName(iYear, voucherRecordTN), filter)
	dao.Logger.DebugContext(ctx, "[VoucherRecord/db/DeleteByMultiCondition] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[VoucherRecord/db/DeleteByMultiCondition] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, values...); err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherRecord/db/DeleteByMultiCondition] [do.ExecContext: %s]", err.Error())
		return err
	}
	return nil
}

// 没有复杂的匹配条件。
func (dao *VoucherRecordDao) SimpleList(ctx context.Context, do DbOperator, filter map[string]interface{},
	iYear, limit, offset, od int, order string) ([]*model.VoucherRecord, error) {
	var voucherRecordSlice []*model.VoucherRecord
	strSql, values := transferListSql(GenTableName(iYear, voucherRecordTN), filter, voucherRecordFields,
		limit, offset, order, od)
	dao.Logger.DebugContext(ctx, "[VoucherRecord/db/SimpleList] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[VoucherRecord/db/SimpleList] [SqlElapsed: %v]", time.Since(start))
	}()
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherRecord/db/SimpleList] [do.Query: %s]", err.Error())
		return voucherRecordSlice, err
	}
	defer result.Close()
	for result.Next() {
		VoucherRecord := new(model.VoucherRecord)
		err = scanVoucherRecord(result, VoucherRecord)
		if err != nil {
			dao.Logger.ErrorContext(ctx, "[VoucherRecord/db/SimpleList] [ScanSnapshot: %s]", err.Error())
			return voucherRecordSlice, err
		}
		voucherRecordSlice = append(voucherRecordSlice, VoucherRecord)
	}
	return voucherRecordSlice, nil
}

// 1、在where条件里增加between ... and
// 2、增加了like
// 3、增加了多列排序。
func (dao *VoucherRecordDao) List(ctx context.Context, do DbOperator, filterNo map[string]interface{},
	filter map[string]interface{}, intervalFilter map[string]interface{}, fuzzyMatchFilter map[string]string,
	orderFiler []*model.OrderItem, iYear, limit, offset int) ([]*model.VoucherRecord, error) {
	strSql, values := makeSelSqlWithMultiCondition(GenTableName(iYear, voucherRecordTN), voucherRecordFields, filterNo,
		filter, intervalFilter, fuzzyMatchFilter, orderFiler, limit, offset)
	dao.Logger.DebugContext(ctx, "[VoucherRecord/db/List] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[VoucherRecord/db/List] [SqlElapsed: %v]", time.Since(start))
	}()
	var voucherRecordSlice []*model.VoucherRecord
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherRecord/db/List] [do.Query: %s]", err.Error())
		return voucherRecordSlice, err
	}
	defer result.Close()
	for result.Next() {
		VoucherRecord := new(model.VoucherRecord)
		err = scanVoucherRecord(result, VoucherRecord)
		if err != nil {
			dao.Logger.ErrorContext(ctx, "[VoucherRecord/db/List] [ScanSnapshot: %s]", err.Error())
			return voucherRecordSlice, err
		}
		voucherRecordSlice = append(voucherRecordSlice, VoucherRecord)
	}
	return voucherRecordSlice, nil
}

func (dao *VoucherRecordDao) UpdateByRecordId(ctx context.Context, do DbOperator, recordId, iYear int,
	params map[string]interface{}) error {
	// strSql := "update " + GenTableName(iYear, voucherRecordTN) + " set "
	// var values []interface{}
	// var first bool = true
	// for key, value := range params {
	// 	dbKey := camelToUnix(key)
	// 	if first {
	// 		strSql += dbKey + "=?"
	// 		first = false
	// 	} else {
	// 		strSql += "," + dbKey + "=?"
	// 	}
	// 	values = append(values, value)
	// }
	// if first {
	// 	return nil
	// }
	// strSql += " where record_id = ?"
	// values = append(values, recordId)
	filter := map[string]any{"record_id": recordId}
	strSql, values := makeUpdateSqlWithMultiCondition(GenTableName(iYear, voucherRecordTN), params, nil, filter, nil, nil)
	start := time.Now()
	dao.Logger.DebugContext(ctx, "[VoucherRecord/db/UpdateByRecordId] [sql: %s, values: %v]", strSql, values)
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[VoucherRecord/db/UpdateByRecordId] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherRecord/db/UpdateByRecordId] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *VoucherRecordDao) UpdateByVoucherId(ctx context.Context, do DbOperator, voucherId, iYear int,
	params map[string]interface{}) error {
	// strSql := "update " + GenTableName(iYear, voucherRecordTN) + " set "
	// var values []interface{}
	// var first bool = true
	// for key, value := range params {
	// 	dbKey := camelToUnix(key)
	// 	if first {
	// 		strSql += dbKey + "=?"
	// 		first = false
	// 	} else {
	// 		strSql += "," + dbKey + "=?"
	// 	}
	// 	values = append(values, value)
	// }
	// if first {
	// 	return nil
	// }
	// strSql += " where voucher_id = ?"
	// values = append(values, voucherId)
	filter := map[string]any{"voucher_id": voucherId}
	strSql, values := makeUpdateSqlWithMultiCondition(GenTableName(iYear, voucherRecordTN), params, nil, filter, nil, nil)
	start := time.Now()
	dao.Logger.DebugContext(ctx, "[VoucherRecord/db/UpdateByVoucherId] [sql: %s, values: %v]", strSql, values)
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[VoucherRecord/db/UpdateByVoucherId] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherRecord/db/UpdateByVoucherId] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}
