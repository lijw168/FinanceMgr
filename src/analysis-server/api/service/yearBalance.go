package service

import (
	"context"
	"database/sql"
	"financeMgr/src/analysis-server/api/db"
	"financeMgr/src/analysis-server/api/utils"
	"financeMgr/src/analysis-server/model"
	cons "financeMgr/src/common/constant"
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
	yearBal := new(model.YearBalance)
	yearBal.CompanyID = *params.CompanyID
	yearBal.SubjectID = *params.SubjectID
	yearBal.Year = *params.Year
	yearBal.Balance = *params.Balance
	//未结算状态
	yearBal.Status = utils.NoAnnualClosing
	yearBal.UpdatedAt = time.Now()
	yearBal.CreatedAt = time.Now()

	// Begin transaction
	bIsRollBack := true
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

func (ys *YearBalanceService) BatchCreateYearBalance(ctx context.Context, params *model.BatchCreateYearBalsParams) CcError {
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
	//以后优化时，可以通过insert into table(col1,...) values ()  进行优化
	yearBal := model.YearBalance{}
	for _, optSubAndBal := range params.OptSubAndBals {
		yearBal.CompanyID = *params.CompanyID
		yearBal.Year = *params.Year
		yearBal.SubjectID = *optSubAndBal.SubjectID
		yearBal.Balance = *optSubAndBal.Balance
		//未结算状态
		yearBal.Status = utils.NoAnnualClosing
		yearBal.UpdatedAt = time.Now()
		yearBal.CreatedAt = time.Now()
		err = ys.YearBalDao.CreateYearBalance(ctx, tx, &yearBal)
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

func (ys *YearBalanceService) GetAccSubYearBalValue(ctx context.Context, params *model.BasicYearBalanceParams,
	requestId string) (float64, CcError) {
	FuncName := "YearBalanceService/service/GetAccSubYearBalValue"
	bIsRollBack := true
	// Begin transaction
	tx, err := ys.Db.Begin()
	if err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return 0, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, ys.Logger, FuncName, tx)
		}
	}()
	var dBalanceValue float64
	if dBalanceValue, err = ys.YearBalDao.GetAccSubYearBalValue(ctx, ys.Db, params); err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, NewCcError(cons.CodeYearBalanceNotExist, ErrYearBalance, ErrNotFound, ErrNull, "the year balance record is not exist")
		default:
			return 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
	}
	if err = tx.Commit(); err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	ys.Logger.InfoContext(ctx, "GetAccSubYearBalValue method end")
	return dBalanceValue, nil
}

func (ys *YearBalanceService) GetYearBalance(ctx context.Context, params *model.BasicYearBalanceParams,
	requestId string) (*model.YearBalanceView, CcError) {

	FuncName := "YearBalanceService/service/GetYearBalance"
	bIsRollBack := true
	// Begin transaction
	tx, err := ys.Db.Begin()
	if err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, ys.Logger, FuncName, tx)
		}
	}()
	var yearBal *model.YearBalance
	if yearBal, err = ys.YearBalDao.GetYearBalance(ctx, ys.Db, params); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, NewCcError(cons.CodeYearBalanceNotExist, ErrYearBalance, ErrNotFound, ErrNull, "the year balance record is not exist")
		default:
			return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
	}
	if err = tx.Commit(); err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	//这里没有把科目的创建和更新时间以及公司ID返回到前段，那两个时间字段，仅在查数据时使用，companyId在前端已经知晓，无需返回。
	yearBalView := new(model.YearBalanceView)
	yearBalView.SubjectID = yearBal.SubjectID
	yearBalView.Balance = yearBal.Balance
	yearBalView.Year = yearBal.Year
	yearBalView.Status = yearBal.Status
	ys.Logger.InfoContext(ctx, "GetYearBalance method end")
	return yearBalView, nil
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

