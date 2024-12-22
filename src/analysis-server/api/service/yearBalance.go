package service

import (
	"context"
	"database/sql"
	"financeMgr/src/analysis-server/api/db"
	"financeMgr/src/analysis-server/model"
	"financeMgr/src/common/log"
	"time"
)

type YearBalanceService struct {
	Logger     *log.Logger
	YearBalDao *db.YearBalanceDao
	Db         *sql.DB
}

func (ys *YearBalanceService) CreateYearBalance(ctx context.Context, params *model.OptYearBalanceParams,
	requestId string) CcError {
	ys.Logger.InfoContext(ctx, "CreateYearBalance method start, create params:%v", params)
	FuncName := "YearBalanceService/yearBalance/CreateYearBalance"
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

	yearBal := new(model.YearBalance)
	yearBal.CompanyID = *params.CompanyID
	yearBal.SubjectID = *params.SubjectID
	yearBal.Year = *params.Year
	yearBal.Balance = *params.Balance
	yearBal.UpdatedAt = time.Now()
	yearBal.CreatedAt = time.Now()

	if err := ys.YearBalDao.CreateYearBalance(ctx, tx, yearBal); err != nil {

		ys.Logger.ErrorContext(ctx, "[%s] [YearBalDao.Create: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	ys.Logger.InfoContext(ctx, "CreateYearBalance method end")
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
	yearBal := new(model.YearBalance)
	for _, param := range params {
		if param.CompanyID == nil || *param.CompanyID <= 0 {
			return NewError(ErrYearBalance, ErrMiss, ErrCompanyId, ErrNull)
		}
		if param.Year == nil || *param.Year <= 0 {
			return NewError(ErrYearBalance, ErrMiss, ErrYear, ErrNull)
		}
		if param.SubjectID == nil || *param.SubjectID <= 0 {
			return NewError(ErrYearBalance, ErrMiss, ErrId, ErrNull)
		}
		if param.Balance == nil {
			return NewError(ErrYearBalance, ErrMiss, ErrBalance, ErrNull)
		}
		yearBal.CompanyID = *param.CompanyID
		yearBal.SubjectID = *param.SubjectID
		yearBal.Year = *param.Year
		yearBal.Balance = *param.Balance
		yearBal.UpdatedAt = time.Now()
		yearBal.CreatedAt = time.Now()
		err := ys.YearBalDao.CreateYearBalance(ctx, tx, yearBal)
		if err != nil {
			return NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
	}
	if err = tx.Commit(); err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	ys.Logger.InfoContext(ctx, "BatchCreateYearBalance method end")
	return nil
}

func (ys *YearBalanceService) GetYearBalance(ctx context.Context, params *model.BasicYearBalanceParams,
	requestId string) (float64, CcError) {
	if dBalanceValue, err := ys.YearBalDao.GetYearBalance(ctx, ys.Db, params); err != nil {
		return 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	} else {
		return dBalanceValue, nil
	}
}

func (ys *YearBalanceService) DeleteYearBalance(ctx context.Context, params *model.BasicYearBalanceParams,
	requestId string) CcError {
	ys.Logger.InfoContext(ctx, "DeleteYearBalance method begin, params:%v", params)
	FuncName := "YearBalanceService/service/DeleteYearBalance"
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

	err = ys.YearBalDao.DeleteYearBalance(ctx, tx, params)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	ys.Logger.InfoContext(ctx, "DeleteYearBalance method end")
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
	filterFields["companyId"] = params.CompanyID
	filterFields["year"] = params.Year
	filterFields["subject_id"] = params.SubjectIDs
	if err = ys.YearBalDao.BatchDeleteYearBalance(ctx, ys.Db, filterFields); err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	ys.Logger.InfoContext(ctx, "BatchDeleteYearBalance method end")
	return nil
}

func (ys *YearBalanceService) UpdateYearBalance(ctx context.Context, params *model.OptYearBalanceParams) CcError {
	ys.Logger.InfoContext(ctx, "UpdateYearBalance method begin, update params:%v", params)
	FuncName := "YearBalanceService/service/UpdateYearBalance"
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
	err = ys.YearBalDao.UpdateBalance(ctx, tx, params)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	if err = tx.Commit(); err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	ys.Logger.InfoContext(ctx, "UpdateYearBalance method end")
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
		if param.CompanyID == nil || *param.CompanyID <= 0 {
			return NewError(ErrYearBalance, ErrMiss, ErrCompanyId, ErrNull)
		}
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
	ys.Logger.InfoContext(ctx, "BatchUpdateYearBalance method end")
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
			case "companyId", "subjectId", "year":
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
		//这里没有把科目的创建和更新时间以及公司ID返回到前段，那两个时间字段，仅在查数据时使用，companyId在前端已经知晓，无需返回。
		yearBalView := new(model.YearBalanceView)
		yearBalView.SubjectID = yearBal.SubjectID
		yearBalView.Balance = yearBal.Balance
		yearBalView.Year = yearBal.Year
		balViewSlice = append(balViewSlice, yearBalView)
	}
	yearBalsCount := len(yearBals)
	return balViewSlice, yearBalsCount, nil
}
