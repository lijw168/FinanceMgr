package service

import (
	"analysis-server/api/db"
	//"analysis-server/api/utils"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
	"context"
	"database/sql"
)

type AccountSubService struct {
	Logger     *log.Logger
	AccSubDao  *db.AccSubDao
	VRecordDao *db.VoucherRecordDao
	CompanyDao *db.CompanyDao
	Db         *sql.DB
}

func (as *AccountSubService) CreateAccSub(ctx context.Context, params *model.CreateSubjectParams,
	requestId string) (*model.AccSubjectView, CcError) {
	//create
	as.Logger.InfoContext(ctx, "CreateAccSub method start, "+"subjectName:%s", *params.SubjectName)
	bIsRollBack := true
	FuncName := "AccountSubService/accountSub/CreateAccSub"
	// Begin transaction
	tx, err := as.Db.Begin()
	if err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, as.Logger, FuncName, tx)
		}
	}()

	conflictCount, err := as.AccSubDao.CheckDuplication(ctx, tx, *params.CompanyID, *params.CommonID, *params.SubjectName)
	if err != nil {
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if conflictCount > 0 {
		return nil, NewError(ErrAccSub, ErrDuplicate, ErrNull, ErrFiledDuplicate)
	}
	//generate account subject
	accSub := new(model.AccSubject)
	accSub.SubjectName = *params.SubjectName
	accSub.SubjectLevel = *params.SubjectLevel
	accSub.CommonID = *params.CommonID
	accSub.CompanyID = *params.CompanyID
	accSub.SubjectDirection = *params.SubjectDirection
	accSub.SubjectType = *params.SubjectType
	accSub.SubjectStyle = *params.SubjectStyle
	if params.MnemonicCode != nil {
		accSub.MnemonicCode = *params.MnemonicCode
	} else {
		accSub.MnemonicCode = ""
	}
	accSub.SubjectID = GIdInfoService.genSubIdInfo.GetNextId()
	if err = as.AccSubDao.Create(ctx, tx, accSub); err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [AccSubDao.Create: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//commit
	if err = tx.Commit(); err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	accSubView := as.AccSubMdelToView(accSub)
	as.Logger.InfoContext(ctx, "CreateAccSub method end, "+"subjectName:%s", *params.SubjectName)
	return accSubView, nil
}

// convert accSubject to accSubjectView ...
func (as *AccountSubService) AccSubMdelToView(accSub *model.AccSubject) *model.AccSubjectView {
	accSubView := new(model.AccSubjectView)
	accSubView.SubjectID = accSub.SubjectID
	accSubView.SubjectName = accSub.SubjectName
	accSubView.SubjectLevel = accSub.SubjectLevel
	accSubView.CommonID = accSub.CommonID
	accSubView.CompanyID = accSub.CompanyID
	accSubView.SubjectType = accSub.SubjectType
	accSubView.SubjectDirection = accSub.SubjectDirection
	accSubView.MnemonicCode = accSub.MnemonicCode
	accSubView.SubjectStyle = accSub.SubjectStyle
	return accSubView
}

func (as *AccountSubService) CopyAccSubTemplate(ctx context.Context, iCompanyId int,
	requestId string) ([]*model.AccSubjectView, int, CcError) {
	as.Logger.InfoContext(ctx, "CopyAccSubTemplate method start, companyId:%s", iCompanyId)
	bIsRollBack := true
	FuncName := "AccountSubService/accountSub/CopyAccSubTemplate"
	// Begin transaction
	tx, err := as.Db.Begin()
	if err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, 0, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, as.Logger, FuncName, tx)
		}
	}()
	//list template account subject
	filterFields := make(map[string]interface{})
	//把超级管理所在公司的会计科目作为会计科目的模板。
	filterFields["companyId"] = 1
	limit, offset := -1, 0
	orderField := ""
	orderDirection := 0
	accSubInfos, err := as.AccSubDao.List(ctx, as.Db, filterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		as.Logger.ErrorContext(ctx, "[AccountSubService/service/ListAccSub] [AccSubDao.List: %s, filterFields: %v]", err.Error(), filterFields)
		return nil, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	accSubViewSlice := make([]*model.AccSubjectView, len(accSubInfos))
	for _, accSubInfo := range accSubInfos {
		//generate new account subject
		accSubInfo.CompanyID = iCompanyId
		accSubInfo.SubjectID = GIdInfoService.genSubIdInfo.GetNextId()
		accSubInfoView := as.AccSubMdelToView(accSubInfo)
		accSubViewSlice = append(accSubViewSlice, accSubInfoView)
		if err = as.AccSubDao.Create(ctx, tx, accSubInfo); err != nil {
			as.Logger.ErrorContext(ctx, "[%s] [AccSubDao.Create: %s]", FuncName, err.Error())
			return nil, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
	}
	//commit
	if err = tx.Commit(); err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	accSubInfoCount := len(accSubInfos)
	as.Logger.InfoContext(ctx, "CopyAccSubTemplate method end, companyId:%s", iCompanyId)
	return accSubViewSlice, accSubInfoCount, nil
}

func (as *AccountSubService) GetAccSubById(ctx context.Context, subjectID int,
	requestId string) (*model.AccSubjectView, CcError) {
	accSubject, err := as.AccSubDao.GetAccSubByID(ctx, as.Db, subjectID)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, NewCcError(cons.CodeAccSubNotExist, ErrAccSub, ErrNotFound, ErrNull, "the account subject is not exist")
	default:
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	accSubView := as.AccSubMdelToView(accSubject)
	return accSubView, nil
}

func (as *AccountSubService) getRefsOfAccSubID(ctx context.Context, subjectID int,
	tx *sql.Tx, isCalTotal bool) (error, int) {
	iTotalCount := 0
	//get company info
	comInfo, err := as.CompanyDao.GetCompanyByAccSubId(ctx, tx, subjectID)
	if err != nil {
		return err, iTotalCount
	}
	//获取有数据的会计年度
	iStartAccountYear := comInfo.StartAccountPeriod / 100
	iLatestAccountYear := comInfo.LatestAccountYear
	yearSlice := make([]int, 0, (iLatestAccountYear - iStartAccountYear + 1))
	for i := iStartAccountYear; i <= iLatestAccountYear; i++ {
		yearSlice = append(yearSlice, i)
	}
	//判断是否在使用
	for _, year := range yearSlice {
		filterFields := make(map[string]interface{})
		filterFields["subId1"] = subjectID
		var iCount int64
		iCount, err = as.VRecordDao.CountByFilter(ctx, tx, year, filterFields)
		if err != nil {
			as.Logger.ErrorContext(ctx, "[AccountSubService/service/JudgeAccSubReferenceBySubID] [VRecordDao.CountByFilter,Error info: %s", err.Error())
			return err, 0
		}
		if iCount > 0 {
			if isCalTotal {
				iTotalCount += int(iCount)
			} else {
				return nil, int(iCount)
			}
		}
	}
	return nil, iTotalCount
}

func (as *AccountSubService) DeleteAccSubByID(ctx context.Context, subjectID int,
	requestId string) CcError {
	as.Logger.InfoContext(ctx, "DeleteAccSubByID method begin, "+"subject:%d", subjectID)
	FuncName := "AccountSubService/accountSub/DeleteAccSubByID"
	bIsRollBack := true
	tx, err := as.Db.Begin()
	if err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, as.Logger, FuncName, tx)
		}
	}()
	iCount := 0
	if err, iCount = as.getRefsOfAccSubID(ctx, subjectID, tx, false); err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if iCount > 0 {
		return NewError(ErrAccSub, ErrError, ErrNull, "the account subjectID is using, don't delete.")
	}
	err = as.AccSubDao.DeleteByID(ctx, tx, subjectID)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, "Delete failed")
	}
	if err = tx.Commit(); err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	as.Logger.InfoContext(ctx, "DeleteAccSubByID method end, "+"subject:%d", subjectID)
	return nil
}