func (ys *YearBalanceService) BatchDeleteYearBalance(ctx context.Context, filter map[string]interface{}) CcError {
	ys.Logger.InfoContext(ctx, "BatchDeleteYearBalance method begin, update params:%v", filter)
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
	if err = ys.YearBalDao.BatchDeleteYearBalance(ctx, ys.Db, filter); err != nil {
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

// 该函数仅仅批量更新balance这一个字段。
func (ys *YearBalanceService) BatchUpdateBals(ctx context.Context, params *model.BatchUpdateBalsParams) CcError {
	ys.Logger.InfoContext(ctx, "BatchUpdateBals method begin, params:%v", params)
	FuncName := "YearBalanceService/service/BatchUpdateBals"
	bIsRollBack := true
	filter := make(map[string]interface{})
	filter["companyId"] = *params.CompanyID
	filter["year"] = *params.Year
	updateField := map[string]interface{}{}
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
	for _, v := range params.OptSubAndBals {
		filter["subjectId"] = *v.SubjectID
		updateField["balance"] = *v.Balance
		err = ys.YearBalDao.UpdateYearBalance(ctx, tx, filter, updateField)
		if err != nil {
			return NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
	}
	if err = tx.Commit(); err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	ys.Logger.InfoContext(ctx, "UpdateYearBalance method end")
	return nil
}

func (ys *YearBalanceService) UpdateYearBalance(ctx context.Context, filter map[string]interface{},
	updateField map[string]interface{}) CcError {
	ys.Logger.InfoContext(ctx, "UpdateYearBalance method begin, filter:%v,updateField:%v", filter, updateField)
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
	if err = ys.YearBalDao.UpdateYearBalance(ctx, tx, filter, updateField); err != nil {
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
	FuncName := "YearBalanceService/service/ListYearBalance"
	bIsRollBack := true
	// Begin transaction
	tx, err := ys.Db.Begin()
	if err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, 0, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, ys.Logger, FuncName, tx)
		}
	}()
	yearBals, err := ys.YearBalDao.ListYearBalance(ctx, tx, filterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		ys.Logger.ErrorContext(ctx, "[AccountSubService/service/ListYearBalance] [AccSubDao.ListYearBalance: %s, filterFields: %v]", err.Error(), filterFields)
		return nil, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	for _, yearBal := range yearBals {
		//这里没有把科目的创建和更新时间以及公司ID返回到前段，那两个时间字段，仅在查数据时使用，companyId在前端已经知晓，无需返回。
		yearBalView := new(model.YearBalanceView)
		yearBalView.SubjectID = yearBal.SubjectID
		yearBalView.Balance = yearBal.Balance
		yearBalView.Year = yearBal.Year
		yearBalView.Status = yearBal.Status
		balViewSlice = append(balViewSlice, yearBalView)
	}
	yearBalsCount := len(yearBals)
	return balViewSlice, yearBalsCount, nil
}

func (ys *YearBalanceService) AnnualClosing(ctx context.Context, params *model.BatchCreateYearBalsParams) CcError {
	ys.Logger.InfoContext(ctx, "AnnualClosing method start, create params:%v", params)

	FuncName := "YearBalanceService/service/AnnualClosing"
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
	//generate next year QC balance
	subIdSlice := []int{}
	yearBal := model.YearBalance{}
	for _, optSubAndBal := range params.OptSubAndBals {
		yearBal.CompanyID = *params.CompanyID
		yearBal.Year = *params.Year
		yearBal.SubjectID = *optSubAndBal.SubjectID
		yearBal.Balance = *optSubAndBal.Balance
		//未结算状态
		yearBal.Status = utils.NoAnnualClosing
		yearBal.UpdatedAt = time.Now()
		yearBal.CreatedAt = time.Now()
		//generate next year QC balance，optimization:use insert into table(...) values (...),(...)
		if err = ys.YearBalDao.CreateYearBalance(ctx, tx, &yearBal); err != nil {
			return NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
		subIdSlice = append(subIdSlice, *optSubAndBal.SubjectID)
	}
	//modify the year data's status,ananual closing status
	filter := make(map[string]interface{})
	filter["companyId"] = *params.CompanyID
	//年度结算是增加的下一年的年度余额，并且更新的上一年的结算状态
	filter["year"] = *params.Year - 1
	filter["subjectId"] = subIdSlice
	updateField := map[string]interface{}{"status": utils.AnnualClosing}
	if err = ys.YearBalDao.UpdateYearBalance(ctx, tx, filter, updateField); err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	ys.Logger.InfoContext(ctx, "BatchCreateYearBalance method end")
	return nil
}

func (ys *YearBalanceService) CancelAnnualClosing(ctx context.Context, params *model.BatchDelYearBalsParams) CcError {
	ys.Logger.InfoContext(ctx, "CancelAnnualClosing method start, create params:%v", params)
	FuncName := "YearBalanceService/service/CancelAnnualClosing"
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
	filterFields["companyId"] = *params.CompanyID
	filterFields["year"] = *params.Year
	filterFields["subjectId"] = params.SubjectIDs
	if err = ys.YearBalDao.BatchDeleteYearBalance(ctx, ys.Db, filterFields); err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//recover the year data's status
	filterFields["year"] = *params.Year - 1
	updateField := map[string]interface{}{"status": utils.NoAnnualClosing}
	if err = ys.YearBalDao.UpdateYearBalance(ctx, tx, filterFields, updateField); err != nil {
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

func (ys *YearBalanceService) GetAnnualClosingStatus(ctx context.Context, companyID, year int) (int, CcError) {
	FuncName := "YearBalanceService/service/GetAnnualClosingStatus"
	bIsRollBack := true
	// Begin transaction
	tx, err := ys.Db.Begin()
	if err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return 0, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, ys.Logger, FuncName, tx)
		}
	}()
	filterFields := make(map[string]interface{})
	filterFields["companyId"] = companyID
	filterFields["year"] = year
	var iStatus int
	if iStatus, err = ys.YearBalDao.GetAccSubYearStatus(ctx, ys.Db, filterFields); err != nil {
		return utils.NoAnnualClosing, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		ys.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	ys.Logger.InfoContext(ctx, "GetAnnualClosingStatus method end")
	return iStatus, nil
}
