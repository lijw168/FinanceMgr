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
	yearBalanceFields   = []string{"subject_id", "year", "balance"}
	scanYearBalanceTask = func(r DbScanner, st *model.YearBalance) error {
		return r.Scan(&st.SubjectID, &st.Year, &st.Balance)
	}
)

func (dao *YearBalanceDao) GetYearBalance(ctx context.Context, do DbOperator, iYear, subjectID int) (float64, error) {
	strSql := "select balance from " + yearBalanceTN + " where subject_id=? and year=?"
	dao.Logger.DebugContext(ctx, "[yearBalance/db/GetYearBalance] [sql: %s ,values: %d,%d]", strSql, subjectID, iYear)
	var dBalanceValue float64
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[accountSubject/db/GetYearBalance] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := do.QueryRowContext(ctx, strSql, subjectID, iYear).Scan(&dBalanceValue); err {
	case nil:
		return dBalanceValue, nil
	case sql.ErrNoRows:
		//根据业务的规则，如果是没有获取到相应的数据，则可以返回0值，0值在前端也表示没有相应的数据。
		return 0, nil
	default:
		dao.Logger.ErrorContext(ctx, "[accountSubject/db/GetAccSubByID] [scanAccSubTask: %s]", err.Error())
		return 0, err
	}
}

func (dao *YearBalanceDao) CreateYearBalance(ctx context.Context, do DbOperator, st *model.OptYearBalanceParams) error {
	strSql := "insert into " + yearBalanceTN +
		" (" + strings.Join(yearBalanceFields, ",") + ") values (?, ?, ?)"
	values := []interface{}{st.SubjectID, st.Year, st.Balance}
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

func (dao *YearBalanceDao) DeleteYearBalance(ctx context.Context, do DbOperator, iYear, subjectID int) error {
	strSql := "delete from " + yearBalanceTN + " where subject_id = ? and year=?"

	dao.Logger.DebugContext(ctx, "[yearBalance/db/DeleteYearBalance] [sql: %s, id:%d, year: %d]", strSql, subjectID, iYear)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[yearBalance/db/DeleteYearBalance] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, subjectID, iYear); err != nil {
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/DeleteYearBalance] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *YearBalanceDao) DeleteYearBalanceOneYear(ctx context.Context, do DbOperator, iYear int) error {
	strSql := "delete from " + yearBalanceTN + " where subject_id = ?"

	dao.Logger.DebugContext(ctx, "[yearBalance/db/DeleteYearBalance] [sql: %s, year: %d]", strSql, iYear)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[yearBalance/db/DeleteYearBalance] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, iYear); err != nil {
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/DeleteYearBalance] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

// 可以更新yearBalanceTN中的所有字段，为以后增加字段保留的接口
func (dao *YearBalanceDao) UpdateYearBalance(ctx context.Context, do DbOperator, iYear, subjectID int,
	params map[string]interface{}) error {
	strSql := "update " + yearBalanceTN + " set "
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
	strSql += " where subject_id = ? and year=?"
	values = append(values, subjectID, iYear)
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
func (dao *YearBalanceDao) UpdateBalance(ctx context.Context, do DbOperator, st *model.OptYearBalanceParams) error {
	strSql := "update " + yearBalanceTN + " set balance = ? where subject_id = ? and year=?"
	values := []interface{}{st.Balance, st.SubjectID, st.Year}
	start := time.Now()
	dao.Logger.DebugContext(ctx, "[yearBalance/db/UpdateBalance] [sql: %s, values: %v]", strSql, values)
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[yearBalance/db/UpdateBalance] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/UpdateBalance] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *YearBalanceDao) ListYearBalance(ctx context.Context, do DbOperator, filter map[string]interface{},
	limit int, offset int, order string, od int) ([]*model.YearBalance, error) {
	var yearBalSlice []*model.YearBalance
	//resFields := []string{"subject_id", "balance"}
	//strSql, values := transferListSqlWithNo(accSubInfoTN, filter, filterNo, resFields, limit, offset, order, od)
	strSql, values := transferListSql(yearBalanceTN, filter, yearBalanceFields, limit, offset, order, od)
	dao.Logger.DebugContext(ctx, "[accountSubject/db/ListYearBalance] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[accountSubject/db/ListYearBalance] [SqlElapsed: %v]", time.Since(start))
	}()
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[accountSubject/db/ListYearBalance] [do.Query: %s]", err.Error())
		return yearBalSlice, err
	}
	defer result.Close()
	for result.Next() {
		yearBal := new(model.YearBalance)
		err = scanYearBalanceTask(result, yearBal)
		if err != nil {
			dao.Logger.ErrorContext(ctx, "[accountSubject/db/List] [ScanSnapshot: %s]", err.Error())
			return yearBalSlice, err
		}
		yearBalSlice = append(yearBalSlice, yearBal)
	}
	return yearBalSlice, nil
}
