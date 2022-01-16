package db

import (
	"context"
	"database/sql"
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
	yearBalanceFields   = []string{"subject_id", "summary", "subjectDirection", "balance"}
	scanYearBalanceTask = func(r DbScanner, st *model.YearBalanceView) error {
		return r.Scan(&st.SubjectID, &st.Summary, &st.SubjectDirection, &st.Balance)
	}
)

func (dao *YearBalanceDao) GetByID(ctx context.Context, do DbOperator, subjectID int) (*model.YearBalanceView, error) {
	strSql := "select " + strings.Join(yearBalanceFields, ",") + " from " + yearBalanceTN + " where subject_id=?"
	dao.Logger.DebugContext(ctx, "[yearBalance/db/GetByID] [sql: %s ,values: %d]", strSql, subjectID)
	var yearBalance = &model.YearBalanceView{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[yearBalance/db/GetByID] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanYearBalanceTask(do.QueryRowContext(ctx, strSql, subjectID), yearBalance); err {
	case nil:
		return yearBalance, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/GetByID] [scanYearBalanceTask: %s]", err.Error())
		return nil, err
	}
}

func (dao *YearBalanceDao) Create(ctx context.Context, do DbOperator, st *model.YearBalanceParams) error {
	strSql := "insert into " + yearBalanceTN +
		" (" + strings.Join(yearBalanceFields, ",") + ") values (?, ?, ?, ?)"
	values := []interface{}{st.SubjectID, st.Summary, st.SubjectDirection, st.Balance}
	dao.Logger.DebugContext(ctx, "[yearBalance/db/Create] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[yearBalance/db/Create] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/Create] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *YearBalanceDao) DeleteByID(ctx context.Context, do DbOperator, subjectID int) error {
	strSql := "delete from " + yearBalanceTN + " where subject_id = ?"

	dao.Logger.DebugContext(ctx, "[yearBalance/db/DeleteByID] [sql: %s, id: %d]", strSql, subjectID)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[yearBalance/db/DeleteByID] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, subjectID); err != nil {
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/DeleteByID] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *YearBalanceDao) UpdateBySubID(ctx context.Context, do DbOperator, subjectID int,
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
	strSql += " where subject_id = ?"
	values = append(values, subjectID)
	start := time.Now()
	dao.Logger.DebugContext(ctx, "[yearBalance/db/UpdateBySubID] [sql: %s, values: %v]", strSql, values)
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[yearBalance/db/UpdateBySubID] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[yearBalance/db/UpdateBySubID] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}
