package service

import (
	"context"
	"database/sql"
	"time"

	"analysis-server/api/db"
	"analysis-server/api/utils"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
	"github.com/go-sql-driver/mysql"
)

type VoucherService struct {
	Logger       *log.Logger
	VInfoDao     *db.VoucherInfoDao
	VRecordDao   *db.VoucherRecordDao
	Db           *sql.DB
	GenVoucherId *utils.GenIdInfo
	GenRecordId  *utils.GenIdInfo
}

func IsDuplicateKeyError(err error) bool {
	if mysqlErr := err.(*mysql.MySQLError); mysqlErr != nil {
		return mysqlErr.Number == 1062
	}
	return false
}

func (vs *VoucherService) CreateVoucher(ctx context.Context, params *model.VoucherParams,
	requestId string) ([]int, CcError) {
	FuncName := "VoucherService/service/CreateVoucher"
	// Begin transaction
	tx, err := vs.Db.Begin()
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer RollbackLog(ctx, vs.Logger, FuncName, tx)
	infoParams := params.InfoParams
	filterFields := make(map[string]interface{})
	filterFields["companyId"] = *infoParams.CompanyID
	filterFields["voucherMonth"] = *infoParams.VoucherMonth
	vs.Logger.InfoContext(ctx, "CreateVoucher method start, "+"companyID:%d,VoucherMonth:%d", *infoParams.CompanyID, *infoParams.VoucherMonth)
	count, err := vs.VInfoDao.CountByFilter(ctx, tx, filterFields)
	if err != nil {
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//第0个元素是voucherId,从第1个元素开始，就是recordId
	var IdValSli []int
	vInfo := new(model.VoucherInfo)
	vInfo.CompanyID = *infoParams.CompanyID
	vInfo.VoucherMonth = *infoParams.VoucherMonth
	vInfo.NumOfMonth = int(count + 1)
	vInfo.VoucherDate = time.Now()
	vInfo.CreatedAt = time.Now()
	vInfo.VoucherID = vs.GenVoucherId.GetNextId()
	IdValSli = append(IdValSli, vInfo.VoucherID)
	if err = vs.VInfoDao.Create(ctx, tx, vInfo); err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [VInfoDao.Create: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//create voucherRecord
	vRecord := new(model.VoucherRecord)
	for _, recParam := range params.RecordsParams {
		vRecord.RecordID = vs.GenRecordId.GetNextId()
		IdValSli = append(IdValSli, vRecord.RecordID)
		vRecord.VoucherID = vInfo.VoucherID
		vRecord.SubjectName = *recParam.SubjectName
		vRecord.DebitMoney = *recParam.DebitMoney
		vRecord.CreditMoney = *recParam.CreditMoney
		vRecord.Summary = *recParam.Summary
		vRecord.BillCount = *recParam.BillCount
		vRecord.SubID1 = *recParam.SubID1
		vRecord.SubID2 = *recParam.SubID2
		vRecord.SubID3 = *recParam.SubID3
		vRecord.SubID4 = *recParam.SubID4
		vRecord.CreatedAt = time.Now()
		if err = vs.VRecordDao.Create(ctx, tx, vRecord); err != nil {
			vs.Logger.ErrorContext(ctx, "[%s] [VRecordDao.Create: %s]", FuncName, err.Error())
			return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
	}
	if err = tx.Commit(); err != nil && IsDuplicateKeyError(err) {
		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: duplicate key conflict]", FuncName)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	} else if err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	vs.Logger.InfoContext(ctx, "CreateVoucher method end ")
	return IdValSli, nil
}

// func (vs *VoucherService) VoucherInfoModelToView(vInfo *model.VoucherInfo) *model.VoucherInfoView {
// 	vInfoView := new(model.VoucherInfoView)
// 	vInfoView.VoucherID = vInfo.VoucherID
// 	vInfoView.CompanyID = vInfo.CompanyID
// 	vInfoView.VoucherDate = vInfo.VoucherDate
// 	vInfoView.VoucherMonth = vInfo.VoucherMonth
// 	vInfoView.NumOfMonth = vInfo.NumOfMonth
// 	return vInfoView
// }

func (vs *VoucherService) DeleteVoucher(ctx context.Context, voucherID int, requestId string) CcError {
	FuncName := "VoucherService/service/DeleteVoucher"
	vs.Logger.InfoContext(ctx, "DeleteVoucher method begin, "+"voucher ID:%d", voucherID)
	// Begin transaction
	tx, err := vs.Db.Begin()
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer RollbackLog(ctx, vs.Logger, FuncName, tx)

	err = vs.VInfoDao.Delete(ctx, tx, voucherID)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, "Delete failed")
	}

	err = vs.VRecordDao.DeleteByVoucherId(ctx, vs.Db, voucherID)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, "Delete failed")
	}
	if err = tx.Commit(); err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	vs.Logger.InfoContext(ctx, "DeleteVoucher method end, "+"voucher ID:%d", voucherID)
	return nil
}

