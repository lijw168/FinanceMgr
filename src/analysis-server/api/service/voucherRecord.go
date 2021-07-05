package service

import (
	"context"
	"database/sql"

	"analysis-server/api/db"
	"analysis-server/api/utils"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
	"time"
)

type VoucherRecordService struct {
	Logger      *log.Logger
	VRecordDao  *db.VoucherRecordDao
	Db          *sql.DB
	GenRecordId *utils.GenIdInfo
}

func (vs *VoucherRecordService) CreateVoucherRecord(ctx context.Context, params *model.CreateVoucherRecordParams,
	requestId string) (int, CcError) {
	//create
	vs.Logger.InfoContext(ctx, "CreateVoucherRecord method start, "+"VoucherId:%d", *params.VoucherID)

	FuncName := "VoucherRecordService/CreateVoucherRecord"
	vRecord := new(model.VoucherRecord)
	vRecord.RecordID = vs.GenRecordId.GetNextId()
	vRecord.VoucherID = *params.VoucherID
	vRecord.SubjectName = *params.SubjectName
	vRecord.DebitMoney = *params.DebitMoney
	vRecord.CreditMoney = *params.CreditMoney
	vRecord.Summary = *params.Summary
	vRecord.BillCount = *params.BillCount
	vRecord.SubID1 = *params.SubID1
	vRecord.SubID2 = *params.SubID2
	vRecord.SubID3 = *params.SubID3
	vRecord.SubID4 = *params.SubID4
	vRecord.CreatedAt = time.Now()

	if err := vs.VRecordDao.Create(ctx, vs.Db, vRecord); err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [VRecordDao.Create: %s]", FuncName, err.Error())
		return -1, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//vRecordView := vs.VoucherRecordModelToView(vRecord)
	vs.Logger.InfoContext(ctx, "CreateVoucherRecord method end ")
	return vRecord.RecordID, nil
}

func (vs *VoucherRecordService) CreateVoucherRecords(ctx context.Context, recordsParams []*model.CreateVoucherRecordParams,
	requestId string) ([]int, CcError) {
	//create
	vs.Logger.InfoContext(ctx, "CreateVoucherRecords method start, "+"requestId:%s", requestId)
	FuncName := "VoucherRecordService/CreateVoucherRecords"
	bIsRollBack := true
	tx, err := vs.Db.Begin()
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func(bRollBack bool) {
		if bRollBack {
			RollbackLog(ctx, vs.Logger, FuncName, tx)
		}
	}(bIsRollBack)

	var IdValSli []int
	for _, itemParam := range recordsParams {
		vRecord := new(model.VoucherRecord)
		vRecord.RecordID = vs.GenRecordId.GetNextId()
		IdValSli = append(IdValSli, vRecord.RecordID)
		vRecord.VoucherID = *itemParam.VoucherID
		vRecord.SubjectName = *itemParam.SubjectName
		vRecord.DebitMoney = *itemParam.DebitMoney
		vRecord.CreditMoney = *itemParam.CreditMoney
		vRecord.Summary = *itemParam.Summary
		vRecord.BillCount = *itemParam.BillCount
		vRecord.SubID1 = *itemParam.SubID1
		vRecord.SubID2 = *itemParam.SubID2
		vRecord.SubID3 = *itemParam.SubID3
		vRecord.SubID4 = *itemParam.SubID4
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
	bIsRollBack = false
	//vRecordView := vs.VoucherRecordModelToView(vRecord)
	vs.Logger.InfoContext(ctx, "CreateVoucherRecords method end ")
	return IdValSli, nil
}

func VoucherRecordModelToView(vRecord *model.VoucherRecord) *model.VoucherRecordView {
	vRecordView := new(model.VoucherRecordView)
	vRecordView.VoucherID = vRecord.VoucherID
	vRecordView.RecordID = vRecord.RecordID
	vRecordView.SubjectName = vRecord.SubjectName
	vRecordView.DebitMoney = vRecord.DebitMoney
	vRecordView.CreditMoney = vRecord.CreditMoney
	vRecordView.Summary = vRecord.Summary
	vRecordView.BillCount = vRecord.BillCount
	vRecordView.SubID1 = vRecord.SubID1
	vRecordView.SubID2 = vRecord.SubID2
	vRecordView.SubID3 = vRecord.SubID3
	vRecordView.SubID4 = vRecord.SubID4
	return vRecordView
}

func (vs *VoucherRecordService) ListVoucherRecords(ctx context.Context,
	params *model.ListParams) ([]*model.VoucherRecordView, int, CcError) {
	recordViewSlice := make([]*model.VoucherRecordView, 0)
	filterFields := make(map[string]interface{})
	limit, offset := -1, 0
	if params.Filter != nil {
		for _, f := range params.Filter {
			switch *f.Field {
			case "recordId", "voucherId", "subjectName", "summary", "subId1", "subId2", "subId3", "subId4":
				filterFields[*f.Field] = f.Value
			default:
				return recordViewSlice, 0, NewError(ErrVoucher, ErrUnsupported, ErrField, *f.Field)
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
	voucherRecords, err := vs.VRecordDao.List(ctx, vs.Db, filterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[VoucherRecordService/service/ListVoucherRecords] [VRecordDao.List: %s, filterFields: %v]", err.Error(), filterFields)
		return recordViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	for _, vouRecord := range voucherRecords {
		vouRecordView := VoucherRecordModelToView(vouRecord)
		recordViewSlice = append(recordViewSlice, vouRecordView)
	}
	vouRecordCount := len(voucherRecords)
	return recordViewSlice, vouRecordCount, nil
}

func (vs *VoucherRecordService) DeleteVoucherRecordByID(ctx context.Context, recordID int, requestId string) CcError {
	vs.Logger.InfoContext(ctx, "DeleteVoucherRecordByID method begin, "+"record ID:%d", recordID)
	err := vs.VRecordDao.Delete(ctx, vs.Db, recordID)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, "Delete failed")
	}
	vs.Logger.InfoContext(ctx, "DeleteVoucherRecordByID method end, "+"voucher ID:%d", recordID)
	return nil
}

func (vs *VoucherRecordService) GetVoucherRecordByID(ctx context.Context, recordID int,
	requestId string) (*model.VoucherRecordView, CcError) {
	vRecord, err := vs.VRecordDao.Get(ctx, vs.Db, recordID)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, NewCcError(cons.CodeVoucherRecordNotExist, ErrVoucher, ErrNotFound, ErrNull, "the voucher record is not exist")
	default:
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	vRecordView := VoucherRecordModelToView(vRecord)
	return vRecordView, nil
}

func (vs *VoucherRecordService) UpdateVoucherRecord(ctx context.Context, recordID int, params map[string]interface{}) CcError {
	FuncName := "VoucherRecordService/UpdateVoucherRecord"
	bIsRollBack := true
	tx, err := vs.Db.Begin()
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func(bRollBack bool) {
		if bRollBack {
			RollbackLog(ctx, vs.Logger, FuncName, tx)
		}
	}(bIsRollBack)
	_, err = vs.VRecordDao.Get(ctx, tx, recordID)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return NewCcError(cons.CodeVoucherRecordNotExist, ErrVoucher, ErrNotFound, ErrNull, "the Voucher record is not exist")
	default:
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//update voucher record
	params["UpdatedAt"] = time.Now()
	err = vs.VRecordDao.Update(ctx, tx, recordID, params)
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
