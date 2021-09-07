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

//默认最多返回100条记录，如果记录超过100条，需要在客户端再把剩余的部分获取出来。
const (
	MaxRecordLimit = 100
)

type VoucherService struct {
	Logger     *log.Logger
	VInfoDao   *db.VoucherInfoDao
	VRecordDao *db.VoucherRecordDao
	Db         *sql.DB
	//GenVoucherId *utils.GenIdInfo
	//GenRecordId  *utils.GenIdInfo
}

func IsDuplicateKeyError(err error) bool {
	if mysqlErr := err.(*mysql.MySQLError); mysqlErr != nil {
		return mysqlErr.Number == 1062
	}
	return false
}

func (vs *VoucherService) CreateVoucher(ctx context.Context, params *model.CreateVoucherParams,
	requestId string) ([]int, CcError) {
	FuncName := "VoucherService/service/CreateVoucher"
	bIsRollBack := true
	// Begin transaction
	tx, err := vs.Db.Begin()
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, vs.Logger, FuncName, tx)
		}
	}()

	vInfo := new(model.VoucherInfo)
	infoParams := params.InfoParams
	if infoParams.VoucherDate != nil {
		iDate := *infoParams.VoucherDate
		iYear := iDate / 10000
		iMonth := (iDate - iYear*10000) / 100
		iDay := iDate % 100
		t := time.Date(iYear, time.Month(iMonth), iDay, 0, 0, 0, 0, time.Local)
		vInfo.VoucherDate = t
	} else {
		vInfo.VoucherDate = time.Now()
	}
	_, month, _ := vInfo.VoucherDate.Date()
	vInfo.VoucherMonth = int(month)
	filterFields := make(map[string]interface{})
	filterFields["companyId"] = *infoParams.CompanyID
	filterFields["voucherMonth"] = vInfo.VoucherMonth
	vs.Logger.InfoContext(ctx, "CreateVoucher method start, "+"companyID:%d,VoucherMonth:%d", *infoParams.CompanyID, vInfo.VoucherMonth)
	count, err := vs.VInfoDao.CountByFilter(ctx, tx, filterFields)
	if err != nil {
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//第0个元素是voucherId,从第1个元素开始，就是recordId
	var IdValSli []int
	vInfo.CompanyID = *infoParams.CompanyID

	vInfo.VoucherFiller = *infoParams.VoucherFiller
	vInfo.NumOfMonth = int(count + 1)
	vInfo.CreatedAt = time.Now()
	vInfo.VoucherID = GIdInfoService.genVouIdInfo.GetNextId()
	IdValSli = append(IdValSli, vInfo.VoucherID)
	if err = vs.VInfoDao.Create(ctx, tx, vInfo); err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [VInfoDao.Create: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//create voucherRecord
	vRecord := new(model.VoucherRecord)
	for _, recParam := range params.RecordsParams {
		vRecord.RecordID = GIdInfoService.genVouRecIdInfo.GetNextId()
		IdValSli = append(IdValSli, vRecord.RecordID)
		vRecord.VoucherID = vInfo.VoucherID
		vRecord.SubjectName = *recParam.SubjectName
		vRecord.DebitMoney = *recParam.DebitMoney
		vRecord.CreditMoney = *recParam.CreditMoney
		vRecord.Status = utils.NoAudit
		if recParam.Summary != nil {
			vRecord.Summary = *recParam.Summary
		}
		if recParam.BillCount != nil {
			vRecord.BillCount = *recParam.BillCount
		}
		if recParam.SubID1 != nil {
			vRecord.SubID1 = *recParam.SubID1
		}
		if recParam.SubID2 != nil {
			vRecord.SubID2 = *recParam.SubID2
		}
		if recParam.SubID3 != nil {
			vRecord.SubID3 = *recParam.SubID3
		}
		if recParam.SubID4 != nil {
			vRecord.SubID4 = *recParam.SubID4
		}
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
	vs.Logger.InfoContext(ctx, "CreateVoucher method end ")
	return IdValSli, nil
}

//UpdateVoucher  该函数用于修改voucher ...
func (vs *VoucherService) UpdateVoucher(ctx context.Context, params *model.UpdateVoucherParams,
	requestID string) ([]int, CcError) {
	vs.Logger.InfoContext(ctx, "UpdateVoucher method begin")
	FuncName := "VoucherService/UpdateVoucher"
	bIsRollBack := true
	tx, err := vs.Db.Begin()
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, vs.Logger, FuncName, tx)
		}
	}()
	//update voucherInfo
	if params.ModifyInfoParams != nil {
		voucherInfoParams := make(map[string]interface{}, 3)
		if params.ModifyInfoParams.VoucherMonth != nil {
			voucherInfoParams["voucherMonth"] = *params.ModifyInfoParams.VoucherMonth
			//如果凭证的月份发生了变化，则该voucherInfo里的凭证号也发生变化。
			iMaxNumOfMonth, err := vs.VInfoDao.GetMaxNumByIdAndMonth(ctx, tx, *params.ModifyInfoParams.VoucherMonth,
				*params.ModifyInfoParams.VoucherID)
			if err != nil {
				return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
			}
			voucherInfoParams["numOfMonth"] = iMaxNumOfMonth + 1
		}
		if params.ModifyInfoParams.VoucherDate != nil {
			voucherInfoParams["voucherDate"] = *params.ModifyInfoParams.VoucherDate
		}
		if params.ModifyInfoParams.VoucherFiller != nil {
			voucherInfoParams["voucherFiller"] = *params.ModifyInfoParams.VoucherFiller
		}
		voucherInfoParams["updatedAt"] = time.Now()
		err = vs.VInfoDao.Update(ctx, tx, *params.ModifyInfoParams.VoucherID, voucherInfoParams)
		if err != nil {
			return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
	}
	//update voucher record
	voucherRecordParams := make(map[string]interface{}, 4)
	for _, recParam := range params.ModifyRecordsParams {
		if recParam.VouRecordID == nil || *recParam.VouRecordID == 0 {
			return nil, NewError(ErrVoucher, ErrMiss, ErrId, ErrNull)
		}
		if recParam.Summary != nil {
			voucherRecordParams["summary"] = *recParam.Summary
		}
		if recParam.SubjectName != nil {
			voucherRecordParams["subjectName"] = *recParam.SubjectName
		}
		if recParam.CreditMoney != nil {
			voucherRecordParams["creditMoney"] = *recParam.CreditMoney
		}
		if recParam.DebitMoney != nil {
			voucherRecordParams["debitMoney"] = *recParam.DebitMoney
		}
		if recParam.SubID1 != nil {
			voucherRecordParams["subId1"] = *recParam.SubID1
		}
		err = vs.VRecordDao.UpdateByRecordId(ctx, tx, *recParam.VouRecordID, voucherRecordParams)
		if err != nil {
			return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
		voucherRecordParams = map[string]interface{}{}
	}
	//delete voucher records
	if len(params.DelRecordsParams) > 0 {
		delConditonParams := make(map[string]interface{})
		delConditonParams["recordId"] = params.DelRecordsParams
		err = vs.VRecordDao.DeleteByMultiCondition(ctx, tx, delConditonParams)
		if err != nil {
			return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
	}
	//add voucher records
	var IdValSli []int
	for _, itemParam := range params.AddRecordsParams {
		vRecord := new(model.VoucherRecord)
		vRecord.RecordID = GIdInfoService.genVouRecIdInfo.GetNextId()
		IdValSli = append(IdValSli, vRecord.RecordID)
		vRecord.VoucherID = *itemParam.VoucherID
		vRecord.SubjectName = *itemParam.SubjectName
		vRecord.DebitMoney = *itemParam.DebitMoney
		vRecord.CreditMoney = *itemParam.CreditMoney
		vRecord.Summary = *itemParam.Summary
		vRecord.BillCount = *itemParam.BillCount
		vRecord.Status = utils.NoAudit
		vRecord.SubID1 = *itemParam.SubID1
		// vRecord.SubID2 = *itemParam.SubID2
		// vRecord.SubID3 = *itemParam.SubID3
		// vRecord.SubID4 = *itemParam.SubID4
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
	vs.Logger.InfoContext(ctx, "UpdateVoucher method end ")
	if len(IdValSli) == 0 {
		return nil, nil
	}
	return IdValSli, nil
}

func (vs *VoucherService) DeleteVoucher(ctx context.Context, voucherID int, requestId string) CcError {
	FuncName := "VoucherService/service/DeleteVoucher"
	bIsRollBack := true
	vs.Logger.InfoContext(ctx, "DeleteVoucher method begin, "+"voucher ID:%d", voucherID)
	// Begin transaction
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
	err = vs.VRecordDao.DeleteByVoucherId(ctx, vs.Db, voucherID)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, "Delete failed")
	}
	err = vs.VInfoDao.Delete(ctx, tx, voucherID)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, "Delete failed")
	}
	if err = tx.Commit(); err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	vs.Logger.InfoContext(ctx, "DeleteVoucher method end, "+"voucher ID:%d", voucherID)
	return nil
}

