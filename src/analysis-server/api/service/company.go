package service

import (
	"context"
	"database/sql"

	"financeMgr/src/analysis-server/api/db"
	"financeMgr/src/analysis-server/api/utils"
	"financeMgr/src/analysis-server/model"
	cons "financeMgr/src/common/constant"
	"financeMgr/src/common/log"
	"fmt"
	"time"
)

type CompanyService struct {
	Logger          *log.Logger
	CompanyDao      *db.CompanyDao
	CompanyGroupDao *db.CompanyGroupDao
	Db              *sql.DB
}

func (cs *CompanyService) CreateCompany(ctx context.Context, params *model.CreateCompanyParams,
	requestId string) (*model.CompanyView, CcError) {
	//create
	cs.Logger.InfoContext(ctx, "CreateCompany method start, "+"companyName:%s", *params.CompanyName)
	FuncName := "CompanyService/Company/CreateCompany"
	//创建新表，如果存在，就不创建了。
	iVoucherYear := (*params.StartAccountPeriod) / 100
	baseTableName := []string{"voucherInfo", "voucherRecordInfo"}
	for _, tn := range baseTableName {
		err := cs.CompanyDao.CreateNewTable(ctx, cs.Db, tn, db.GenTableName(iVoucherYear, tn))
		if err != nil {
			errMsg := fmt.Sprintf("CreateNewTable,failed;errInfo:%s", err.Error())
			return nil, NewError(ErrSystem, ErrError, ErrNull, errMsg)
		}
	}
	bIsRollBack := true
	// Begin transaction
	tx, err := cs.Db.Begin()
	if err != nil {
		cs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, cs.Logger, FuncName, tx)
		}
	}()

	filterFields := make(map[string]interface{})
	filterFields["companyName"] = *params.CompanyName
	conflictCount, err := cs.CompanyDao.CountByFilter(ctx, tx, filterFields)
	if err != nil {
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if conflictCount > 0 {
		return nil, NewError(ErrCompany, ErrConflict, ErrNull, ErrRecordExist)
	}
	//generate company, 有潜在的bug，没有判断指针是否为空，就赋值了。
	comInfo := new(model.CompanyInfo)
	comInfo.CompanyName = *params.CompanyName
	comInfo.AbbrevName = *params.AbbrevName
	comInfo.Corporator = *params.Corporator
	comInfo.Phone = *params.Phone
	comInfo.Email = *params.Email
	comInfo.CompanyAddr = *params.CompanyAddr
	comInfo.Backup = *params.Backup
	comInfo.CreatedAt = time.Now()
	comInfo.CompanyID = GIdInfoService.genComIdInfo.GetNextId()
	comInfo.StartAccountPeriod = *params.StartAccountPeriod
	comInfo.LatestAccountYear = (*params.StartAccountPeriod) / 100
	if err = cs.CompanyDao.Create(ctx, tx, comInfo); err != nil {
		cs.Logger.ErrorContext(ctx, "[%s] [CompanyDao.Create: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		cs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	comView := cs.CompanyModelToView(comInfo)
	cs.Logger.InfoContext(ctx, "CreateCompany method end, "+"companyName:%s", *params.CompanyName)
	return comView, nil
}

// CompanyModelToView... ,convert accSubject to CompanyView ...
func (cs *CompanyService) CompanyModelToView(comInfo *model.CompanyInfo) *model.CompanyView {
	comView := new(model.CompanyView)
	comView.CompanyName = comInfo.CompanyName
	comView.AbbrevName = comInfo.AbbrevName
	comView.Corporator = comInfo.Corporator
	comView.Phone = comInfo.Phone
	comView.Email = comInfo.Email
	comView.CompanyAddr = comInfo.CompanyAddr
	comView.Backup = comInfo.Backup
	comView.StartAccountPeriod = comInfo.StartAccountPeriod
	comView.LatestAccountYear = comInfo.LatestAccountYear
	comView.CompanyID = comInfo.CompanyID
	comView.CreatedAt = comInfo.CreatedAt
	comView.UpdatedAt = comInfo.UpdatedAt
	comView.CompanyGroupID = comInfo.CompanyGroupID
	return comView
}

func (cs *CompanyService) GetCompanyById(ctx context.Context, companyId int,
	requestId string) (*model.CompanyView, CcError) {
	comInfo, err := cs.CompanyDao.Get(ctx, cs.Db, companyId)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, NewCcError(cons.CodeComInfoNotExist, ErrCompany, ErrNotFound, ErrNull, "the company information is not exist")
	default:
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	comInfoView := cs.CompanyModelToView(comInfo)
	return comInfoView, nil
}

// func (cs *CompanyService) DeleteCompanyByName(ctx context.Context, strCompanyName string,
// 	requestId string) CcError {
// 	cs.Logger.InfoContext(ctx, "DeleteCompanyByName method begin, "+"company Name:%s", strCompanyName)
// 	err := cs.CompanyDao.DeleteByName(ctx, cs.Db, strSubName)
// 	if err != nil {
// 		return NewError(ErrSystem, ErrError, ErrNull, "Delete failed")
// 	}
// 	cs.Logger.InfoContext(ctx, "DeleteCompanyByName method end, "+"company Name:%s", strCompanyName)
// 	return nil
// }

func (cs *CompanyService) DeleteCompanyByID(ctx context.Context, companyID int,
	requestId string) CcError {
	cs.Logger.InfoContext(ctx, "DeleteCompanyByID method begin, company ID:%d", companyID)
	err := cs.CompanyDao.Delete(ctx, cs.Db, companyID)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, "Delete failed")
	}
	cs.Logger.InfoContext(ctx, "DeleteCompanyByID method end, company ID:%d", companyID)
	return nil
}

