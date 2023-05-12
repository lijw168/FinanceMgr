package service

import (
	"context"
	"database/sql"
	"financeMgr/src/analysis-server/api/db"
	"financeMgr/src/analysis-server/model"
	"financeMgr/src/common/log"
)

type YearBalanceService struct {
	Logger     *log.Logger
	YearBalDao *db.YearBalanceDao
	Db         *sql.DB
}

func (ys *YearBalanceService) CreateYearBalance(ctx context.Context, params *model.OptYearBalanceParams,
	requestId string) CcError {
	ys.Logger.InfoContext(ctx, "CreateYearBalance method start, create params:%v", params)
	if err := ys.YearBalDao.CreateYearBalance(ctx, ys.Db, params); err != nil {
		FuncName := "YearBalanceService/yearBalance/CreateYearBalance"
		ys.Logger.ErrorContext(ctx, "[%s] [YearBalDao.Create: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	ys.Logger.InfoContext(ctx, "CreateYearBalance method end, create params:%v", params)
	return nil
}

func (ys *YearBalanceService) BatchCreateYearBalance(ctx context.Context, params []*model.OptYearBalanceParams,
	requestId string) CcError {
	ys.Logger.InfoContext(ctx, "BatchCreateYearBalance method start, create params:%v", params)

	FuncName := "YearBalanceService/service/BatchUpdateYearBalance"
	bIsRollBack := true
	// Begin transaction
	tx, err := ys.Db.Begin()
	if err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, ys.Logger, FuncName, tx)
		}
	}()
	for _, param := range params {
		if param.Year == nil || *(param.Year) <= 0 {
			return NewError(ErrYearBalance, ErrMiss, ErrYear, ErrNull)
		}
		if param.SubjectID == nil || *(param.SubjectID) <= 0 {
			return NewError(ErrYearBalance, ErrMiss, ErrId, ErrNull)
		}
		if param.Balance == nil {
			return NewError(ErrYearBalance, ErrMiss, ErrBalance, ErrNull)
		}
		err := ys.YearBalDao.CreateYearBalance(ctx, tx, param)
		if err != nil {
			return NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
	}
	if err = tx.Commit(); err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	ys.Logger.InfoContext(ctx, "BatchCreateYearBalance method end, create params:%v", params)
	return nil
}

func (ys *YearBalanceService) GetYearBalance(ctx context.Context, iYear, subjectID int,
	requestId string) (float64, CcError) {
	if dBalanceValue, err := ys.YearBalDao.GetYearBalance(ctx, ys.Db, iYear, subjectID); err != nil {
		return 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	} else {
		return dBalanceValue, nil
	}
}

func (ys *YearBalanceService) DeleteYearBalance(ctx context.Context, iYear, subjectID int,
	requestId string) CcError {
	ys.Logger.InfoContext(ctx, "DeleteYearBalance method begin, year:%d,Id:%v", iYear, subjectID)
	err := ys.YearBalDao.DeleteYearBalance(ctx, ys.Db, iYear, subjectID)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	ys.Logger.InfoContext(ctx, "DeleteYearBalance method end, year:%d,Id:%v", iYear, subjectID)
	return nil
}

func (ys *YearBalanceService) BatchDeleteYearBalance(ctx context.Context, params *model.BatchDelYearBalsParams,
	requestId string) CcError {
	ys.Logger.InfoContext(ctx, "BatchDeleteYearBalance method begin, update params:%v", params)
	FuncName := "YearBalanceService/service/BatchDeleteYearBalance"
	bIsRollBack := true
	// Begin transaction
	tx, err := ys.Db.Begin()
	if err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, ys.Logger, FuncName, tx)
		}
	}()
	filterFields := make(map[string]interface{})
	filterFields["subject_id"] = params.SubjectIDs
	filterFields["year"] = params.Year
	if err = ys.YearBalDao.BatchDeleteYearBalance(ctx, ys.Db, filterFields); err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	ys.Logger.InfoContext(ctx, "BatchDeleteYearBalance method end, update params:%v", params)
	return nil
}

func (ys *YearBalanceService) UpdateYearBalance(ctx context.Context, iYear, subjectID int,
	params map[string]interface{}) CcError {
	err := ys.YearBalDao.UpdateYearBalance(ctx, ys.Db, iYear, subjectID, params)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	return nil
}

func (ys *YearBalanceService) BatchUpdateYearBalance(ctx context.Context, params []*model.OptYearBalanceParams,
	requestId string) CcError {
	ys.Logger.InfoContext(ctx, "BatchUpdateYearBalance method begin, update params:%v", params)
	FuncName := "YearBalanceService/service/BatchUpdateYearBalance"
	bIsRollBack := true
	// Begin transaction
	tx, err := ys.Db.Begin()
	if err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, ys.Logger, FuncName, tx)
		}
	}()
	for _, param := range params {
		if param.SubjectID == nil || *param.SubjectID <= 0 {
			return NewError(ErrYearBalance, ErrMiss, ErrId, ErrNull)
		}
		if param.Year == nil || *param.Year <= 0 {
			return NewError(ErrYearBalance, ErrMiss, ErrYear, ErrNull)
		}
		if param.Balance == nil {
			return NewError(ErrYearBalance, ErrMiss, ErrChangeContent, ErrNull)
		}
		err := ys.YearBalDao.UpdateBalance(ctx, tx, param)
		if err != nil {
			return NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
	}
	if err = tx.Commit(); err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	ys.Logger.InfoContext(ctx, "BatchUpdateYearBalance method end, update params:%v", params)
	return nil
}

func (ys *YearBalanceService) ListYearBalance(ctx context.Context,
	params *model.ListParams) ([]*model.YearBalanceView, int, CcError) {
	balViewSlice := make([]*model.YearBalanceView, 0)
	filterFields := make(map[string]interface{})
	limit, offset := -1, 0
	if params.Filter != nil {
		for _, f := range params.Filter {
			switch *f.Field {
			case "subjectId", "year":
				filterFields[*f.Field] = f.Value
			default:
				return balViewSlice, 0, NewError(ErrAccSub, ErrUnsupported, ErrField, *f.Field)
			}
		}
	}
	if params.DescLimit != nil {
		limit = *params.DescLimit
		if params.DescOffset != nil {
			offset = *params.DescOffset
		}
	}
	orderField := ""
	orderDirection := 0
	if params.Order != nil {
		orderField = *params.Order[0].Field
		orderDirection = *params.Order[0].Direction
	}
	yearBals, err := ys.YearBalDao.ListYearBalance(ctx, ys.Db, filterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		ys.Logger.ErrorContext(ctx, "[AccountSubService/service/ListYearBalance] [AccSubDao.ListYearBalance: %s, filterFields: %v]", err.Error(), filterFields)
		return balViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	for _, yearBal := range yearBals {
		//这里没有把科目的创建和更新时间返回到前段，这两个字段，仅在查数据时使用。
		yearBalView := new(model.YearBalanceView)
		yearBalView.SubjectID = yearBal.SubjectID
		yearBalView.Balance = yearBal.Balance
		yearBalView.Year = yearBal.Year
		balViewSlice = append(balViewSlice, yearBalView)
	}
	yearBalsCount := len(yearBals)
	return balViewSlice, yearBalsCount, nil
}
