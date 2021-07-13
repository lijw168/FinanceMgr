package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"analysis-server/model"
	"common/log"
)

type LoginInfoDao struct {
	Logger log.ILog
}

var (
	loginInfoTN     = "userLoginInfo"
	loginInfoFields = []string{"operator_id", "name", "status", "client_ip", "begined_at", "ended_at"}
	scanloginInfo   = func(r DbScanner, st *model.LoginInfo) error {
		return r.Scan(&st.OperatorID, &st.Name, &st.Status, &st.ClientIp, &st.BeginedAt, &st.EndedAt)
	}
)

func (dao *LoginInfoDao) Get(ctx context.Context, do DbOperator, optID int) (*model.LoginInfo, error) {
	strSql := "select " + strings.Join(loginInfoFields, ",") + " from " +
		loginInfoTN + " where operator_id=?"
	dao.Logger.DebugContext(ctx, "[LoginInfo/db/Get] [sql: %s ,values: %d]", strSql, optID)
	var optInfo = &model.LoginInfo{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[LoginInfo/db/Get] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanloginInfo(do.QueryRowContext(ctx, strSql, optID), optInfo); err {
	case nil:
		return optInfo, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[LoginInfo/db/Get] [scanloginInfo: %s]", err.Error())
		return nil, err
	}
}

//list count by filter
func (dao *LoginInfoDao) CountByFilter(ctx context.Context, do DbOperator, filter map[string]interface{}) (int64, error) {
	var c int64
	strSql, values := transferCountSql(loginInfoTN, filter)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, values...).Scan(&c)
	dao.Logger.InfoContext(ctx, "[LoginInfo/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

func (dao *LoginInfoDao) Create(ctx context.Context, do DbOperator, st *model.LoginInfo) error {
	strSql := "insert into " + loginInfoTN + " (" + strings.Join(loginInfoFields, ",") +
		") values (?, ?, ?, ?, ?, ?)"
	values := []interface{}{st.OperatorID, st.Name, st.Status, st.ClientIp, st.BeginedAt, st.EndedAt}
	dao.Logger.DebugContext(ctx, "[LoginInfo/db/Create] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[LoginInfo/db/Create] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[LoginInfo/db/Create] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *LoginInfoDao) Delete(ctx context.Context, do DbOperator, optID int) error {
	strSql := "delete from " + loginInfoTN + " where operator_id = ?"

	dao.Logger.DebugContext(ctx, "[LoginInfo/db/Delete] [sql: %s, id: %s]", strSql, optID)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[LoginInfo/db/Delete] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, optID); err != nil {
		dao.Logger.ErrorContext(ctx, "[LoginInfo/db/Delete] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *LoginInfoDao) List(ctx context.Context, do DbOperator, filter map[string]interface{},
	limit int, offset int, order string, od int) ([]*model.LoginInfo, error) {
	var LoginInfoSlice []*model.LoginInfo
	strSql, values := transferListSql(loginInfoTN, filter, loginInfoFields, limit, offset, order, od)
	dao.Logger.DebugContext(ctx, "[LoginInfo/db/List] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[LoginInfo/db/List] [SqlElapsed: %v]", time.Since(start))
	}()
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[LoginInfo/db/List] [do.Query: %s]", err.Error())
		return LoginInfoSlice, err
	}
	defer result.Close()
	for result.Next() {
		LoginInfo := new(model.LoginInfo)
		err = scanloginInfo(result, LoginInfo)
		if err != nil {
			dao.Logger.ErrorContext(ctx, "[LoginInfo/db/List] [ScanSnapshot: %s]", err.Error())
			return LoginInfoSlice, err
		}
		LoginInfoSlice = append(LoginInfoSlice, LoginInfo)
	}
	return LoginInfoSlice, nil
}

//由于该表不能靠表中的一个字段，就可以确定一条记录，需要多个字段来确定一条记录，所以才有如下的实现方式。
func (dao *LoginInfoDao) Update(ctx context.Context, do DbOperator, filter map[string]interface{},
	updateField map[string]interface{}) error {
	strSql, values := transferUpdateSql(loginInfoTN, filter, updateField)
	dao.Logger.DebugContext(ctx, "[LoginInfo/db/Update] sql %s with values %v", strSql, values)
	start := time.Now()
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[LoginInfo/db/Update] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[LoginInfo/db/Update] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}
