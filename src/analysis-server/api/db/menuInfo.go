package db

import (
	"context"
	"database/sql"
	"financeMgr/src/common/log"
	"strings"
	"time"

	"financeMgr/src/analysis-server/model"
)

type MenuInfoDao struct {
	// Logger *log.Logger
	Logger log.ILog
}

var (
	menuInfoTN       = "menuInfo"
	menuInfoFields   = []string{"menu_id", "menu_name", "menu_level", "parent_menu_id", "menu_serial_num"}
	scanMenuInfoTask = func(r DbScanner, st *model.MenuInfo) error {
		return r.Scan(&st.MenuID, &st.MenuName, &st.MenuLevel, &st.ParentMenuID, &st.MenuSerialNum)
	}
)

func (dao *MenuInfoDao) GetMenuInfoByID(ctx context.Context, do DbOperator, menuID int) (*model.MenuInfo, error) {
	strSql := "select " + strings.Join(menuInfoFields, ",") + " from " + menuInfoTN + " where menu_id=?"
	dao.Logger.DebugContext(ctx, "[menuInfo/db/GetMenuInfoByID] [sql: %s ,values: %d]", strSql, menuID)
	var menuInfo = &model.MenuInfo{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[menuInfo/db/GetMenuInfoByID] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanMenuInfoTask(do.QueryRowContext(ctx, strSql, menuID), menuInfo); err {
	case nil:
		return menuInfo, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[menuInfo/db/GetMenuInfoByID] [scanMenuInfoTask: %s]", err.Error())
		return nil, err
	}
}

//list count by filter
func (dao *MenuInfoDao) CountByFilter(ctx context.Context, do DbOperator, filter map[string]interface{}) (int64, error) {
	var c int64
	strSql, values := transferCountSql(menuInfoTN, filter)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, values...).Scan(&c)
	dao.Logger.InfoContext(ctx, "[menuInfo/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

//list count by filter
func (dao *MenuInfoDao) Count(ctx context.Context, do DbOperator) (int64, error) {
	var c int64
	strSql := "select count(1) from " + menuInfoTN
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql).Scan(&c)
	dao.Logger.InfoContext(ctx, "[menuInfo/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

func (dao *MenuInfoDao) Create(ctx context.Context, do DbOperator, st *model.MenuInfo) error {
	strSql := "insert into " + menuInfoTN +
		" (" + strings.Join(menuInfoFields, ",") + ") values (?, ?, ?, ?, ?)"
	values := []interface{}{st.MenuID, st.MenuName, st.MenuLevel, st.ParentMenuID, st.MenuSerialNum}
	dao.Logger.DebugContext(ctx, "[menuInfo/db/Create] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[menuInfo/db/Create] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[menuInfo/db/Create] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *MenuInfoDao) DeleteByID(ctx context.Context, do DbOperator, menuID int) error {
	strSql := "delete from " + menuInfoTN + " where menu_id = ?"

	dao.Logger.DebugContext(ctx, "[menuInfo/db/DeleteByID] [sql: %s, id: %d]", strSql, menuID)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[menuInfo/db/DeleteByID] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, menuID); err != nil {
		dao.Logger.ErrorContext(ctx, "[menuInfo/db/DeleteByID] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *MenuInfoDao) List(ctx context.Context, do DbOperator, filter map[string]interface{}, limit int,
	offset int, order string, od int) ([]*model.MenuInfo, error) {
	var menuInfoSlice []*model.MenuInfo
	strSql, values := transferListSql(menuInfoTN, filter, menuInfoFields, limit, offset, order, od)
	dao.Logger.DebugContext(ctx, "[menuInfo/db/List] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[menuInfo/db/List] [SqlElapsed: %v]", time.Since(start))
	}()
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[menuInfo/db/List] [do.Query: %s]", err.Error())
		return menuInfoSlice, err
	}
	defer result.Close()
	for result.Next() {
		menuInfo := new(model.MenuInfo)
		err = scanMenuInfoTask(result, menuInfo)
		if err != nil {
			dao.Logger.ErrorContext(ctx, "[menuInfo/db/List] [ScanSnapshot: %s]", err.Error())
			return menuInfoSlice, err
		}
		menuInfoSlice = append(menuInfoSlice, menuInfo)
	}
	return menuInfoSlice, nil
}

func (dao *MenuInfoDao) UpdateBySubID(ctx context.Context, do DbOperator, menuID int,
	params map[string]interface{}) error {
	strSql := "update " + menuInfoTN + " set "
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
	strSql += " where menu_id = ?"
	values = append(values, menuID)
	start := time.Now()
	dao.Logger.DebugContext(ctx, "[menuInfo/db/UpdateBySubID] [sql: %s, values: %v]", strSql, values)
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[menuInfo/db/UpdateBySubID] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[menuInfo/db/UpdateBySubID] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}
