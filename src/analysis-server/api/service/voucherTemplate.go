package service

import (
	"context"
	"database/sql"
	"financeMgr/src/analysis-server/api/db"
	"financeMgr/src/analysis-server/model"
	cons "financeMgr/src/common/constant"
	"fmt"
	"time"
)

type VoucherTemplateService struct {
	VTemplateDao *db.VoucherTemplateDao
}

func (vs *VoucherTemplateService) CreateVoucherTemplate(ctx context.Context, params *model.VoucherTemplateParams,
	requestId string) (int, CcError) {
	FuncName := "VoucherTemplateService/service/CreateVoucherTemplate"
	bIsRollBack := true
	// Begin transaction
	tx, err := gDb.Begin()
	if err != nil {
		gLogger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return 0, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, FuncName, tx)
		}
	}()

	vTemplate := new(model.VoucherTemplate)
	vTemplate.CompanyID = *params.CompanyID
	vTemplate.VoucherTemplateID = gIdInfoService.genvVouTempIdInfo.GetNextId()
	vTemplate.RefVoucherID = *params.RefVoucherID
	vTemplate.VoucherYear = *params.VoucherYear
	vTemplate.Illustration = *params.Illustration
	vTemplate.CreatedAt = time.Now()
	if err = vs.VTemplateDao.Create(ctx, tx, vTemplate); err != nil {
		gLogger.ErrorContext(ctx, "[%s] [VTemplateDao.Create: %s]", FuncName, err.Error())
		return 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil && IsDuplicateKeyError(err) {
		gLogger.ErrorContext(ctx, "[%s] [Commit Err: duplicate key conflict]", FuncName)
		return 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	} else if err != nil {
		gLogger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	gLogger.InfoContext(ctx, "CreateVoucherTemplate method end ")
	return vTemplate.VoucherTemplateID, nil
}

func (vs *VoucherTemplateService) DeleteVoucherTemplate(ctx context.Context, voucherTemplateID int,
	requestId string) CcError {
	gLogger.InfoContext(ctx, "DeleteVoucherTemplate method begin")
	//delete voucher template
	err := vs.VTemplateDao.Delete(ctx, gDb, voucherTemplateID)
	if err != nil {
		errMsg := fmt.Sprintf("Delete voucher template failed,errInfo:%s", err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, errMsg)
	}
	gLogger.InfoContext(ctx, "DeleteVoucherTemplate method end")
	return nil
}

func (vs *VoucherTemplateService) GetVoucherTemplate(ctx context.Context, voucherTemplateID int,
	requestId string) (*model.VoucherTemplateView, CcError) {
	vTemplate, err := vs.VTemplateDao.Get(ctx, gDb, voucherTemplateID)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, NewCcError(cons.CodeVoucherTemplateNotExist, ErrVoucherTemplate, ErrNotFound, ErrNull, "the VoucherTemplate is not exist")
	default:
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	vTemplateView := VoucherTemplateModelToView(vTemplate)
	return vTemplateView, nil
}

func (vs *VoucherTemplateService) ListVoucherTemplate(ctx context.Context, params *model.ListParams) ([]*model.VoucherTemplateView, int, CcError) {
	voucherTemplateSlice := make([]*model.VoucherTemplateView, 0)
	filterFields := make(map[string]interface{})
	limit, offset := -1, 0
	if params.Filter != nil {
		for _, f := range params.Filter {
			switch *f.Field {
			case "voucherTemplateId", "companyId", "referenceVoucherId", "voucherYear", "illustration", "status":
				filterFields[*f.Field] = f.Value
			default:
				return voucherTemplateSlice, 0, NewError(ErrVoucherTemplate, ErrUnsupported, ErrField, *f.Field)
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
	voucherTemps, err := vs.VTemplateDao.SimpleList(ctx, gDb, filterFields, limit, offset, orderDirection, orderField)
	if err != nil {
		gLogger.ErrorContext(ctx, "[VoucherTemplateService/service/ListVoucherTemplate] [VTemplateDao.SimpleList: %s, filterFields: %v]", err.Error(), filterFields)
		return voucherTemplateSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	for _, voucherTemp := range voucherTemps {
		vouTempView := VoucherTemplateModelToView(voucherTemp)
		voucherTemplateSlice = append(voucherTemplateSlice, vouTempView)
	}
	vouRecordCount := len(voucherTemps)
	return voucherTemplateSlice, vouRecordCount, nil
}

func VoucherTemplateModelToView(vTemp *model.VoucherTemplate) *model.VoucherTemplateView {
	vTempView := new(model.VoucherTemplateView)
	vTempView.VoucherTemplateID = vTemp.VoucherTemplateID
	vTempView.RefVoucherID = vTemp.RefVoucherID
	vTempView.VoucherYear = vTemp.VoucherYear
	vTempView.Illustration = vTemp.Illustration
	vTempView.CreatedAt = vTemp.CreatedAt
	return vTempView
}