func (vs *VoucherService) GetVoucherByVoucherID(ctx context.Context, voucherID int,
	requestId string) (*model.VoucherView, CcError) {
	FuncName := "VoucherService/service/GetVoucherByVoucherID"
	vs.Logger.InfoContext(ctx, "GetVoucherByVoucherID method begin, "+"voucher ID:%d", voucherID)
	//Begin transaction
	tx, err := vs.Db.Begin()
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer RollbackLog(ctx, vs.Logger, FuncName, tx)

	vInfo, err := vs.VInfoDao.Get(ctx, tx, voucherID)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, NewCcError(cons.CodeVoucherInfoNotExist, ErrVoucherInfo, ErrNotFound, ErrNull, "the VoucherInfo is not exist")
	default:
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	voucherView := &model.VoucherView{}
	voucherView.VouInfoView = *(VoucherInfoModelToView(vInfo))
	recordViewSlice := make([]model.VoucherRecordView, 0)
	filterFields := make(map[string]interface{})
	limit, offset := -1, 0
	filterFields["voucherId"] = voucherID
	orderField := "recordId"
	orderDirection := cons.Order_Asc

	voucherRecords, err := vs.VRecordDao.List(ctx, tx, filterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[VoucherService/service/GetVoucherByVoucherID] [VRecordDao.List: %s, filterFields: %v]", err.Error(), filterFields)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//默认最多返回100条记录，如果记录超过100条，需要在客户端再把剩余的部分获取出来。
	limit = 100
	for index, vouRecord := range voucherRecords {
		if index >= 100 {
			break
		}
		vouRecordView := *(VoucherRecordModelToView(vouRecord))
		recordViewSlice = append(recordViewSlice, vouRecordView)
	}
	voucherView.VouRecordTotalCount = len(voucherRecords)
	voucherView.VouRecordViewSli = append(voucherView.VouRecordViewSli, recordViewSlice...)

	if err = tx.Commit(); err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	vs.Logger.InfoContext(ctx, "GetVoucherByVoucherID method end, "+"voucher ID:%d", voucherID)
	return voucherView, nil
}

// func (vs *VoucherService) UpdateVoucherInfo(ctx context.Context, voucherID int, params map[string]interface{}) service.CcError {
// 	//insure the volume exist
// 	volume, err := vs.VInfoDao.Get(ctx, vs.Db, voucherID)
// 	switch err {
// 	case nil:
// 	case sql.ErrNoRows:
// 		return NewCcError(cons.CodeVoucherInfoNotExist, ErrVoucherInfo, ErrNotFound, ErrNull, "the VoucherInfo is not exist")
// 	default:
// 		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
// 	}
// 	//update info
// 	// params["DeletedAt"] = time.Now()
// 	// err = vs.VInfoDao.Update(ctx, vs.Db, voucherID, params)
// 	// if err != nil {
// 	// 	return service.NewError(service.ErrSystem, service.ErrError, service.ErrNull, err.Error())
// 	// }
// 	return nil
// }
