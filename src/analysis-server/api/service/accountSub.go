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
	Logger    *log.Logger
	AccSubDao *db.AccSubDao
	Db        *sql.DB
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
	return accSubView
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

func (as *AccountSubService) DeleteAccSubByID(ctx context.Context, subjectID int,
	requestId string) CcError {
	as.Logger.InfoContext(ctx, "DeleteAccSubByID method begin, "+"subject:%d", subjectID)
	err := as.AccSubDao.DeleteByID(ctx, as.Db, subjectID)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, "Delete failed")
	}
	as.Logger.InfoContext(ctx, "DeleteAccSubByID method end, "+"subject:%d", subjectID)
	return nil
}

func (as *AccountSubService) UpdateAccSubById(ctx context.Context, subjectID int, params map[string]interface{}) CcError {
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

	_, err = as.AccSubDao.GetAccSubByID(ctx, tx, subjectID)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return NewCcError(cons.CodeAccSubNotExist, ErrAccSub, ErrNotFound, ErrNull, "the account subject is not exist")
	default:
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//update info
	//params["UpdatedAt"] = time.Now()
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

func (as *AccountSubService) ListAccSub(ctx context.Context,
	params *model.ListSubjectParams) ([]*model.AccSubjectView, int, CcError) {
	accSubViewSlice := make([]*model.AccSubjectView, 0)
	filterFields := make(map[string]interface{})
	limit, offset := -1, 0
	if params.Filter != nil {
		for _, f := range params.Filter {
			switch *f.Field {
			case "subjectId", "companyId", "subjectName", "subjectLevel", "commonId":
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
