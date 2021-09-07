package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"analysis-server/model"
	"common/log"
)

type VoucherRecordDao struct {
	// Logger *log.Logger
	Logger log.ILog
}

var (
	voucherRecordTN     = "voucherRecordInfo"
	voucherRecordFields = []string{"record_id", "voucher_id", "subject_name", "debit_money", "credit_money",
		"summary", "sub_id1", "sub_id2", "sub_id3", "sub_id4", "bill_count", "status", "created_at", "updated_at"}
	scanVoucherRecord = func(r DbScanner, st *model.VoucherRecord) error {
		return r.Scan(&st.RecordID, &st.VoucherID, &st.SubjectName, &st.DebitMoney, &st.CreditMoney,
			&st.Summary, &st.SubID1, &st.SubID2, &st.SubID3, &st.SubID4, &st.BillCount, &st.Status, &st.CreatedAt, &st.UpdatedAt)
	}
)

func (dao *VoucherRecordDao) Get(ctx context.Context, do DbOperator, recordId int) (*model.VoucherRecord, error) {
	strSql := "select " + strings.Join(voucherRecordFields, ",") + " from " + voucherRecordTN + " where record_id=?"
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

//get the count of the table
func (dao *VoucherRecordDao) Count(ctx context.Context, do DbOperator) (int64, error) {
	var c int64
	strSql := "select count(1) from " + voucherInfoTN
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, nil).Scan(&c)
	dao.Logger.InfoContext(ctx, "[VoucherRecord/db/Count] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

//list count by filter
func (d *VoucherRecordDao) CountByFilter(ctx context.Context, do DbOperator, filter map[string]interface{}) (int64, error) {
	var c int64
	strSql, values := transferCountSql(voucherRecordTN, filter)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, values...).Scan(&c)
	d.Logger.InfoContext(ctx, "[VoucherRecord/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

func (dao *VoucherRecordDao) Create(ctx context.Context, do DbOperator, st *model.VoucherRecord) error {
	strSql := "insert into " + voucherRecordTN + " (" + strings.Join(voucherRecordFields, ",") +
		") values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	values := []interface{}{st.RecordID, st.VoucherID, st.SubjectName, st.DebitMoney, st.CreditMoney,
		st.Summary, st.SubID1, st.SubID2, st.SubID3, st.SubID4, st.BillCount, st.Status, st.CreatedAt, st.UpdatedAt}
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

func (dao *VoucherRecordDao) Delete(ctx context.Context, do DbOperator, recordId int) error {
	strSql := "delete from " + voucherRecordTN + " where record_id = ?"

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

func (dao *VoucherRecordDao) DeleteByVoucherId(ctx context.Context, do DbOperator, voucherId int) error {
	strSql := "delete from " + voucherRecordTN + " where voucher_id = ?"

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
	filter map[string]interface{}) error {
	strSql, values := transferDeleteSql(voucherRecordTN, filter)
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

//1、在where条件里增加between ... and
//2、增加了like
func (dao *VoucherRecordDao) List(ctx context.Context, do DbOperator, filter map[string]interface{},
	intervalFilter map[string]interface{}, fuzzyMatchFilter map[string]string, limit int,
	offset int, order string, od int) ([]*model.VoucherRecord, error) {
	var voucherRecordSlice []*model.VoucherRecord
	//strSql, values := transferListSql(voucherRecordTN, filter, voucherRecordFields, limit, offset, order, od)
	strSql, values := transferListSqlWithMutiCondition(voucherRecordTN, filter, intervalFilter, fuzzyMatchFilter,
		voucherRecordFields, limit, offset, order, od)
	dao.Logger.DebugContext(ctx, "[VoucherRecord/db/List] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[VoucherRecord/db/List] [SqlElapsed: %v]", time.Since(start))
	}()
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

func (dao *VoucherRecordDao) UpdateByRecordId(ctx context.Context, do DbOperator, recordId int,
	params map[string]interface{}) error {
	strSql := "update " + voucherRecordTN + " set "
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
	strSql += " where record_id = ?"
	values = append(values, recordId)
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

func (dao *VoucherRecordDao) UpdateByVoucherId(ctx context.Context, do DbOperator, voucherId int,
	params map[string]interface{}) error {
	strSql := "update " + voucherRecordTN + " set "
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
	strSql += " where voucher_id = ?"
	values = append(values, voucherId)
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
