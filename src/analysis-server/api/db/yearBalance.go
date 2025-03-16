package db

import (
	"context"
	"database/sql"
	"financeMgr/src/common/log"
	"strings"
	"time"

	"financeMgr/src/analysis-server/model"
)

type YearBalanceDao struct {
	Logger log.ILog
}

var (
	yearBalanceTN       = "beginOfYearBalance"
	yearBalanceFields   = []string{"company_id", "subject_id", "year", "balance", "status", "created_at", "updated_at"}
	scanYearBalanceTask = func(r DbScanner, st *model.YearBalance) error {
		return r.Scan(&st.CompanyID, &st.SubjectID, &st.Year, &st.Balance, &st.Status, &st.CreatedAt, &st.UpdatedAt)
	}
)

func (dao *YearBalanceDao) GetYearBalance(ctx context.Context, do DbOperator,
	params *model.BasicYearBalanceParams) (*model.YearBalance, error) {
	strSql := "select " + strings.Join(yearBalanceFields, ",") + " from " + yearBalanceTN + " where company_id = ? and year=? and subject_id=?"
	dao.Logger.DebugContext(ctx, "[yearBalance/db/GetYearBalance] [sql: %s ,comId:%d,year:%d,subId:%d]",
		strSql, *params.CompanyID, *params.Year, *params.SubjectID)
	yearBal := &model.YearBalance{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[yearBalance/db/GetYearBalance] [SqlElapsed: %v]", time.Since(start))
	}()
	err := scanYearBalanceTask(do.QueryRowContext(ctx, strSql, *params.CompanyID, *params.Year, *params.SubjectID), yearBal)
	switch err {
	case nil:
		return yearBal, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/GetYearBalance] [Scan: %s]", err.Error())
		return nil, err
	}
}

// get a account subject's balance value
func (dao *YearBalanceDao) GetAccSubYearBalValue(ctx context.Context, do DbOperator,
	params *model.BasicYearBalanceParams) (float64, error) {
	strSql := "select balance from " + yearBalanceTN + " where company_id = ? and year=? and subject_id=?"
	dao.Logger.DebugContext(ctx, "[yearBalance/db/GetAccSubYearBalValue] [sql: %s ,comId:%d,year:%d,subId:%d]",
		strSql, *params.CompanyID, *params.Year, *params.SubjectID)
	var dBalanceValue float64
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[yearBalance/db/GetAccSubYearBalValue] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := do.QueryRowContext(ctx, strSql, *params.CompanyID, *params.Year, *params.SubjectID).Scan(&dBalanceValue); err {
	case nil:
		return dBalanceValue, nil
	//因为目前financeHelper客户端接收不到所返回的错误信息及错误码，其错误码在client api库中被替换掉了。所以该错误不能返回到client api库。
	//等以后修改了改设计后，再返回相应的错误。
	case sql.ErrNoRows:
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/GetAccSubYearBalValue] [scanAccSubTask: %s]", err.Error())
		return 0, nil
	// case sql.ErrNoRows:
	// 	return 0, err
	default:
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/GetAccSubYearBalValue] [scanAccSubTask: %s]", err.Error())
		return 0, err
	}
}