func (as *AccountSubService) UpdateAccSubById(ctx context.Context, subjectID int,
	params map[string]interface{}) CcError {
	FuncName := "AccountSubService/accountSub/UpdateAccSubById"
	bIsRollBack := true
	tx, err := as.Db.Begin()
	if err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, as.Logger, FuncName, tx)
		}
	}()

	iCount := 0
	if err, iCount = as.getRefsOfAccSubID(ctx, subjectID, tx, false); err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if iCount > 0 {
		return NewError(ErrAccSub, ErrError, ErrNull, "the account subjectID is using, don't updated.")
	}
	err = as.AccSubDao.UpdateBySubID(ctx, tx, subjectID, params)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	return nil
}

func (as *AccountSubService) UpdateYearBalanceById(ctx context.Context, subjectID int,
	dYearBalance float64) CcError {
	FuncName := "AccountSubService/accountSub/UpdateYearBalanceById"
	bIsRollBack := true
	tx, err := as.Db.Begin()
	if err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, as.Logger, FuncName, tx)
		}
	}()
	updateFields := make(map[string]interface{})
	updateFields["balance"] = dYearBalance
	err = as.AccSubDao.UpdateBySubID(ctx, tx, subjectID, updateFields)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	return nil
}

func (as *AccountSubService) ListYearBalance(ctx context.Context,
	params *model.ListParams) ([]*model.YearBalanceView, int, CcError) {
	balViewSlice := make([]*model.YearBalanceView, 0)
	filterFields := make(map[string]interface{})
	limit, offset := -1, 0
	if params.Filter != nil {
		for _, f := range params.Filter {
			switch *f.Field {
			case "subjectId", "companyId", "subjectName", "subjectLevel", "commonId":
				fallthrough
			case "subjectDirection", "subjectType", "mnemonicCode", "subjectStyle":
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
	yearBals, err := as.AccSubDao.ListYearBalance(ctx, as.Db, filterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		as.Logger.ErrorContext(ctx, "[AccountSubService/service/ListYearBalance] [AccSubDao.ListYearBalance: %s, filterFields: %v]", err.Error(), filterFields)
		return balViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	for _, yearBal := range yearBals {
		yearBalView := new(model.YearBalanceView)
		yearBalView.SubjectID = yearBal.SubjectID
		yearBalView.Balance = yearBal.Balance
		balViewSlice = append(balViewSlice, yearBalView)
	}
	yearBalsCount := len(yearBals)
	return balViewSlice, yearBalsCount, nil
}

func (as *AccountSubService) ListAccSub(ctx context.Context,
	params *model.ListSubjectParams) ([]*model.AccSubjectView, int, CcError) {
	accSubViewSlice := make([]*model.AccSubjectView, 0)
	filterFields := make(map[string]interface{})
	limit, offset := -1, 0
	if params.Filter != nil {
		for _, f := range params.Filter {
			switch *f.Field {
			case "subjectId", "companyId", "subjectName", "subjectLevel", "commonId":
				fallthrough
			case "subjectDirection", "subjectType", "mnemonicCode", "subjectStyle":
				filterFields[*f.Field] = f.Value
			default:
				return accSubViewSlice, 0, NewError(ErrAccSub, ErrUnsupported, ErrField, *f.Field)
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
	accSubInfos, err := as.AccSubDao.List(ctx, as.Db, filterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		as.Logger.ErrorContext(ctx, "[AccountSubService/service/ListAccSub] [AccSubDao.List: %s, filterFields: %v]", err.Error(), filterFields)
		return accSubViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	for _, accSubInfo := range accSubInfos {
		accSubInfoView := as.AccSubMdelToView(accSubInfo)
		accSubViewSlice = append(accSubViewSlice, accSubInfoView)
	}
	accSubInfoCount := len(accSubInfos)
	return accSubViewSlice, accSubInfoCount, nil
}

func (as *AccountSubService) QueryAccSubReferenceBySubID(ctx context.Context, subjectID int,
	requestId string) (int, CcError) {
	as.Logger.InfoContext(ctx, "QueryAccSubReferenceBySubID method begin, "+"subject:%d", subjectID)
	FuncName := "AccountSubService/accountSub/QueryAccSubReferenceBySubID"
	bIsRollBack := true
	tx, err := as.Db.Begin()
	if err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return 0, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, as.Logger, FuncName, tx)
		}
	}()

	iCount := 0
	if err, iCount = as.getRefsOfAccSubID(ctx, subjectID, tx, true); err != nil {
		return 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	as.Logger.InfoContext(ctx, "QueryAccSubReferenceBySubID method end, "+"subject:%d", subjectID)
	return iCount, nil
}

func (as *AccountSubService) GetYearBalanceById(ctx context.Context, subjectID int,
	requestId string) (float64, CcError) {
	if dBalanceValue, err := as.AccSubDao.GetBalanceByID(ctx, as.Db, subjectID); err != nil {
		return 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	} else {
		return dBalanceValue, nil
	}
}
