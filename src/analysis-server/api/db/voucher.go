package db

import (
	"context"
	"database/sql"
	"time"

	"analysis-server/model"
	"common/log"
)

//process voucher;use voucherRecordInfo and voucherInfo
type VoucherDao struct {
	Logger log.ILog
}

//计算某个科目在某段时间内的累计的贷方和借方金额。该函数用于银行明细账中的累计部分。
func (dao *VoucherDao) CalcAccuMoney(ctx context.Context, do DbOperator,
	params *model.CalAccuMoneyParams) (*model.AccuMoneyValueView, error) {
	voucherInfoTable := GenTableName(*params.VoucherYear, voucherInfoTN)
	voucherRecordTable := GenTableName(*params.VoucherYear, voucherRecordTN)
	var strSql string
	var values []interface{}
	if *params.Status == 0 {
		strSql = "select COALESCE(abs(sum(debit_money)), 0), COALESCE(abs(sum(credit_money)), 0) from " + voucherRecordTable +
			" where  sub_id1 = ? and voucher_id in (select voucher_id from " + voucherInfoTable +
			" where company_id = ? and voucher_month = ?)"
		values = []interface{}{*params.SubjectID, *params.CompanyID, *params.VoucherMonth}
	} else {
		strSql = "select COALESCE(abs(sum(debit_money)), 0), COALESCE(abs(sum(credit_money)), 0) from " + voucherRecordTable +
			"where  sub_id1 = ? and voucher_id in (select voucher_id from " + voucherInfoTable +
			" where company_id = ? and voucher_month = ? and status = ?)"
		values = []interface{}{*params.SubjectID, *params.CompanyID, *params.VoucherMonth, *params.Status}
	}

	dao.Logger.DebugContext(ctx, "[Voucher/db/CalcAccumulateAccSubSum] [sql: %s ,values: %d]", strSql, values)
	var accValue = &model.AccuMoneyValueView{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[VoucherRecord/db/Get] [SqlElapsed: %v]", time.Since(start))
	}()
	row := do.QueryRowContext(ctx, strSql, values...)
	switch err := row.Scan(&accValue.AccuDebitMoney, &accValue.AccuCreditMoney); err {
	case nil:
		return accValue, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[Voucher/db/Get] [Scan: %s]", err.Error())
		return nil, err
	}
}
