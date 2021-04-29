package service

import (
	"context"
	"database/sql"

	"analysis-server/api/db"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
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
			case "voucher_id", "company_id", "voucher_month", "num_of_month", "voucher_date":
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

func VoucherInfoModelToView(vInfo *model.VoucherInfo) *model.VoucherInfoView {
	vInfoView := new(model.VoucherInfoView)
	vInfoView.VoucherID = vInfo.VoucherID
	vInfoView.CompanyID = vInfo.CompanyID
	vInfoView.VoucherDate = vInfo.VoucherDate
	vInfoView.VoucherMonth = vInfo.VoucherMonth
	vInfoView.NumOfMonth = vInfo.NumOfMonth
	return vInfoView
}
