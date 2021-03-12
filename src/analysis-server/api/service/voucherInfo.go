package service

import (
	"context"
	"database/sql"

	cons "common/constant"
	"common/log"
	"analysis-server/api/db"
	"analysis-server/model"
)

type VoucherInfoService struct {
	Logger   *log.Logger
	VInfoDao *db.VoucherInfoDao
	Db       *sql.DB
}

func VoucherInfoModelToView(vInfo *model.VoucherInfo) *model.VoucherInfoView {
	vInfoView := new(model.VoucherInfoView)
	vInfoView.VoucherID = vInfo.VoucherID
	vInfoView.CompanyID = vInfo.CompanyID
	vInfoView.VoucherDate = vInfo.VoucherDate
	vInfoView.VoucherMonth = vInfo.VoucherMonth
	vInfoView.NumOfMonth = vInfo.NumOfMonth
	return vInfoView
}

func (vs *VoucherInfoService) GetVoucherInfoByID(ctx context.Context, voucherID int,
	requestId string) (*model.VoucherInfoView, CcError) {
	vInfo, err := vs.VInfoDao.Get(ctx, vs.Db, voucherID)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return NewCcError(cons.CodeVoucherInfoNotExist, ErrVoucherInfo, ErrNotFound, ErrNull, "the VoucherInfo is not exist")
	default:
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
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
			case "voucherId", "companyId", "voucherMonth", "numOfMonth", "voucherDate":
				filterFields[*f.Field] = f.Value
			default:
				return vouInfoViewSlice, 0, NewError(ErrDesc, ErrUnsupported, ErrField, *f.Field)
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
	voucherInfos, err := vs.VInfoDao.List(ctx, ps.Db, filterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[VoucherInfoService/service/ListVoucherInfo] [VInfoDao.List: %s, filterFields: %v]", err.Error(), filterFields)
		return vouInfoViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	for _, vouInfo := range voucherInfos {
		vouInfoView := VoucherInfoModelToView(vouInfo)
		vouInfoViewSlice = append(vouInfoViewSlice, vouInfoView)
	}
	vouRecordCount := len(voucherInfos)
	//volumeCount, CcErr := vs.CountByFilter(ctx, vs.Db, filterFields)
	// if CcErr != nil {
	// 	return nil, 0, CcErr
	// }
	return vouInfoViewSlice, vouRecordCount, nil
}
