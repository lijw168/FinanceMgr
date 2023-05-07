package service

import (
	"context"
	"database/sql"

	"financeMgr/src/analysis-server/api/db"
	"financeMgr/src/analysis-server/model"
	"financeMgr/src/common/log"
)

type VoucherRecordService struct {
	Logger     *log.Logger
	VRecordDao *db.VoucherRecordDao
	Db         *sql.DB
	//GenRecordId *utils.GenIdInfo
}

// func (vs *VoucherRecordService) CreateVoucherRecord(ctx context.Context, params *model.CreateVoucherRecordParams,
// 	requestId string) (int, CcError) {
// 	//create
// 	vs.Logger.InfoContext(ctx, "CreateVoucherRecord method start, "+"VoucherId:%d", *params.VoucherID)

// 	FuncName := "VoucherRecordService/CreateVoucherRecord"
// 	vRecord := new(model.VoucherRecord)
// 	vRecord.RecordID = GIdInfoService.genVouRecIdInfo.GetNextId()
// 	vRecord.VoucherID = *params.VoucherID
// 	vRecord.SubjectName = *params.SubjectName
// 	vRecord.DebitMoney = *params.DebitMoney
// 	vRecord.CreditMoney = *params.CreditMoney
// 	vRecord.Summary = *params.Summary
// 	vRecord.SubID1 = *params.SubID1
// 	vRecord.SubID2 = *params.SubID2
// 	vRecord.SubID3 = *params.SubID3
// 	vRecord.SubID4 = *params.SubID4
// 	vRecord.CreatedAt = time.Now()

// 	if err := vs.VRecordDao.Create(ctx, vs.Db, vRecord); err != nil {
// 		vs.Logger.ErrorContext(ctx, "[%s] [VRecordDao.Create: %s]", FuncName, err.Error())
// 		return -1, NewError(ErrSystem, ErrError, ErrNull, err.Error())
// 	}
// 	//vRecordView := vs.VoucherRecordModelToView(vRecord)
// 	vs.Logger.InfoContext(ctx, "CreateVoucherRecord method end ")
// 	return vRecord.RecordID, nil
// }

// func (vs *VoucherRecordService) CreateVoucherRecords(ctx context.Context, recordsParams []*model.CreateVoucherRecordParams,
// 	requestId string) ([]int, CcError) {
// 	//create
// 	vs.Logger.InfoContext(ctx, "CreateVoucherRecords method start, "+"requestId:%s", requestId)
// 	FuncName := "VoucherRecordService/CreateVoucherRecords"
// 	bIsRollBack := true
// 	tx, err := vs.Db.Begin()
// 	if err != nil {
// 		vs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
// 		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
// 	}
// 	defer func() {
// 		if bIsRollBack {
// 			RollbackLog(ctx, vs.Logger, FuncName, tx)
// 		}
// 	}()

// 	var IdValSli []int
// 	for _, itemParam := range recordsParams {
// 		vRecord := new(model.VoucherRecord)
// 		vRecord.RecordID = GIdInfoService.genVouRecIdInfo.GetNextId()
// 		IdValSli = append(IdValSli, vRecord.RecordID)
// 		vRecord.VoucherID = *itemParam.VoucherID
// 		vRecord.SubjectName = *itemParam.SubjectName
// 		vRecord.DebitMoney = *itemParam.DebitMoney
// 		vRecord.CreditMoney = *itemParam.CreditMoney
// 		vRecord.Summary = *itemParam.Summary
// 		vRecord.SubID1 = *itemParam.SubID1
// 		// vRecord.SubID2 = *itemParam.SubID2
// 		// vRecord.SubID3 = *itemParam.SubID3
// 		// vRecord.SubID4 = *itemParam.SubID4
// 		vRecord.CreatedAt = time.Now()
// 		if err = vs.VRecordDao.Create(ctx, tx, vRecord); err != nil {
// 			vs.Logger.ErrorContext(ctx, "[%s] [VRecordDao.Create: %s]", FuncName, err.Error())
// 			return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
// 		}
// 	}
// 	if err = tx.Commit(); err != nil && IsDuplicateKeyError(err) {
// 		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: duplicate key conflict]", FuncName)
// 		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
// 	} else if err != nil {
// 		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
// 		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
// 	}
// 	bIsRollBack = false
// 	//vRecordView := vs.VoucherRecordModelToView(vRecord)
// 	vs.Logger.InfoContext(ctx, "CreateVoucherRecords method end ")
// 	return IdValSli, nil
// }