func (dao *YearBalanceDao) CreateYearBalance(ctx context.Context, do DbOperator, st *model.YearBalance) error {
	strSql := "insert into " + yearBalanceTN +
		" (" + strings.Join(yearBalanceFields, ",") + ") values (?, ?, ?, ?, ?, ?, ?)"
	values := []interface{}{st.CompanyID, st.SubjectID, st.Year, st.Balance, st.Status, st.CreatedAt, st.UpdatedAt}
	dao.Logger.DebugContext(ctx, "[yearBalance/db/CreateYearBalance] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[yearBalance/db/CreateYearBalance] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/CreateYearBalance] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *YearBalanceDao) DeleteYearBalance(ctx context.Context, do DbOperator, params *model.BasicYearBalanceParams) error {
	strSql := "delete from " + yearBalanceTN + " where company_id = ? and year = ? and subject_id = ?"

	dao.Logger.DebugContext(ctx, "[yearBalance/db/DeleteYearBalance] [sql: %s, comId:%d,year:%d,subId:%d]",
		strSql, *params.CompanyID, *params.Year, *params.SubjectID)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[yearBalance/db/DeleteYearBalance] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, *params.CompanyID, *params.Year, *params.SubjectID); err != nil {
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/DeleteYearBalance] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *YearBalanceDao) BatchDeleteYearBalance(ctx context.Context, do DbOperator, filter map[string]interface{}) error {
	strSql, values := transferDeleteSql(yearBalanceTN, filter)
	dao.Logger.DebugContext(ctx, "[yearBalance/db/BatchDeleteYearBalance] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[yearBalance/db/BatchDeleteYearBalance] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, values...); err != nil {
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/BatchDeleteYearBalance] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

// 可以更新yearBalanceTN中的所有字段，所以可以把UpdateBalance这个函数去掉。
func (dao *YearBalanceDao) UpdateYearBalance(ctx context.Context, do DbOperator,
	filter map[string]interface{}, updateField map[string]interface{}) error {
	//因为需要跟新一下updated_at字段,并且在service中，有多处调用该函数，所以在此添加更新该字段比较合适，所以把这个字段添加到updateField
	updateField["updated_at"] = time.Now()
	strSql, values := makeUpdateSqlWithMultiCondition(yearBalanceTN, updateField, nil, filter, nil, nil)
	start := time.Now()
	dao.Logger.DebugContext(ctx, "[yearBalance/db/UpdateYearBalance] [sql: %s, values: %v]", strSql, values)
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[yearBalance/db/UpdateYearBalance] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/UpdateYearBalance] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

// 仅更新年初余额这一个字段
// func (dao *YearBalanceDao) UpdateBalance(ctx context.Context, do DbOperator, st *model.OptYearBalanceParams) error {
// 	strSql := "update " + yearBalanceTN + " set balance = ? , updated_at=? where company_id = ? and year=? and subject_id = ?"
// 	values := []interface{}{st.Balance, time.Now(), st.CompanyID, st.Year, st.SubjectID}
// 	start := time.Now()
// 	dao.Logger.DebugContext(ctx, "[yearBalance/db/UpdateBalance] [sql: %s, values: %v]", strSql, values)
// 	_, err := do.ExecContext(ctx, strSql, values...)
// 	dao.Logger.InfoContext(ctx, "[yearBalance/db/UpdateBalance] [SqlElapsed: %v]", time.Since(start))
// 	if err != nil {
// 		dao.Logger.ErrorContext(ctx, "[yearBalance/db/UpdateBalance] [do.Exec: %s]", err.Error())
// 		return err
// 	}
// 	return nil
// }

func (dao *YearBalanceDao) ListYearBalance(ctx context.Context, do DbOperator, filter map[string]interface{},
	limit int, offset int, order string, od int) ([]*model.YearBalance, error) {
	var yearBalSlice []*model.YearBalance
	//resFields := []string{"subject_id", "balance"}
	//strSql, values := transferListSqlWithNo(accSubInfoTN, filter, filterNo, resFields, limit, offset, order, od)
	strSql, values := transferListSql(yearBalanceTN, filter, yearBalanceFields, limit, offset, order, od)
	dao.Logger.DebugContext(ctx, "[yearBalance/db/ListYearBalance] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[yearBalance/db/ListYearBalance] [SqlElapsed: %v]", time.Since(start))
	}()
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/ListYearBalance] [do.Query: %s]", err.Error())
		return yearBalSlice, err
	}
	defer result.Close()
	for result.Next() {
		yearBal := new(model.YearBalance)
		err = scanYearBalanceTask(result, yearBal)
		if err != nil {
			dao.Logger.ErrorContext(ctx, "[yearBalance/db/List] [ScanSnapshot: %s]", err.Error())
			return yearBalSlice, err
		}
		yearBalSlice = append(yearBalSlice, yearBal)
	}
	return yearBalSlice, nil
}

// list count by filter
func (dao *YearBalanceDao) CountByFilter(ctx context.Context, do DbOperator, filter map[string]interface{}) (int, error) {
	var c int
	strSql, values := transferCountSql(yearBalanceTN, filter)
	dao.Logger.DebugContext(ctx, "[yearBalance/db/CountByFilter] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[yearBalance/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := do.QueryRowContext(ctx, strSql, values...).Scan(&c); err {
	case nil:
		return c, nil
	//因为目前financeHelper客户端接收不到所返回的错误信息及错误码，其错误码在client api库中被替换掉了。所以该错误不能返回到client api库。
	//等以后修改了改设计后，再返回相应的错误。
	case sql.ErrNoRows:
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/CountByFilter] [Scan: %s]", err.Error())
		return 0, nil
	// case sql.ErrNoRows:
	// 	return 0, err
	default:
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/CountByFilter] [Scan: %s]", err.Error())
		return 0, err
	}
}
