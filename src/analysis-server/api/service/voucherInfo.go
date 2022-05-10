package service

import (
	"context"
	"database/sql"

	"analysis-server/api/db"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
	"time"
)

type VoucherInfoService struct {
	Logger   *log.Logger
	VInfoDao *db.VoucherInfoDao
	Db       *sql.DB
}

func (vs *VoucherInfoService) GetVoucherInfoByID(ctx context.Context, voucherID, iYear int,
	requestId string) (*model.VoucherInfoView, CcError) {
	vInfo, err := vs.VInfoDao.Get(ctx, vs.Db, voucherID, iYear)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, NewCcError(cons.CodeVoucherInfoNotExist, ErrVoucherInfo, ErrNotFound, ErrNull, "the VoucherInfo is not exist")
	default:
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	vInfoView := VoucherInfoModelToView(vInfo)
	return vInfoView, nil
}

func (vs *VoucherInfoService) ListVoucherInfo(ctx context.Context, params *model.ListParams) ([]*model.VoucherInfoView, int, CcError) {
	vouInfoViewSlice := make([]*model.VoucherInfoView, 0)
	filterNo := make(map[string]interface{})
	filterFields := make(map[string]interface{})
	intervalFilterFields := make(map[string]interface{})
	limit, offset := -1, 0
	iVoucherYear := 0
	if params.Filter != nil {
		for _, f := range params.Filter {
			switch *f.Field {
			case "voucherYear":
				{
					switch f.Value.(type) {
					case float64:
						//从客户端发过来的，解析json时，会解析成float64
						iVoucherYear = int(f.Value.(float64))
					case int:
						//从cli发过来的，解析json时，会解析成int
						iVoucherYear = f.Value.(int)
					}
				}
			case "voucherId", "companyId", "voucherMonth", "numOfMonth", "voucherDate":
				fallthrough
			case "voucherAuditor", "voucherFiller", "status", "billCount":
				filterFields[*f.Field] = f.Value
			case "status_no":
				filterNo["status"] = f.Value
			case "numOfMonth_interval":
				intervalFilterFields["numOfMonth"] = f.Value
			case "voucherDate_interval":
				intervalFilterFields["voucherDate"] = f.Value
			default:
				return vouInfoViewSlice, 0, NewError(ErrVoucherInfo, ErrUnsupported, ErrField, *f.Field)
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
	voucherInfos, err := vs.VInfoDao.List(ctx, vs.Db, filterNo, filterFields, intervalFilterFields, iVoucherYear,
		limit, offset, orderDirection, orderField)
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[VoucherInfoService/service/ListVoucherInfo] [VInfoDao.List: %s, filterFields: %v]", err.Error(), filterFields)
		return vouInfoViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	for _, vouInfo := range voucherInfos {
		vouInfoView := VoucherInfoModelToView(vouInfo)
		vouInfoViewSlice = append(vouInfoViewSlice, vouInfoView)
	}
	vouRecordCount := len(voucherInfos)
	return vouInfoViewSlice, vouRecordCount, nil
}

func (vs *VoucherInfoService) GetLatestVoucherInfoByCompanyID(ctx context.Context, iYear, iCompanyID int,
	requestId string) ([]*model.VoucherInfoView, int, CcError) {
	vouInfoViewSlice := make([]*model.VoucherInfoView, 0)
	voucherInfos, err := vs.VInfoDao.GetLatestVoucherInfoByCompanyID(ctx, vs.Db, iYear, iCompanyID)
	if err != nil {
		FunctionName := "VoucherInfoService/service/GetLatestVoucherInfo"
		vs.Logger.ErrorContext(ctx, "[%s] [Error: %s, companyID: %d]", FunctionName, err.Error(), iCompanyID)
		return vouInfoViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	for _, vouInfo := range voucherInfos {
		vouInfoView := VoucherInfoModelToView(vouInfo)
		vouInfoViewSlice = append(vouInfoViewSlice, vouInfoView)
	}
	vouRecordCount := len(voucherInfos)
	return vouInfoViewSlice, vouRecordCount, nil
}

func VoucherInfoModelToView(vInfo *model.VoucherInfo) *model.VoucherInfoView {
	vInfoView := new(model.VoucherInfoView)
	vInfoView.VoucherID = vInfo.VoucherID
	vInfoView.CompanyID = vInfo.CompanyID
	vInfoView.VoucherDate = vInfo.VoucherDate
	vInfoView.VoucherMonth = vInfo.VoucherMonth
	vInfoView.NumOfMonth = vInfo.NumOfMonth
	vInfoView.VoucherFiller = vInfo.VoucherFiller
	vInfoView.VoucherAuditor = vInfo.VoucherAuditor
	vInfoView.BillCount = vInfo.BillCount
	vInfoView.Status = vInfo.Status
	return vInfoView
}

//该函数也可以用于审核、取消审核、作废凭证等功能 ...
func (vs *VoucherInfoService) UpdateVoucherInfoByID(ctx context.Context, voucherID, iYear int,
	params map[string]interface{}) CcError {
	FuncName := "VoucherInfoService/UpdateVoucherInfoByID"
	bIsRollBack := true
	tx, err := vs.Db.Begin()
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, vs.Logger, FuncName, tx)
		}
	}()
	//insure the voucherInfo exist
	_, err = vs.VInfoDao.Get(ctx, tx, voucherID, iYear)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return NewCcError(cons.CodeVoucherInfoNotExist, ErrVoucherInfo, ErrNotFound, ErrNull, "the VoucherInfo is not exist")
	default:
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	params["updatedAt"] = time.Now()
	err = vs.VInfoDao.Update(ctx, tx, voucherID, iYear, params)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	return nil
}

func (vs *VoucherInfoService) BatchAuditVoucherInfo(ctx context.Context, params *model.BatchAuditParams) CcError {
	FuncName := "VoucherInfoService/BatchAuditVoucherInfo"
	bIsRollBack := true
	tx, err := vs.Db.Begin()
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, vs.Logger, FuncName, tx)
		}
	}()
	//insure the voucherInfo exist  由于是批量更新，所以就不一一判断了。
	err = vs.VInfoDao.BatchUpdate(ctx, tx, *params.VoucherYear, *params.Status, *params.VoucherAuditor, params.IDs)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	return nil
}

func (vs *VoucherInfoService) GetMaxNumOfMonthByContion(ctx context.Context,
	params *model.QueryMaxNumOfMonthParams, requestId string) (int64, CcError) {
	filterFields := make(map[string]interface{})
	filterFields["companyId"] = *params.CompanyID
	filterFields["voucherMonth"] = *params.VoucherMonth
	count, err := vs.VInfoDao.CountByFilter(ctx, vs.Db, *params.VoucherYear, filterFields)
	if err != nil {
		return 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	return count, nil
}
