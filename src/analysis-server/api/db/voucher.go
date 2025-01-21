package db

import (
	"context"
	"database/sql"
	"financeMgr/src/common/log"
	"time"

	"financeMgr/src/analysis-server/model"
)

// process voucher;use voucherRecordInfo and voucherInfo
type VoucherDao struct {
	Logger log.ILog
}

// 计算某个科目截止到某个月份的累计的贷方和借方金额。该函数用于银行明细账中的累计部分。
func (dao *VoucherDao) CalcAccuMoney(ctx context.Context, do DbOperator,
	params *model.CalAccuMoney) (*model.AccuMoneyValueView, error) {
	voucherInfoTable := GenTableName(params.VoucherYear, voucherInfoTN)
	voucherRecordTable := GenTableName(params.VoucherYear, voucherRecordTN)
	var strSql string
	var values []interface{}
	//该判断应该是兼容之前，没有设置凭证状态的数据。
	if params.Status == 0 {
		strSql = "select COALESCE(abs(sum(debit_money)), 0), COALESCE(abs(sum(credit_money)), 0) from " + voucherRecordTable +
			" where  sub_id1 = ? and voucher_id in (select voucher_id from " + voucherInfoTable +
			" where company_id = ? and voucher_month <= ?)"
		values = []interface{}{params.SubjectID, params.CompanyID, params.VoucherMonth}
	} else {
		strSql = "select COALESCE(abs(sum(debit_money)), 0), COALESCE(abs(sum(credit_money)), 0) from " + voucherRecordTable +
			" where  sub_id1 = ? and voucher_id in (select voucher_id from " + voucherInfoTable +
			" where company_id = ? and voucher_month <= ? and status = ?)"
		values = []interface{}{params.SubjectID, params.CompanyID, params.VoucherMonth, params.Status}
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

var scanPartialVouRecords = func(r DbScanner, st *model.AccountOfPeriod) error {
	return r.Scan(&st.SubjectID, &st.PeriodDebitMoney, &st.PeriodCreditMoney)
}

func handleArrFilter(arr []interface{}, s *string) (values []interface{}) {
	for i, ki := range arr {
		if i == 0 {
			*s += "?"
		} else {
			*s += ", ?"
		}
		values = append(values, ki)
	}
	return
}

// 该函数用于计算余额表中本期发生额。
func (dao *VoucherDao) GetPartialVouRecords(ctx context.Context, do DbOperator,
	params *model.CalAmountOfPeriodParams) ([]*model.AccountOfPeriod, error) {
	voucherInfoTable := GenTableName(*params.VoucherYear, voucherInfoTN)
	voucherRecordTable := GenTableName(*params.VoucherYear, voucherRecordTN)

	var values []interface{}
	tmpK := "sub_id1"
	tmpK += " in ("
	arr := []interface{}{}
	for _, ki := range params.SubjectIDArr {
		arr = append(arr, ki)
	}
	tmpFv := handleArrFilter(arr, &tmpK)
	tmpK += ")"
	values = append(values, tmpFv...)

	var strSql string
	if *params.Status == 0 {
		strSql = "select sub_id1, debit_money, credit_money from " + voucherRecordTable +
			" where  " + tmpK + " and voucher_id in (select voucher_id from " + voucherInfoTable +
			" where company_id = ? and voucher_month between ? and ?)"
		values = append(values, *params.CompanyID, *params.StartMonth, *params.EndMonth)
	} else {
		strSql = "select sub_id1, debit_money, credit_money from " + voucherRecordTable +
			" where  " + tmpK + " and voucher_id in (select voucher_id from " + voucherInfoTable +
			" where company_id = ? and status = ? and voucher_month between ? and ?)"
		values = append(values, *params.CompanyID, *params.Status, *params.StartMonth, *params.EndMonth)
	}
	accuMoneyValSlice := []*model.AccountOfPeriod{}
	dao.Logger.DebugContext(ctx, "[Voucher/db/GetPartialVouRecords] [sql: %s ,values: %d]", strSql, values)
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[Voucher/db/GetPartialVouRecords] [do.Query: %s]", err.Error())
		return accuMoneyValSlice, err
	}
	defer result.Close()
	for result.Next() {
		moneyVal := new(model.AccountOfPeriod)
		err = scanPartialVouRecords(result, moneyVal)
		if err != nil {
			dao.Logger.ErrorContext(ctx, "[Voucher/db/GetPartialVouRecords] [ScanSnapshot: %s]", err.Error())
			return accuMoneyValSlice, err
		}
		accuMoneyValSlice = append(accuMoneyValSlice, moneyVal)
	}
	return accuMoneyValSlice, nil
}
