package service

import (
	"context"
	"database/sql"

	"financeMgr/src/analysis-server/api/db"
	"financeMgr/src/analysis-server/model"
	cons "financeMgr/src/common/constant"
	"financeMgr/src/common/log"
	"time"
)

type CompanyGroupService struct {
	Logger      *log.Logger
	ComGroupDao *db.CompanyGroupDao
	Db          *sql.DB
}

func (cg *CompanyGroupService) CreateCompanyGroup(ctx context.Context, params *model.CreateCompanyGroupParams,
	requestId string) (*model.CompanyGroupView, CcError) {
	//create
	cg.Logger.InfoContext(ctx, "CreateCompanyGroup method start, "+"groupName:%s", *params.GroupName)
	FuncName := "CompanyGroupService/Company/CreateCompanyGroup"
	bIsRollBack := true
	// Begin transaction
	tx, err := cg.Db.Begin()
	if err != nil {
		cg.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, cg.Logger, FuncName, tx)
		}
	}()

	filterFields := make(map[string]interface{})
	filterFields["groupName"] = *params.GroupName
	conflictCount, err := cg.ComGroupDao.CountByFilter(ctx, tx, filterFields)
	if err != nil {
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if conflictCount > 0 {
		return nil, NewError(ErrComGroup, ErrConflict, ErrNull, ErrRecordExist)
	}
	//generate company
	comGroup := new(model.CompanyGroup)
	comGroup.GroupName = *params.GroupName
	comGroup.GroupStatus = *params.GroupStatus
	comGroup.CreatedAt = time.Now()
	comGroup.CompanyGroupID = GIdInfoService.genComGroupIdInfo.GetNextId()
	if err = cg.ComGroupDao.Create(ctx, tx, comGroup); err != nil {
		cg.Logger.ErrorContext(ctx, "[%s] [ComGroupDao.Create: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		cg.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	comView := cg.CompanyGroupModelToView(comGroup)
	cg.Logger.InfoContext(ctx, "CreateCompanyGroup method end, "+"companyGroupName:%s", *params.GroupName)
	return comView, nil
}

// CompanyGroupModelToView... ,convert CompanyGroup to CompanyGroupView ...
func (cg *CompanyGroupService) CompanyGroupModelToView(comGroup *model.CompanyGroup) *model.CompanyGroupView {
	comGroupView := new(model.CompanyGroupView)
	comGroupView.GroupName = comGroup.GroupName
	comGroupView.GroupStatus = comGroup.GroupStatus
	comGroupView.CompanyGroupID = comGroup.CompanyGroupID
	comGroupView.CreatedAt = comGroup.CreatedAt
	comGroupView.UpdatedAt = comGroup.UpdatedAt
	return comGroupView
}

func (cg *CompanyGroupService) GetCompanyGroupById(ctx context.Context, companyGroupId int,
	requestId string) (*model.CompanyGroupView, CcError) {
	comGroupInfo, err := cg.ComGroupDao.Get(ctx, cg.Db, companyGroupId)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, NewCcError(cons.CodeCompanyGroupNotExist, ErrComGroup, ErrNotFound, ErrNull, "the company group is not exist")
	default:
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	comGroupView := cg.CompanyGroupModelToView(comGroupInfo)
	return comGroupView, nil
}

func (cg *CompanyGroupService) DeleteCompanyGroupByID(ctx context.Context, companyGroupId int,
	requestId string) CcError {
	cg.Logger.InfoContext(ctx, "DeleteCompanyGroupByID method begin, company ID:%d", companyGroupId)
	err := cg.ComGroupDao.Delete(ctx, cg.Db, companyGroupId)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, "Delete failed")
	}
	cg.Logger.InfoContext(ctx, "DeleteCompanyGroupByID method end, company ID:%d", companyGroupId)
	return nil
}

func (cg *CompanyGroupService) UpdateCompanyGroupById(ctx context.Context, companyGroupId int,
	params map[string]interface{}) CcError {
	FuncName := "CompanyGroupService/CompanyGroup/UpdateCompanyGroupById"
	bIsRollBack := true
	// Begin transaction
	tx, err := cg.Db.Begin()
	if err != nil {
		cg.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, cg.Logger, FuncName, tx)
		}
	}()
	//insure the company exist
	_, err = cg.ComGroupDao.Get(ctx, tx, companyGroupId)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return NewCcError(cons.CodeCompanyGroupNotExist, ErrComGroup, ErrNotFound, ErrNull, "the company group is not exist")
	default:
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//update info
	params["UpdatedAt"] = time.Now()
	err = cg.ComGroupDao.Update(ctx, tx, companyGroupId, params)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		cg.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	return nil
}

func (cg *CompanyGroupService) ListCompanyGroup(ctx context.Context,
	params *model.ListParams) ([]*model.CompanyGroupView, int, CcError) {
	comGroupViewSlice := make([]*model.CompanyGroupView, 0)
	filterFields := make(map[string]interface{})
	limit, offset := -1, 0
	if params.Filter != nil {
		for _, f := range params.Filter {
			switch *f.Field {
			case "companyGroupId", "groupName", "groupStatus":
				filterFields[*f.Field] = f.Value
			default:
				return comGroupViewSlice, 0, NewError(ErrComGroup, ErrUnsupported, ErrField, *f.Field)
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
	comGroups, err := cg.ComGroupDao.List(ctx, cg.Db, filterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		cg.Logger.ErrorContext(ctx, "[CompanyGroupService/service/ListCompanyGroup] [ComGroupDao.List: %s, filterFields: %v]",
			err.Error(), filterFields)
		return comGroupViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	for _, comGroup := range comGroups {
		comGroupView := cg.CompanyGroupModelToView(comGroup)
		comGroupViewSlice = append(comGroupViewSlice, comGroupView)
	}
	comInfoCount := len(comGroupViewSlice)
	return comGroupViewSlice, comInfoCount, nil
}
