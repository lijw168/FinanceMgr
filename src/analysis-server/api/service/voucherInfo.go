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

func (vs *VoucherInfoService) GetVoucherInfoByID(ctx context.Context, voucherID int,
	requestId string) (*model.VoucherInfoView, CcError) {
	vInfo, err := vs.VInfoDao.Get(ctx, vs.Db, voucherID)
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
	filterFields := make(map[string]interface{})
	limit, offset := -1, 0
	if params.Filter != nil {
		for _, f := range params.Filter {
			switch *f.Field {
			case "voucherId", "companyId", "voucherMonth", "numOfMonth", "voucherDate", "voucherFiller", "voucherAuditor":
				filterFields[*f.Field] = f.Value
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
	voucherInfos, err := vs.VInfoDao.List(ctx, vs.Db, filterFields, limit, offset, orderField, orderDirection)
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

func (vs *VoucherInfoService) GetLatestVoucherInfoByCompanyID(ctx context.Context, iCompanyID int,
	requestId string) ([]*model.VoucherInfoView, int, CcError) {
	vouInfoViewSlice := make([]*model.VoucherInfoView, 0)
	voucherInfos, err := vs.VInfoDao.GetLatestVoucherInfoByCompanyID(ctx, vs.Db, iCompanyID)
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
	return vInfoView
}

func (vs *VoucherInfoService) UpdateVoucherInfo(ctx context.Context, voucherID int,
	params map[string]interface{}) CcError {
	FuncName := "VoucherInfoService/UpdateVoucherInfo"
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
	_, err = vs.VInfoDao.Get(ctx, tx, voucherID)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return NewCcError(cons.CodeVoucherInfoNotExist, ErrVoucherInfo, ErrNotFound, ErrNull, "the VoucherInfo is not exist")
	default:
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	params["updatedAt"] = time.Now()
	err = vs.VInfoDao.Update(ctx, tx, voucherID, params)
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
