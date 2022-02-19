package db

import (
	"context"
	//"database/sql"
	"strings"
	"time"

	"analysis-server/model"
	"common/log"
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
	err := do.QueryRowContext(ctx, strSql, subjectID).Scan(&dBalanceValue)
	return dBalanceValue, err
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

// var scanBalanceTask = func(r DbScanner, st *model.YearBalance) error {
// 	return r.Scan(&st.SubjectID, &st.Balance)
// }

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
