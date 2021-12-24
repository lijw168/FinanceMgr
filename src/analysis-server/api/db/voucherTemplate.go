package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"analysis-server/model"
	"common/log"
)

type VoucherTemplateDao struct {
	// Logger *log.Logger
	Logger log.ILog
}

var (
	voucherTemplateTN     = "voucherTemplate"
	voucherTemplateFields = []string{"voucher_template_id", "reference_voucher_id", "voucher_year", "illustration", "created_at"}
	scanVoucherTemplate   = func(r DbScanner, st *model.VoucherTemplate) error {
		return r.Scan(&st.VoucherTemplateID, &st.RefVoucherID, &st.VoucherYear, &st.Illustration, &st.CreatedAt)
	}
)

func (dao *VoucherTemplateDao) Get(ctx context.Context, do DbOperator, voucherTemplateID int) (*model.VoucherTemplate, error) {
	strSql := "select " + strings.Join(voucherTemplateFields, ",") + " from " + voucherTemplateTN + " where voucher_template_id=?"
	dao.Logger.DebugContext(ctx, "[VoucherTemplate/db/Get] [sql: %s ,values: %d]", strSql, voucherTemplateID)
	var comVoucher = &model.VoucherTemplate{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[/db/Get] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanVoucherTemplate(do.QueryRowContext(ctx, strSql, voucherTemplateID), comVoucher); err {
	case nil:
		return comVoucher, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[VoucherTemplate/db/Get] [scanVoucherTemplate: %s]", err.Error())
		return nil, err
	}
}

//get the count of the table
func (dao *VoucherTemplateDao) Count(ctx context.Context, do DbOperator) (int64, error) {
	var c int64
	strSql := "select count(1) from " + voucherTemplateTN
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql).Scan(&c)
	dao.Logger.InfoContext(ctx, "[VoucherTemplate/db/Count] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

//list count by filter ...
func (dao *VoucherTemplateDao) CountByFilter(ctx context.Context, do DbOperator,
	filter map[string]interface{}) (int64, error) {
	var c int64
	strSql, values := transferCountSql(voucherTemplateTN, filter)
	dao.Logger.DebugContext(ctx, "[VoucherTemplate/db/CountByFilter] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, values...).Scan(&c)
	dao.Logger.InfoContext(ctx, "[voucherInfo/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

func (dao *VoucherTemplateDao) Create(ctx context.Context, do DbOperator, st *model.VoucherTemplate) error {
	strSql := "insert into " + voucherTemplateTN + " (" + strings.Join(voucherTemplateFields, ",") +
		") values (?, ?, ?, ?, ?)"
	values := []interface{}{st.VoucherTemplateID, st.RefVoucherID, st.VoucherYear, st.Illustration, st.CreatedAt}
	dao.Logger.DebugContext(ctx, "[VoucherTemplate/db/Create] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[VoucherTemplate/db/Create] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherTemplate/db/Create] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}
func (dao *VoucherTemplateDao) Delete(ctx context.Context, do DbOperator, voucherTemplateID int) error {
	strSql := "delete from " + voucherTemplateTN + " where voucher_template_id=?"

	dao.Logger.DebugContext(ctx, "[VoucherTemplate/db/Delete] [sql: %s, id: %d]", strSql, voucherTemplateID)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[VoucherTemplate/db/Delete] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, voucherTemplateID); err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherTemplate/db/Delete] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

//没有复杂的匹配条件
func (dao *VoucherTemplateDao) SimpleList(ctx context.Context, do DbOperator, filter map[string]interface{},
	limit, offset, od int, order string) ([]*model.VoucherTemplate, error) {
	var comVoucherSlice []*model.VoucherTemplate
	strSql, values := transferListSql(voucherTemplateTN, filter, voucherTemplateFields, limit, offset, order, od)
	dao.Logger.DebugContext(ctx, "[VoucherTemplate/db/SimpleList] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[VoucherTemplate/db/SimpleList] [SqlElapsed: %v]", time.Since(start))
	}()
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[VoucherTemplate/db/SimpleList] [do.Query: %s]", err.Error())
		return comVoucherSlice, err
	}
	defer result.Close()
	for result.Next() {
		comVoucher := new(model.VoucherTemplate)
		err = scanVoucherTemplate(result, comVoucher)
		if err != nil {
			dao.Logger.ErrorContext(ctx, "[VoucherTemplate/db/SimpleList] [ScanSnapshot: %s]", err.Error())
			return comVoucherSlice, err
		}
		comVoucherSlice = append(comVoucherSlice, comVoucher)
	}
	return comVoucherSlice, nil
}
