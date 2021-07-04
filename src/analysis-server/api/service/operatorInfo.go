package service

import (
	"analysis-server/api/db"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
	"context"
	"database/sql"
	"time"
)

type OperatorInfoService struct {
	Logger     *log.Logger
	OptInfoDao *db.OperatorInfoDao
	Db         *sql.DB
}

func (ps *OperatorInfoService) CreateOptInfo(ctx context.Context, params *model.CreateOptInfoParams,
	requestId string) (*model.OperatorInfoView, CcError) {
	ps.Logger.InfoContext(ctx, "CreateOptInfo method start, "+"operator Name:%s", *params.Name)

	FuncName := "OperatorInfoService/operater/CreateOptInfo"
	bIsRollBack := true
	tx, err := ps.Db.Begin()
	if err != nil {
		ps.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func(bRollBack bool) {
		if bRollBack {
			RollbackLog(ctx, ps.Logger, FuncName, tx)
		}
	}(bIsRollBack)

	filterFields := make(map[string]interface{})
	filterFields["name"] = *params.Name
	conflictCount, err := ps.OptInfoDao.CountByFilter(ctx, tx, filterFields)
	if err != nil {
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if conflictCount > 0 {
		return nil, NewError(ErrCompany, ErrConflict, ErrNull, ErrRecordExist)
	}
	//generate company
	optInfo := new(model.OperatorInfo)
	optInfo.Name = *params.Name
	optInfo.Password = *params.Password
	optInfo.CompanyID = *params.CompanyID
	optInfo.Job = *params.Job
	optInfo.Department = *params.Department
	optInfo.Status = *params.Status
	optInfo.Role = *params.Role
	optInfo.CreatedAt = time.Now()

	if err = ps.OptInfoDao.Create(ctx, tx, optInfo); err != nil {
		ps.Logger.ErrorContext(ctx, "[%s] [OptInfoDao.Create: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		ps.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	optView := ps.OperateInfoMdelToView(optInfo)
	ps.Logger.InfoContext(ctx, "CreateOptInfo method end, "+"operator Name:%s", *params.Name)
	return optView, nil
}

func (ps *OperatorInfoService) ListOperators(ctx context.Context,
	params *model.ListOperatorsParams) ([]*model.OperatorInfoView, int, CcError) {
	OptViewSlice := make([]*model.OperatorInfoView, 0)
	filterFields := make(map[string]interface{})
	//fuzzyMatchFields := make(map[string]string)
	limit, offset := -1, 0
	if params.Filter != nil {
		for _, f := range params.Filter {
			switch *f.Field {
			case "name", "company_id", "job", "department", "status", "role":
				filterFields[*f.Field] = f.Value
			default:
				return OptViewSlice, 0, NewError(ErrOperator, ErrUnsupported, ErrField, *f.Field)
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
	optInfos, err := ps.OptInfoDao.List(ctx, ps.Db, filterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		ps.Logger.ErrorContext(ctx, "[OperatorInfoService/service/ListOperators] [OptInfoDao.List: %s, filterFields: %v]", err.Error(), filterFields)
		return OptViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	for _, optInfo := range optInfos {
		optInfoView := ps.OperateInfoMdelToView(optInfo)
		OptViewSlice = append(OptViewSlice, optInfoView)
	}
	optInfoCount := len(optInfos)
	return OptViewSlice, optInfoCount, nil
}

// convert accSubject to OperatorInfoView ...
func (ps *OperatorInfoService) OperateInfoMdelToView(optInfo *model.OperatorInfo) *model.OperatorInfoView {
	optView := new(model.OperatorInfoView)
	optView.Name = optInfo.Name
	optView.Password = optInfo.Password
	optView.CompanyID = optInfo.CompanyID
	optView.Job = optInfo.Job
	optView.Department = optInfo.Department
	optView.Status = optInfo.Status
	optView.Role = optInfo.Role
	optView.CreatedAt = optInfo.CreatedAt
	optView.UpdatedAt = optInfo.UpdatedAt
	return optView
}

func (ps *OperatorInfoService) GetOperatorInfoByName(ctx context.Context, strOperatorName string,
	requestId string) (*model.OperatorInfoView, CcError) {
	optInfo, err := ps.OptInfoDao.Get(ctx, ps.Db, strOperatorName)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, NewCcError(cons.CodeOptInfoNotExist, ErrOperator, ErrNotFound, ErrNull, "the operator information is not exist")
	default:
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	optInfoView := ps.OperateInfoMdelToView(optInfo)
	return optInfoView, nil
}

func (ps *OperatorInfoService) DeleteOperatorInfoByName(ctx context.Context, strOperatorName string,
	requestId string) CcError {
	ps.Logger.InfoContext(ctx, "DeleteOperatorInfoByName method begin, "+"operator Name:%s", strOperatorName)
	err := ps.OptInfoDao.Delete(ctx, ps.Db, strOperatorName)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, "Delete failed")
	}
	ps.Logger.InfoContext(ctx, "DeleteOperatorInfoByName method end, "+"operator Name:%s", strOperatorName)
	return nil
}

func (ps *OperatorInfoService) UpdateOperator(ctx context.Context, strOptName string, params map[string]interface{}) CcError {
	FuncName := "OperatorInfoService/UpdateOperator"
	bIsRollBack := true
	tx, err := ps.Db.Begin()
	if err != nil {
		ps.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func(bRollBack bool) {
		if bRollBack {
			RollbackLog(ctx, ps.Logger, FuncName, tx)
		}
	}(bIsRollBack)
	_, err = ps.OptInfoDao.Get(ctx, tx, strOptName)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return NewCcError(cons.CodeOptInfoNotExist, ErrOperator, ErrNotFound, ErrNull, "the operator information is not exist")
	default:
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//update info
	params["UpdatedAt"] = time.Now()
	err = ps.OptInfoDao.Update(ctx, tx, strOptName, params)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		ps.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	return nil
}
