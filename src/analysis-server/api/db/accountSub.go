package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"common/log"
	"analysis-server/model"
)

type AccSubDao struct {
	// Logger *log.Logger
	Logger log.ILog
}

var (
	accSubInfoTN     = "accountSubject"
	accSubInfoFields = []string{"subjectId", "subjectName", "subjectLevel"}
	scanAccSubTask   = func(r DbScanner, st *model.AccSubject) error {
		return r.Scan(&st.SubjectID, &st.SubjectName, &st.SubjectLevel)
	}
)

func (dao *AccSubDao) GetAccSubByName(ctx context.Context, do DbOperator, strName string) (*model.AccSubject, error) {
	strSql := "select " + strings.Join(accSubInfoFields, ",") + " from " + accSubInfoTN + " where subjectName=?"
	dao.Logger.DebugContext(ctx, "[accountSubject/db/GetAccSubByName] [sql: %s ,values: %s]", strSql, strName)
	var accSub = &model.AccSubject{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[accountSubject/db/GetAccSubByName] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanAccSubTask(do.QueryRowContext(ctx, strSql, strName), accSub); err {
	case nil:
		return accSub, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[accountSubject/db/GetAccSubByName] [scanAccSubTask: %s]", err.Error())
		return nil, err
	}
}

func (dao *AccSubDao) GetAccSubByID(ctx context.Context, do DbOperator, strID string) (*model.AccSubject, error) {
	strSql := "select " + strings.Join(accSubInfoFields, ",") + " from " + accSubInfoTN + " where subjectId=?"
	dao.Logger.DebugContext(ctx, "[accountSubject/db/GetAccSubByID] [sql: %s ,values: %s]", strSql, strID)
	var accSub = &model.AccSubject{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[accountSubject/db/GetAccSubByID] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanAccSubTask(do.QueryRowContext(ctx, strSql, strID), accSub); err {
	case nil:
		return accSub, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[accountSubject/db/GetAccSubByID] [scanAccSubTask: %s]", err.Error())
		return nil, err
	}
}

//list count by filter
func (dao *AccSubDao) CountByFilter(ctx context.Context, do DbOperator, filter map[string]interface{}) (int64, error) {
	var c int64
	strSql, values := transferCountSql(accSubInfoTN, filter)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, values...).Scan(&c)
	dao.Logger.InfoContext(ctx, "[accountSubject/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

//list count by filter
func (dao *AccSubDao) Count(ctx context.Context, do DbOperator) (int64, error) {
	var c int64
	strSql := "select count(1) from " + accSubInfoTN
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, nil).Scan(&c)
	dao.Logger.InfoContext(ctx, "[accountSubject/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

func (dao *AccSubDao) Create(ctx context.Context, do DbOperator, st *model.AccSubject) error {
	strSql := "insert into " + accSubInfoTN + " (" + strings.Join(accSubInfoFields, ",") + ") values (?, ?, ?)"
	values := []interface{}{st.SubjectID, st.SubjectName, st.SubjectLevel}
	dao.Logger.DebugContext(ctx, "[accountSubject/db/Create] [sql: %s, values: %v]", strSql, values)
	start := time.Now()
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[accountSubject/db/Create] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[accountSubject/db/Create] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *AccSubDao) DeleteByName(ctx context.Context, do DbOperator, strName string) error {
	strSql := "delete from " + accSubInfoTN + " where subjectName = ?"

	dao.Logger.DebugContext(ctx, "[accountSubject/db/DeleteByName] [sql: %s, id: %s]", strSql, strName)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[accountSubject/db/DeleteByName] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, strName); err != nil {
		dao.Logger.ErrorContext(ctx, "[accountSubject/db/DeleteByName] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *AccSubDao) DeleteByID(ctx context.Context, do DbOperator, subjectID int) error {
	strSql := "delete from " + accSubInfoTN + " where subjectId = ?"

	dao.Logger.DebugContext(ctx, "[accountSubject/db/DeleteByID] [sql: %s, id: %s]", strSql, subjectID)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[accountSubject/db/DeleteByID] [SqlElapsed: %v]", time.Since(start))
	}()
	if _, err := do.ExecContext(ctx, strSql, subjectID); err != nil {
		dao.Logger.ErrorContext(ctx, "[accountSubject/db/DeleteByID] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *AccSubDao) List(ctx context.Context, do DbOperator, filter map[string]interface{}, limit int,
						   offset int, order string, od int) ([]*model.AccSubject, error) {
	var accountSubjectSlice []*model.AccSubject
	strSql, values := transferListSql(accSubInfoTN, filter, accSubInfoFields, limit, offset, order, od)
	dao.Logger.DebugContext(ctx, "[accountSubject/db/List] sql %s with values %v", strSql, values)
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[accountSubject/db/List] [SqlElapsed: %v]", time.Since(start))
	}()
	result, err := do.QueryContext(ctx, strSql, values...)
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[accountSubject/db/List] [do.Query: %s]", err.Error())
		return accountSubjectSlice, err
	}
	defer result.Close()
	for result.Next() {
		accountSubject := new(model.AccSubject)
		err = scanAccSubTask(result, accountSubject)
		if err != nil {
			dao.Logger.ErrorContext(ctx, "[accountSubject/db/List] [ScanSnapshot: %s]", err.Error())
			return accountSubjectSlice, err
		}
		accountSubjectSlice = append(accountSubjectSlice, accountSubject)
	}
	return accountSubjectSlice, nil
}

// func (dao *AccSubDao) ListWithFilterNo(ctx context.Context, do DbOperator, filter map[string]interface{},
// 	filterNo map[string]interface{}, limit int, offset int, order string, od int) ([]*model.accountSubject, error) {
// 	var accountSubjectSlice []*model.accountSubject
// 	strSql, values := transferListSqlWithNo(accSubInfoTN, filter, filterNo, accSubInfoFields, limit, offset, order, od)
// 	dao.Logger.DebugContext(ctx, "[accountSubject/db/ListWithFilterNo] sql %s with values %v", strSql, values)
// 	start := time.Now()
// 	defer func() {
// 		dao.Logger.InfoContext(ctx, "[accountSubject/db/ListWithFilterNo] [SqlElapsed: %v]", time.Since(start))
// 	}()
// 	result, err := do.QueryContext(ctx, strSql, values...)
// 	if err != nil {
// 		dao.Logger.ErrorContext(ctx, "[accountSubject/db/ListWithFilterNo] [do.Query: %s]", err.Error())
// 		return accountSubjectSlice, err
// 	}
// 	defer result.Close()
// 	for result.Next() {
// 		accountSubject := new(model.accountSubject)
// 		err = scanAccSubTask(result, accountSubject)
// 		if err != nil {
// 			dao.Logger.ErrorContext(ctx, "[accountSubject/db/ListWithFilterNo] [ScanSnapshot: %s]", err.Error())
// 			return accountSubjectSlice, err
// 		}
// 		accountSubjectSlice = append(accountSubjectSlice, accountSubject)
// 	}
// 	return accountSubjectSlice, nil
// }

func (dao *AccSubDao) UpdateBySubID(ctx context.Context, do DbOperator, strSubID string,
	params map[string]interface{}) error {
	var keyMap = map[string]string{"SubjectID": "subjectId", "SubjectName": "subjectName", "SubjectLevel": "subjectLevel"}
	strSql := "update " + accSubInfoTN + " set "
	var values []interface{}
	var first bool = true
	for key, value := range params {
		if dbKey, ok := keyMap[key]; ok {
			if first {
				strSql += dbKey + "=?"
				first = false
			} else {
				strSql += "," + dbKey + "=?"
			}
			values = append(values, value)
		}
	}
	if first {
		return nil
	}
	strSql += " where subjectId = ?"
	values = append(values, strSubID)
	start := time.Now()
	dao.Logger.DebugContext(ctx, "[accountSubject/db/UpdateBySubID] [sql: %s, values: %v]", strSql, values)
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[accountSubject/db/UpdateBySubID] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[accountSubject/db/UpdateBySubID] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}

func (dao *AccSubDao) UpdateByName(ctx context.Context, do DbOperator, strSubName string,
	params map[string]interface{}) error {
	var keyMap = map[string]string{"SubjectID": "subjectId", "SubjectName": "subjectName", "SubjectLevel": "subjectLevel"}
	strSql := "update " + accSubInfoTN + " set "
	var values []interface{}
	var first bool = true
	for key, value := range params {
		if dbKey, ok := keyMap[key]; ok {
			if first {
				strSql += dbKey + "=?"
				first = false
			} else {
				strSql += "," + dbKey + "=?"
			}
			values = append(values, value)
		}
	}
	if first {
		return nil
	}
	strSql += " where subjectName = ?"
	values = append(values, strSubName)
	start := time.Now()
	dao.Logger.DebugContext(ctx, "[accountSubject/db/UpdateByName] [sql: %s, values: %v]", strSql, values)
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[accountSubject/db/UpdateByName] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[accountSubject/db/UpdateByName] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}
