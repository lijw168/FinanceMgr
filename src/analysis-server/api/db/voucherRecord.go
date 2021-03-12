package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"common/log"
	"analysis-server/model"
)

type VoucherRecordDao struct {
	// Logger *log.Logger
	Logger log.ILog
}

var (
	voucherRecordTN     = "voucherRecordInfo_2020"
	voucherRecordFields = []string{"recordId", "voucherId", "subjectName", "debitMoney", "creditMoney", 
									"summary", "subId1", "subId2", "subId3", "subId4","billCount"}
	scanVoucherRecord   = func(r DbScanner, st *model.VoucherRecord) error {
		return r.Scan(&st.RecordID, &st.VoucherID, &st.SubjectName, &st.DebitMoney, &st.CreditMoney,
					 &st.Summary, &st.SbuID1,&st.SbuID2,&st.SbuID3,&st.SbuID4,&st.BillCount)
	}
)

func (dao *VoucherRecordDao) Get(ctx context.Context, do DbOperator, recordId int) (*model.VoucherRecord, error) {
	strSql := "select " + strings.Join(voucherRecordFields, ",") + " from " + voucherRecordTN + " where recordId=?"
	dao.Logger.DebugContext(ctx, "[VoucherRecord/db/Get] [sql: %s ,values: %s]", strSql, recordId)
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
func (dao *VoucherInfoDao) Count(ctx context.Context, do DbOperator) (int64, error) {
	var c int64
	strSql := "select count(1) from " + voucherInfoTN
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, nil).Scan(&c)
	dao.Logger.InfoContext(ctx, "[accountSubject/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

//list count by filter
// func (d *VoucherRecordDao) CountByFilter(ctx context.Context, do DbOperator, filter map[string]interface{}) (int64, error) {
// 	var c int64
// 	strSql, values := transferCountSql(voucherRecordTN, filter)
// 	start := time.Now()
// 	err := do.QueryRowContext(ctx, strSql, values...).Scan(&c)
// 	d.Logger.InfoContext(ctx, "[VoucherRecord/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
// 	return c, err
// }

func (dao *VoucherRecordDao) Create(ctx context.Context, do DbOperator, st *model.VoucherRecord) error {
	strSql := "insert into " + voucherRecordTN + " (" + strings.Join(voucherRecordFields, ",") + ") 
	           values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	values := []interface{}{st.RecordID, st.VoucherID, st.SubjectName, st.DebitMoney, st.CreditMoney,
		                    st.Summary, st.SbuID1,st.SbuID2,st.SbuID3,st.SbuID4,st.BillCount}
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
	strSql := "delete from " + voucherRecordTN + " where recordId = ?"

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
	strSql := "delete from " + voucherRecordTN + " where voucherId = ?"

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

func (dao *VoucherRecordDao) List(ctx context.Context, do DbOperator, filter map[string]interface{}, limit int,
	offset int, order string, od int) ([]*model.VoucherRecord, error) {
	var voucherRecordSlice []*model.VoucherRecord
	strSql, values := transferListSql(voucherRecordTN, filter, voucherRecordFields, limit, offset, order, od)
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

func (dao *VoucherRecordDao) Update(ctx context.Context, do DbOperator, strName string,
	                               params map[string]interface{}) error {
	var keyMap = map[string]string{"VoucherID": "voucherId","SubjectName": "subjectName", "DebitMoney": "debitMoney",
							"CreditMoney": "creditMoney", "Summary": "summary", "SubID1": "subId1", "SubID2": "subId2",
							"SubID3": "subId3", "RecordID": "recordId"}
	strSql := "update " + voucherRecordTN + " set "
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
	strSql += " where name = ?"
	values = append(values, strName)
	start := time.Now()
	dao.Logger.DebugContext(ctx, "[VoucherRecord/db/Update] [sql: %s, values: %v]", strSql, values)
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[VoucherRecord/db/Update] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherRecord/db/Update] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}