func (vs *VoucherService) GetVoucherByVoucherID(ctx context.Context, voucherID int,
	requestId string) (*model.VoucherView, CcError) {
	FuncName := "VoucherService/service/GetVoucherByVoucherID"
	bIsRollBack := true
	vs.Logger.InfoContext(ctx, "GetVoucherByVoucherID method begin, "+"voucher ID:%d", voucherID)
	//Begin transaction
	tx, err := vs.Db.Begin()
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, vs.Logger, FuncName, tx)
		}
	}()
	//get voucher information
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
	//获取一个voucherID 所对应的所有的voucher records count
	//目前前端没有限制记录的条数，后端暂时先不限制。
	filterFields["voucherId"] = voucherID
	// count, err := vs.VRecordDao.CountByFilter(ctx, tx, filterFields)
	// if err != nil {
	// 	return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	// }
	limit, offset := -1, 0
	// if count > MaxRecordLimit {
	// 	limit = MaxRecordLimit
	// }
	orderField := "recordId"
	orderDirection := cons.Order_Asc
	intervalFilterFields := make(map[string]interface{})
	fuzzyMatchFields := make(map[string]string)
	voucherRecords, err := vs.VRecordDao.List(ctx, tx, filterFields, intervalFilterFields, fuzzyMatchFields,
		limit, offset, orderField, orderDirection)
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[VoucherService/service/GetVoucherByVoucherID] [VRecordDao.List: %s, filterFields: %v]", err.Error(), filterFields)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	for _, vouRecord := range voucherRecords {
		vouRecordView := *(VoucherRecordModelToView(vouRecord))
		recordViewSlice = append(recordViewSlice, vouRecordView)
	}
	//voucherView.VouRecordTotalCount = int(count)
	voucherView.VouRecordTotalCount = len(recordViewSlice)
	voucherView.VouRecordViewSli = append(voucherView.VouRecordViewSli, recordViewSlice...)

	if err = tx.Commit(); err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	vs.Logger.InfoContext(ctx, "GetVoucherByVoucherID method end, "+"voucher ID:%d", voucherID)
	return voucherView, nil
}