func (cs *CompanyService) UpdateCompanyById(ctx context.Context, companyId int, params map[string]interface{}) CcError {
	FuncName := "CompanyService/Company/UpdateCompanyById"
	bIsRollBack := true
	// Begin transaction
	tx, err := cs.Db.Begin()
	if err != nil {
		cs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, cs.Logger, FuncName, tx)
		}
	}()
	//insure the company exist
	_, err = cs.CompanyDao.Get(ctx, tx, companyId)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return NewCcError(cons.CodeComInfoNotExist, ErrCompany, ErrNotFound, ErrNull, "the company is not exist")
	default:
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//update info
	params["UpdatedAt"] = time.Now()
	err = cs.CompanyDao.Update(ctx, tx, companyId, params)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		cs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	return nil
}

func (cs *CompanyService) ListCompany(ctx context.Context,
	params *model.ListCompanyParams) ([]*model.CompanyView, int, CcError) {
	comViewSlice := make([]*model.CompanyView, 0)
	filterFields := make(map[string]interface{})
	limit, offset := -1, 0
	if params.Filter != nil {
		for _, f := range params.Filter {
			switch *f.Field {
			// case "fuzzy_name":
			// 	volName := "%" + f.Value.(string) + "%"
			// 	fuzzyMatchFields["volume_name"] = volName
			case "companyId", "companyName", "abbreName", "corporator", "phone", "e_mail", "companyAddr":
				filterFields[*f.Field] = f.Value
			case "backup", "startAccountPeriod", "latestAccountYear", "companyGroupId":
				filterFields[*f.Field] = f.Value
			default:
				return comViewSlice, 0, NewError(ErrCompany, ErrUnsupported, ErrField, *f.Field)
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
	comInfos, err := cs.CompanyDao.List(ctx, cs.Db, filterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		cs.Logger.ErrorContext(ctx, "[CompanyService/service/ListCompany] [CompanyDao.List: %s, filterFields: %v]",
			err.Error(), filterFields)
		return comViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	for _, comInfo := range comInfos {
		comInfoView := cs.CompanyModelToView(comInfo)
		comViewSlice = append(comViewSlice, comInfoView)
	}
	comInfoCount := len(comViewSlice)
	//volumeCount, CcErr := vs.CountByFilter(ctx, vs.Db, filterFields)
	// if CcErr != nil {
	// 	return nil, 0, CcErr
	// }
	return comViewSlice, comInfoCount, nil
}

func (cs *CompanyService) AssociatedCompanyGroup(ctx context.Context, params *model.AssociatedCompanyGroupParams,
	requestId string) CcError {
	FuncName := "CompanyService/Company/AssociatedCompanyGroup"
	bIsRollBack := true
	// Begin transaction
	tx, err := cs.Db.Begin()
	if err != nil {
		cs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, cs.Logger, FuncName, tx)
		}
	}()
	//insure the company group exist
	filterFields := make(map[string]interface{})
	filterFields["companyGroupId"] = *params.CompanyGroupID
	filterFields["groupStatus"] = utils.ValidCompanyGroup
	conflictCount, err := cs.CompanyGroupDao.CountByFilter(ctx, tx, filterFields)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if conflictCount == 0 {
		return NewCcError(cons.CodeCompanyGroupNotExist, ErrCompany, ErrNotFound, ErrNull, "the company group is not exist")
	}
	//update company info
	updateFields := make(map[string]interface{})
	if *params.IsAttach {
		updateFields["companyGroupId"] = *params.CompanyGroupID
	} else {
		updateFields["companyGroupId"] = 0
	}
	updateFields["UpdatedAt"] = time.Now()
	err = cs.CompanyDao.Update(ctx, tx, *params.CompanyID, updateFields)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		cs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	return nil
}