func VoucherRecordModelToView(vRecord *model.VoucherRecord) *model.VoucherRecordView {
	vRecordView := new(model.VoucherRecordView)
	vRecordView.VoucherID = vRecord.VoucherID
	vRecordView.RecordID = vRecord.RecordID
	vRecordView.SubjectName = vRecord.SubjectName
	vRecordView.DebitMoney = vRecord.DebitMoney
	vRecordView.CreditMoney = vRecord.CreditMoney
	vRecordView.Summary = vRecord.Summary
	vRecordView.SubID1 = vRecord.SubID1
	vRecordView.SubID2 = vRecord.SubID2
	vRecordView.SubID3 = vRecord.SubID3
	vRecordView.SubID4 = vRecord.SubID4
	return vRecordView
}

func (vs *VoucherRecordService) ListVoucherRecords(ctx context.Context,
	params *model.ListParams) ([]*model.VoucherRecordView, int, CcError) {
	recordViewSlice := make([]*model.VoucherRecordView, 0)
	filterNo := make(map[string]interface{})
	filterFields := make(map[string]interface{})
	intervalFilterFields := make(map[string]interface{})
	fuzzyMatchFields := make(map[string]string)
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
			case "recordId", "voucherId", "subjectName", "summary", "subId1":
				fallthrough
			case "subId2", "subId3", "subId4":
				filterFields[*f.Field] = f.Value
			case "debitMoney_interval":
				intervalFilterFields["debitMoney"] = f.Value
			case "creditMoney_interval":
				intervalFilterFields["creditMoney"] = f.Value
			case "subjectName_fuzzy":
				fuzzyMatchFields["subjectName"] = f.Value.(string)
			case "summary_fuzzy":
				fuzzyMatchFields["summary"] = f.Value.(string)
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
	voucherRecords, err := vs.VRecordDao.List(ctx, vs.Db, filterNo, filterFields, intervalFilterFields, fuzzyMatchFields,
		iVoucherYear, limit, offset, orderDirection, orderField)
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

// func (vs *VoucherRecordService) DeleteVoucherRecordByID(ctx context.Context, recordID int,
// 	requestId string) CcError {
// 	vs.Logger.InfoContext(ctx, "DeleteVoucherRecordByID method begin, "+"record ID:%d", recordID)
// 	err := vs.VRecordDao.Delete(ctx, vs.Db, recordID)
// 	if err != nil {
// 		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
// 	}
// 	vs.Logger.InfoContext(ctx, "DeleteVoucherRecordByID method end.")
// 	return nil
// }

// func (vs *VoucherRecordService) DeleteVoucherRecords(ctx context.Context, params *model.IDsParams,
// 	requestId string) CcError {
// 	vs.Logger.InfoContext(ctx, "DeleteVoucherRecords method begin, "+"record IDs:%v", params.IDs)
// 	//var vouIds = []int{}
// 	delConditonParams := make(map[string]interface{})
// 	delConditonParams["recordId"] = params.IDs
// 	err := vs.VRecordDao.DeleteByMultiCondition(ctx, vs.Db, delConditonParams)
// 	if err != nil {
// 		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
// 	}
// 	vs.Logger.InfoContext(ctx, "DeleteVoucherRecords method end.")
// 	return nil
// }

// func (vs *VoucherRecordService) GetVoucherRecordByID(ctx context.Context, recordID int,
// 	requestId string) (*model.VoucherRecordView, CcError) {
// 	vRecord, err := vs.VRecordDao.Get(ctx, vs.Db, recordID)
// 	switch err {
// 	case nil:
// 	case sql.ErrNoRows:
// 		return nil, NewCcError(cons.CodeVoucherRecordNotExist, ErrVoucher, ErrNotFound, ErrNull, "the voucher record is not exist")
// 	default:
// 		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
// 	}
// 	vRecordView := VoucherRecordModelToView(vRecord)
// 	return vRecordView, nil
// }

// func (vs *VoucherRecordService) UpdateVoucherRecordByID(ctx context.Context, recordID int,
// 	params map[string]interface{}) CcError {
// 	FuncName := "VoucherRecordService/UpdateVoucherRecordByID"
// 	bIsRollBack := true
// 	tx, err := vs.Db.Begin()
// 	if err != nil {
// 		vs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
// 		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
// 	}
// 	defer func() {
// 		if bIsRollBack {
// 			RollbackLog(ctx, vs.Logger, FuncName, tx)
// 		}
// 	}()
// 	_, err = vs.VRecordDao.Get(ctx, tx, recordID)
// 	switch err {
// 	case nil:
// 	case sql.ErrNoRows:
// 		return NewCcError(cons.CodeVoucherRecordNotExist, ErrVoucher, ErrNotFound, ErrNull, "the Voucher record is not exist")
// 	default:
// 		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
// 	}
// 	//update voucher record
// 	params["updatedAt"] = time.Now()
// 	err = vs.VRecordDao.UpdateByRecordId(ctx, tx, recordID, params)
// 	if err != nil {
// 		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
// 	}
// 	if err = tx.Commit(); err != nil {
// 		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
// 		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
// 	}
// 	bIsRollBack = false
// 	return nil
// }
