package service

import (
	"context"
	"database/sql"

	cons "common/constant"
	"common/log"
	"analysis-server/api/db"
	"analysis-server/api/utils"
	"analysis-server/model"
)

type CompanyService struct {
	Logger     *log.Logger
	CompanyDao *db.CompanyDao
	Db         *sql.DB
	GenComId   *utils.GenIdInfo
}

func (cs *CompanyService) CreateCompany(ctx context.Context, params *model.CreateCompanyParams,
	requestId string) (*model.CompanyView, CcError) {
	//create
	cs.Logger.InfoContext(ctx, "CreateCompany method start, "+"companyName:%s", *params.CompanyName)
	FuncName := "CompanyService/Company/CreateCompany"
	// Begin transaction
	tx, err := cs.Db.Begin()
	if err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer RollbackLog(ctx, cs.Logger, FuncName, tx)
	filterFields := make(map[string]interface{})
	filterFields["companyName"] = *params.CompanyName
	conflictCount, err := cs.CompanyDao.CountByFilter(ctx, tx, filterFields)
	if err != nil {
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if conflictCount > 0 {
		return nil, NewError(ErrCompany, ErrConflict, ErrNull, err.Error())
	}
	//get the count of the table company
	// var comCount int
	// comCount, err = cs.CompanyDao.Count(ctx, tx)
	// if err != nil {
	// 	return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	// }
	//generate company
	comInfo := new(model.CompanyInfo)
	comInfo.CompanyName = *params.CompanyName
	comInfo.AbbrevName = *params.AbbrevName
	comInfo.Corporator = *params.Corporator
	comInfo.Phone = *params.Phone
	comInfo.Summary = *params.Summary
	comInfo.Email = *params.Email
	comInfo.CompanyAddr = *params.CompanyAddr
	comInfo.Backup = *params.Backup
	//>100,as subjectId
	comInfo.CompanyID = cs.GenComId.GetId()
	if err = cs.CompanyDao.create(ctx, tx, comInfo); err != nil {
		cs.Logger.ErrorContext(ctx, "[%s] [CompanyDao.Create: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		ps.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	comView := cs.CompanyModelToView(comInfo)
	cs.Logger.InfoContext(ctx, "CreateCompany method end, "+"companyName:%s", *params.CompanyName)
	return comView, nil
}

// convert accSubject to CompanyView ...
func (cs *CompanyService) CompanyModelToView(comInfo *model.CompanyInfo) *model.CompanyView {
	comView := new(model.CompanyView)
	comView.CompanyName = comInfo.CompanyName
	comView.AbbrevName = comInfo.AbbrevName
	comViewCorporator = comInfo.Corporator
	comView.Phone = comInfo.Phone
	comView.Summary = comInfo.Summary
	comView.Email = comInfo.Email
	comView.CompanyAddr = comInfo.CompanyAddr
	comView.Backup = comInfo.Backup
	comView.CompanyID = comInfo.CompanyID
	return comView
}

func (cs *CompanyService) GetCompanyById(ctx context.Context, companyId string,
	requestId string) (*model.CompanyView, CcError) {
	comInfo, err := cs.CompanyDao.Get(ctx, vs.Db, companyId)
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
	cs.Logger.InfoContext(ctx, "DeleteCompanyByID method begin, "+"company ID:%s", companyID)
	err := cs.CompanyDao.Delete(ctx, cs.Db, companyID)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, "Delete failed")
	}
	cs.Logger.InfoContext(ctx, "DeleteCompanyByID method end, "+"company ID:%s", companyID)
	return nil
}

func (cs *CompanyService) UpdateCompanyById(ctx context.Context, companyId string, params map[string]interface{}) CcError {
	FuncName := "CompanyService/Company/UpdateCompanyById"
	// Begin transaction
	tx, err := cs.Db.Begin()
	if err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer RollbackLog(ctx, cs.Logger, FuncName, tx)
	//insure the volume exist
	company, err := cs.CompanyDao.Get(ctx, tx, companyId)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return NewCcError(cons.CodeComInfoNotExist, ErrCompany, ErrNotFound, ErrNull, "the company is not exist")
	default:
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//update info
	//params["UpdatedAt"] = time.Now()
	err = cs.CompanyDao.Update(ctx, tx, companyId, params)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		ps.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
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
			case "companyId", "companyName", "abbreviationName", "corporator":
				filterFields[*f.Field] = f.Value
			case "phone", "e_mail", "companyAddr", "backup":
				filterFields[*f.Field] = f.Value
			default:
				return comViewSlice, 0, NewError(ErrDesc, ErrUnsupported, ErrField, *f.Field)
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
		vs.Logger.ErrorContext(ctx, "[CompanyService/service/ListCompany] [CompanyDao.List: %s, filterFields: %v]",
			err.Error(), filterFields)
		return comViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	for _, comInfo := range comInfos {
		comInfoView := cs.AccSubMdelToView(comInfo)
		comViewSlice = append(comViewSlice, comInfoView)
	}
	comInfoCount := len(comViewSlice)
	//volumeCount, CcErr := vs.CountByFilter(ctx, vs.Db, filterFields)
	if CcErr != nil {
		return nil, 0, CcErr
	}
	return comViewSlice, comInfoCount, nil
}
