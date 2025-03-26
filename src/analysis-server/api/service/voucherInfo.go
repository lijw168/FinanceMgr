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

	iVoucherYear := 0
	if params.Filter != nil {
		var intervalValSlice []int
		for _, f := range params.Filter {
			if *f.Field == "numOfMonth_interval" || *f.Field == "voucherDate_interval" ||
				*f.Field == "voucherMonth_interval" {
				err := FormatData(f.Value, &intervalValSlice)
				if err != nil {
					return vouInfoViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
				}
			}
			switch *f.Field {
			case "voucherYear":
				{
					switch f.Value.(type) {
					case float64:
						//从客户端发过来的，解析json时，会解析成float64 (经验证该结论是错误的)
						//正确的结论是，文档显示当把json解析成interface{}时，把number解析成float64
						iVoucherYear = int(f.Value.(float64))
						//测试代码
						//vs.Logger.ErrorContext(ctx, "the iVoucherYear is float64")
					}
				}
			case "voucherId", "companyId", "voucherMonth", "numOfMonth", "voucherDate":
				fallthrough
			case "voucherAuditor", "voucherFiller", "status", "billCount":
				filterFields[*f.Field] = f.Value
			case "status_no":
				filterNo["status"] = f.Value
			case "numOfMonth_interval":
				intervalFilterFields["numOfMonth"] = intervalValSlice
				intervalValSlice = intervalValSlice[0:0]
			case "voucherDate_interval":
				intervalFilterFields["voucherDate"] = intervalValSlice
				intervalValSlice = intervalValSlice[0:0]
			case "voucherMonth_interval":
				intervalFilterFields["voucherMonth"] = intervalValSlice
				intervalValSlice = intervalValSlice[0:0]
			default:
				return vouInfoViewSlice, 0, NewError(ErrVoucherInfo, ErrUnsupported, ErrField, *f.Field)
			}
		}
	}
	limit, offset := -1, 0
	if params.DescLimit != nil {
		limit = *params.DescLimit
		if params.DescOffset != nil {
			offset = *params.DescOffset
		}
	}
	// orderFilter := make(map[string]int)
	// for _, v := range params.Order {
	// 	orderFilter[*v.Field] = *v.Direction
	// }
	voucherInfos, err := vs.VInfoDao.List(ctx, vs.Db, filterNo, filterFields, intervalFilterFields,
		nil, params.Order, iVoucherYear, limit, offset)
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

func (vs *VoucherInfoService) GetLatestVoucherInfoByCompanyID(ctx context.Context, iMonth, iYear, iCompanyID int,
	requestId string) ([]*model.VoucherInfoView, int, CcError) {
	vouInfoViewSlice := make([]*model.VoucherInfoView, 0)
	voucherInfos, err := vs.VInfoDao.GetLatestVoucherInfo(ctx, vs.Db, iMonth, iYear, iCompanyID)
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

// 该函数也可以用于审核、取消审核、作废凭证等功能 ...
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
	params *model.QueryMonthlyVoucherInfoAttrParameters, requestId string) (int64, CcError) {
	filterFields := make(map[string]interface{})
	filterFields["companyId"] = *params.CompanyID
	filterFields["voucherMonth"] = *params.VoucherMonth
	count, err := vs.VInfoDao.CountByFilter(ctx, vs.Db, *params.VoucherYear, filterFields)
	if err != nil {
		return 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	return count, nil
}

func (vs *VoucherInfoService) GetNoAuditedVoucherInfoCountByContion(ctx context.Context,
	params *model.QueryVoucherInfoStatusParams, requestId string) (int64, CcError) {
	filterFields := make(map[string]interface{})
	filterFields["companyId"] = *params.CompanyID
	filterFields["status"] = *params.Status
	count, err := vs.VInfoDao.CountByFilter(ctx, vs.Db, *params.VoucherYear, filterFields)
	if err != nil {
		return 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	return count, nil
}
