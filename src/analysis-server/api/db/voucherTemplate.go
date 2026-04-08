package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"financeMgr/src/analysis-server/model"
)

type VoucherTemplateDao struct {
	// Logger *log.Logger
	//Logger log.ILog
}

var (
	voucherTemplateTN     = "voucherTemplate"
	voucherTemplateFields = []string{"voucher_template_id", "company_id", "reference_voucher_id", "voucher_year", "illustration", "created_at"}
	scanVoucherTemplate   = func(r DbScanner, st *model.VoucherTemplate) error {
		return r.Scan(&st.VoucherTemplateID, &st.CompanyID, &st.RefVoucherID, &st.VoucherYear, &st.Illustration, &st.CreatedAt)
	}
)

func (dao *VoucherTemplateDao) Get(ctx context.Context, do DbOperator, voucherTemplateID int) (*model.VoucherTemplate, error) {
	strSql := "select " + strings.Join(voucherTemplateFields, ",") + " from " + voucherTemplateTN + " where voucher_template_id=?"
	gLogger.DebugContext(ctx, "[VoucherTemplate/db/Get] [sql: %s ,values: %d]", strSql, voucherTemplateID)
	var comVoucher = &model.VoucherTemplate{}
	start := time.Now()
	defer func() {
		gLogger.InfoContext(ctx, "[/db/Get] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanVoucherTemplate(do.QueryRowContext(ctx, strSql, voucherTemplateID), comVoucher); err {
	case nil:
		return comVoucher, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		gLogger.ErrorContext(ctx, "[VoucherTemplate/db/Get] [scanVoucherTemplate: %s]", err.Error())
		return nil, err
	}
}

// get the count of the table
func (dao *VoucherTemplateDao) Count(ctx context.Context, do DbOperator) (int64, error) {
	var c int64
	strSql := "select count(1) from " + voucherTemplateTN
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql).Scan(&c)
	gLogger.InfoContext(ctx, "[VoucherTemplate/db/Count] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

// list count by filter ...
func (dao *VoucherTemplateDao) CountByFilter(ctx context.Context, do DbOperator,
	filter map[string]interface{}) (int64, error) {
	var c int64
	strSql, values := transferCountSql(voucherTemplateTN, filter)
	gLogger.DebugContext(ctx, "[VoucherTemplate/db/CountByFilter] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, values...).Scan(&c)
	gLogger.InfoContext(ctx, "[voucherInfo/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

func (dao *VoucherTemplateDao) Create(ctx context.Context, do DbOperator, st *model.VoucherTemplate) error {
	strSql := "insert into " + voucherTemplateTN + " (" + strings.Join(voucherTemplateFields, ",") +
		") values (?, ?, ?, ?, ?, ?)"
	values := []interface{}{st.VoucherTemplateID, st.CompanyID, st.RefVoucherID, st.VoucherYear, st.Illustration, st.CreatedAt}
	gLogger.DebugContext(ctx, "[VoucherTemplate/db/Create] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	_, err := do.ExecContext(ctx, strSql, values...)
	gLogger.InfoContext(ctx, "[VoucherTemplate/db/Create] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		gLogger.ErrorContext(ctx, "[VoucherTemplate/db/Create] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}
func (dao *VoucherTemplateDao) Delete(ctx context.Context, do DbOperator, voucherTemplateID int) error {
	strSql := "delete from " + voucherTemplateTN + " where voucher_template_id=?"

	gLogger.DebugContext(ctx, "[VoucherTemplate/db/Delete] [sql: %s, id: %d]", strSql, voucherTemplateID)
	start := time.Now()
	defer func() {
		gLogger.InfoContext(ctx, "[VoucherTemplate/db/Delete] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, voucherTemplateID); err != nil {
		gLogger.ErrorContext(ctx, "[VoucherTemplate/db/Delete] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

// 没有复杂的匹配条件
func (dao *VoucherTemplateDao) SimpleList(ctx context.Context, do DbOperator, filter map[string]interface{},
	limit, offset, od int, order string) ([]*model.VoucherTemplate, error) {
	var comVoucherSlice []*model.VoucherTemplate
	strSql, values := transferListSql(voucherTemplateTN, filter, voucherTemplateFields, limit, offset, order, od)
	gLogger.DebugContext(ctx, "[VoucherTemplate/db/SimpleList] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		gLogger.InfoContext(ctx, "[VoucherTemplate/db/SimpleList] [SqlElapsed: %v]", time.Since(start))
	}()
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		gLogger.ErrorContext(ctx, "[VoucherTemplate/db/SimpleList] [do.Query: %s]", err.Error())
		return comVoucherSlice, err
	}
	defer result.Close()
	for result.Next() {
		comVoucher := new(model.VoucherTemplate)
		err = scanVoucherTemplate(result, comVoucher)
		if err != nil {
			gLogger.ErrorContext(ctx, "[VoucherTemplate/db/SimpleList] [ScanSnapshot: %s]", err.Error())
			return comVoucherSlice, err
		}
		comVoucherSlice = append(comVoucherSlice, comVoucher)
	}
	return comVoucherSlice, nil
}

// get max voucherTemplateId
func (dao *VoucherTemplateDao) GetMaxVoucherTemplateId(ctx context.Context, do DbOperator) (int, error) {
	//use NullInt to avoid the error when the table is empty
	var maxID sql.NullInt64
	strSql := "select max(voucher_template_id) from " + voucherTemplateTN
	gLogger.DebugContext(ctx, "[VoucherTemplate/db/GetMaxVoucherTemplateId] [sql: %s]", strSql)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql).Scan(&maxID)
	gLogger.InfoContext(ctx, "[VoucherTemplate/db/GetMaxVoucherTemplateId] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		gLogger.ErrorContext(ctx, "[VoucherTemplate/db/GetMaxVoucherTemplateId] [do.QueryRowContext: %s]", err.Error())
		return 0, err
	}
	if maxID.Valid {
		return int(maxID.Int64), nil
	}
	return 0, nil
}