//VoucherAudit  该函数用于审核和取消审核 ...
func (vs *VoucherService) VoucherAudit(ctx context.Context, voucherID int,
	params map[string]interface{}, requestID string) CcError {
	FuncName := "VoucherService/VoucherAudit"
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
	voucherParams := make(map[string]interface{}, 2)
	voucherParams["voucherAuditor"] = params["voucherAuditor"]
	voucherParams["updatedAt"] = time.Now()
	//update voucher information
	err = vs.VInfoDao.Update(ctx, tx, voucherID, voucherParams)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//update voucher record
	delete(voucherParams, "voucherAuditor")
	voucherParams["status"] = params["status"]
	err = vs.VRecordDao.UpdateByVoucherId(ctx, tx, voucherID, voucherParams)
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

//之所以放在这，是因为list voucherInfo时，有时，可能需要访问voucher record这个表。
func (vs *VoucherService) ListVoucherInfoByMulCondition(ctx context.Context,
	params *model.ListVoucherInfoParams) ([]*model.VoucherInfoView, int, CcError) {
	vouInfoViewSlice := make([]*model.VoucherInfoView, 0)
	filterFields := make(map[string]interface{})
	intervalFilterFields := make(map[string]interface{})
	limit, offset := -1, 0
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
	FuncName := "VoucherService/ListVoucherInfoByMulCondition"
	bIsRollBack := true
	tx, err := vs.Db.Begin()
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return vouInfoViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, vs.Logger, FuncName, tx)
		}
	}()
	if params.AuxiFilter != nil {
		fuzzyMatchFields := make(map[string]string)
		for _, f := range params.AuxiFilter {
			switch *f.Field {
			case "debitMoney_interval":
				intervalFilterFields["debitMoney"] = f.Value
			case "creditMoney_interval":
				intervalFilterFields["creditMoney"] = f.Value
			case "subjectName_fuzzy":
				fuzzyMatchFields["subjectName"] = f.Value.(string)
			case "summary_fuzzy":
				fuzzyMatchFields["summary"] = f.Value.(string)
			default:
				return vouInfoViewSlice, 0, NewError(ErrVoucherInfo, ErrUnsupported, ErrField, *f.Field)
			}
		}
		voucherRecords, err := vs.VRecordDao.List(ctx, tx, filterFields, intervalFilterFields, fuzzyMatchFields,
			limit, offset, orderField, orderDirection)
		if err != nil {
			vs.Logger.ErrorContext(ctx, "[VoucherService/service/ListVoucherInfoByMulCondition] [VRecordDao.List: %s, filterFields: %v]",
				err.Error(), filterFields)
			return vouInfoViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
		voucherIds := make([]int, 1)
		for _, vouRecord := range voucherRecords {
			voucherIds = append(voucherIds, vouRecord.VoucherID)
		}
		if len(voucherIds) > 0 {
			filterFields["voucherId"] = voucherIds
		}
	}

	delete(intervalFilterFields, "debitMoney")
	delete(intervalFilterFields, "creditMoney")
	if params.BasicFilter != nil {
		for _, f := range params.BasicFilter {
			switch *f.Field {
			case "voucherId", "companyId", "voucherMonth", "numOfMonth", "voucherDate", "voucherFiller", "voucherAuditor":
				filterFields[*f.Field] = f.Value
			case "numOfMonth_interval":
				intervalFilterFields["numOfMonth"] = f.Value
			case "voucherDate_interval":
				intervalFilterFields["voucherDate"] = f.Value
			default:
				return vouInfoViewSlice, 0, NewError(ErrVoucherInfo, ErrUnsupported, ErrField, *f.Field)
			}
		}
	}

	voucherInfos, err := vs.VInfoDao.List(ctx, tx, filterFields, intervalFilterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		vs.Logger.ErrorContext(ctx, "[VoucherInfoService/service/ListVoucherInfo] [VInfoDao.List: %s, filterFields: %v]", err.Error(), filterFields)
		return vouInfoViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return vouInfoViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false

	for _, vouInfo := range voucherInfos {
		vouInfoView := VoucherInfoModelToView(vouInfo)
		vouInfoViewSlice = append(vouInfoViewSlice, vouInfoView)
	}
	vouRecordCount := len(voucherInfos)
	return vouInfoViewSlice, vouRecordCount, nil
}
