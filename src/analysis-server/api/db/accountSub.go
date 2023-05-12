package db

import (
	"context"
	"database/sql"
	"financeMgr/src/common/log"
	"strings"
	"time"

	"financeMgr/src/analysis-server/model"
)

type AccSubDao struct {
	// Logger *log.Logger
	Logger log.ILog
}

var (
	accSubInfoTN     = "accountSubject"
	accSubInfoFields = []string{"subject_id", "company_id", "common_id", "subject_name", "subject_level",
		"subject_direction", "subject_type", "mnemonic_code", "subject_style"}
	scanAccSubTask = func(r DbScanner, st *model.AccSubject) error {
		return r.Scan(&st.SubjectID, &st.CompanyID, &st.CommonID, &st.SubjectName, &st.SubjectLevel,
			&st.SubjectDirection, &st.SubjectType, &st.MnemonicCode, &st.SubjectStyle)
	}
)

func (dao *AccSubDao) GetAccSubByID(ctx context.Context, do DbOperator, subjectID int) (*model.AccSubject, error) {
	strSql := "select " + strings.Join(accSubInfoFields, ",") + " from " + accSubInfoTN + " where subject_id=?"
	dao.Logger.DebugContext(ctx, "[accountSubject/db/GetAccSubByID] [sql: %s ,values: %d]", strSql, subjectID)
	var accSub = &model.AccSubject{}
	start := time.Now()
	defer func() {
		dao.Logger.InfoContext(ctx, "[accountSubject/db/GetAccSubByID] [SqlElapsed: %v]", time.Since(start))
	}()
	switch err := scanAccSubTask(do.QueryRowContext(ctx, strSql, subjectID), accSub); err {
	case nil:
		return accSub, nil
	case sql.ErrNoRows:
		return nil, err
	default:
		dao.Logger.ErrorContext(ctx, "[accountSubject/db/GetAccSubByID] [scanAccSubTask: %s]", err.Error())
		return nil, err
	}
}

// list count by filter
func (dao *AccSubDao) CountByFilter(ctx context.Context, do DbOperator, filter map[string]interface{}) (int64, error) {
	var c int64
	strSql, values := transferCountSql(accSubInfoTN, filter)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, values...).Scan(&c)
	dao.Logger.InfoContext(ctx, "[accountSubject/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

// list count by filter
func (dao *AccSubDao) Count(ctx context.Context, do DbOperator) (int64, error) {
	var c int64
	strSql := "select count(1) from " + accSubInfoTN
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql).Scan(&c)
	dao.Logger.InfoContext(ctx, "[accountSubject/db/CountByFilter] [SqlElapsed: %v]", time.Since(start))
	return c, err
}

// 检查在同一个company内是否subjectName和commonId是否有重复的记录
// 会计科目的名称，1级科目名称不能重复，但二级以后的科目名称是可以重复的。
func (dao *AccSubDao) CheckDuplication(ctx context.Context, do DbOperator, companyId int,
	commonId, subjectName string) (int64, error) {
	var c int64
	strSql := "select count(1) from " + accSubInfoTN +
		" where company_id = ? and (common_id = ?  or (subject_name = ? and subject_level = 1))"
	dao.Logger.DebugContext(ctx, "[accountSubject/db/CheckDuplication] [sql:%s,company_id: %d,commonId:%s,subject_name:%s]",
		strSql, companyId, commonId, subjectName)
	start := time.Now()
	err := do.QueryRowContext(ctx, strSql, companyId, commonId, subjectName).Scan(&c)
	dao.Logger.InfoContext(ctx, "[accountSubject/db/CheckDuplication] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[accountSubject/db/CheckDuplication] [do.Exec: %s]", err.Error())
		return 0, err
	}
	return c, err
}

func (dao *AccSubDao) Create(ctx context.Context, do DbOperator, st *model.AccSubject) error {
	strSql := "insert into " + accSubInfoTN +
		" (" + strings.Join(accSubInfoFields, ",") + ") values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	values := []interface{}{st.SubjectID, st.CompanyID, st.CommonID, st.SubjectName, st.SubjectLevel,
		st.SubjectDirection, st.SubjectType, st.MnemonicCode, st.SubjectStyle}
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

func (dao *AccSubDao) DeleteByID(ctx context.Context, do DbOperator, subjectID int) error {
	strSql := "delete from " + accSubInfoTN + " where subject_id = ?"

	dao.Logger.DebugContext(ctx, "[accountSubject/db/DeleteByID] [sql: %s, id: %d]", strSql, subjectID)
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

func (dao *AccSubDao) UpdateBySubID(ctx context.Context, do DbOperator, subjectID int,
	params map[string]interface{}) error {
	//var keyMap = map[string]string{"SubjectID": "subject_id", "SubjectName": "subject_name", "SubjectLevel": "subject_level"}
	strSql := "update " + accSubInfoTN + " set "
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
	dao.Logger.DebugContext(ctx, "[accountSubject/db/UpdateBySubID] [sql: %s, values: %v]", strSql, values)
	_, err := do.ExecContext(ctx, strSql, values...)
	dao.Logger.InfoContext(ctx, "[accountSubject/db/UpdateBySubID] [SqlElapsed: %v]", time.Since(start))
	if err != nil {
		dao.Logger.ErrorContext(ctx, "[accountSubject/db/UpdateBySubID] [do.Exec: %s]", err.Error())
		return err
	}
	return nil
}
