package service

import (
	"context"
	"database/sql"
	"financeMgr/src/analysis-server/api/db"
	"financeMgr/src/analysis-server/api/utils"
	"financeMgr/src/analysis-server/model"
	cons "financeMgr/src/common/constant"
	"financeMgr/src/common/log"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

// 默认最多返回100条记录，如果记录超过100条，需要在客户端再把剩余的部分获取出来。
const (
	MaxRecordLimit = 100
)

type VoucherService struct {
	Logger     *log.Logger
	VInfoDao   *db.VoucherInfoDao
	VRecordDao *db.VoucherRecordDao
	VouDao     *db.VoucherDao
	Db         *sql.DB
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
	var iYear, iMonth int
	if infoParams.VoucherDate != nil {
		iDate := *infoParams.VoucherDate
		iYear = iDate / 10000
		iMonth = (iDate - iYear*10000) / 100
		vInfo.VoucherDate = iDate
	} else {
		curTime := time.Now()
		var iDay int
		var month time.Month
		iYear, month, iDay = curTime.Date()
		iMonth = int(month)
		vInfo.VoucherDate = iYear*10000 + iMonth*100 + iDay
	}
	vInfo.VoucherMonth = iMonth
	filterFields := make(map[string]interface{})
	filterFields["companyId"] = *infoParams.CompanyID
	filterFields["voucherMonth"] = vInfo.VoucherMonth
	vs.Logger.InfoContext(ctx, "CreateVoucher method start, "+"companyID:%d,VoucherMonth:%d",
		*infoParams.CompanyID, vInfo.VoucherMonth)
	count, err := vs.VInfoDao.CountByFilter(ctx, tx, iYear, filterFields)
	if err != nil {
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//第0个元素是voucherId,从第1个元素开始，就是recordId
	var IdValSli []int
	vInfo.CompanyID = *infoParams.CompanyID
	vInfo.Status = utils.NoAuditVoucher
	if infoParams.BillCount != nil {
		vInfo.BillCount = *infoParams.BillCount
	}
	vInfo.VoucherFiller = *infoParams.VoucherFiller
	vInfo.NumOfMonth = int(count + 1)
	vInfo.CreatedAt = time.Now()
	vInfo.UpdatedAt = time.Now()
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
		if recParam.Summary != nil {
			vRecord.Summary = *recParam.Summary
		}
		if recParam.SubID1 != nil {
			vRecord.SubID1 = *recParam.SubID1
		}
		// if recParam.SubID2 != nil {
		// 	vRecord.SubID2 = *recParam.SubID2
		// }
		// if recParam.SubID3 != nil {
		// 	vRecord.SubID3 = *recParam.SubID3
		// }
		// if recParam.SubID4 != nil {
		// 	vRecord.SubID4 = *recParam.SubID4
		// }
		vRecord.CreatedAt = time.Now()
		vRecord.UpdatedAt = time.Now()
		if err = vs.VRecordDao.Create(ctx, tx, iYear, vRecord); err != nil {
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

// UpdateVoucher  该函数用于修改voucher ...
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
	iVoucherYear := *params.VoucherYear
	//update voucherInfo,目前在client端，暂不支持修改已存凭证的月份或者制作日期，所以下面修改凭证信息的代码，也就无用了。等以后删除。
	if params.ModifyInfoParams != nil {
		voucherInfoParams := make(map[string]interface{}, 3)
		if params.ModifyInfoParams.VoucherMonth != nil {
			voucherInfoParams["voucherMonth"] = *params.ModifyInfoParams.VoucherMonth
			//如果凭证的月份发生了变化，则该voucherInfo里的凭证号也发生变化。
			iMaxNumOfMonth, err := vs.VInfoDao.GetMaxNumByIdAndMonth(ctx, tx, *params.ModifyInfoParams.VoucherMonth,
				iVoucherYear, *params.ModifyInfoParams.VoucherID)
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
		err = vs.VInfoDao.Update(ctx, tx, *params.ModifyInfoParams.VoucherID, iVoucherYear, voucherInfoParams)
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
		voucherRecordParams["updatedAt"] = time.Now()
		err = vs.VRecordDao.UpdateByRecordId(ctx, tx, *recParam.VouRecordID, iVoucherYear, voucherRecordParams)
		if err != nil {
			return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
		voucherRecordParams = map[string]interface{}{}
	}
	//delete voucher records
	if len(params.DelRecordsParams) > 0 {
		delConditonParams := make(map[string]interface{})
		delConditonParams["recordId"] = params.DelRecordsParams
		err = vs.VRecordDao.DeleteByMultiCondition(ctx, tx, iVoucherYear, delConditonParams)
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
		vRecord.SubID1 = *itemParam.SubID1
		// vRecord.SubID2 = *itemParam.SubID2
		// vRecord.SubID3 = *itemParam.SubID3
		// vRecord.SubID4 = *itemParam.SubID4
		vRecord.CreatedAt = time.Now()
		if err = vs.VRecordDao.Create(ctx, tx, iVoucherYear, vRecord); err != nil {
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

func (vs *VoucherService) DeleteVoucher(ctx context.Context, voucherID, iYear int, requestId string) CcError {
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
	err = vs.VRecordDao.DeleteByVoucherId(ctx, tx, voucherID, iYear)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, "Delete failed")
	}
	err = vs.VInfoDao.Delete(ctx, tx, voucherID, iYear)
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

func (vs *VoucherService) GetVoucherByVoucherID(ctx context.Context, voucherID, iYear int,
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
	vInfo, err := vs.VInfoDao.Get(ctx, tx, voucherID, iYear)
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
	orderDirection := utils.OrderAsc
	voucherRecords, err := vs.VRecordDao.SimpleList(ctx, tx, filterFields, iYear, limit, offset, orderDirection, orderField)
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

// VoucherArrange
func (vs *VoucherService) ArrangeVoucher(ctx context.Context, params *model.VoucherArrangeParams,
	requestID string) CcError {
	vs.Logger.InfoContext(ctx, "ArrangeVoucher method begin,companyID:%d ,month:%d",
		*params.CompanyID, *params.VoucherMonth)
	err := vs.deleteInvalidVoucher(ctx, *params.VoucherYear, *params.CompanyID, *params.VoucherMonth)
	if err == nil {
		//update the voucher Num
		if params.ArrangeVoucherNum != nil || *params.ArrangeVoucherNum {
			err = vs.arrangeVoucherNum(ctx, *params.VoucherYear, *params.CompanyID, *params.VoucherMonth)
		}
	}
	vs.Logger.InfoContext(ctx, "ArrangeVoucher method end")
	return err
}

func (vs *VoucherService) deleteInvalidVoucher(ctx context.Context, iVoucherYear, companyID, voucherMonth int) CcError {
	FuncName := "VoucherService/service/deleteInvalidVoucher"
	bIsRollBack := true
	vs.Logger.InfoContext(ctx, "deleteInvalidVoucher method begin,companyID:%d ,month:%d", companyID, voucherMonth)
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
	filterFields := make(map[string]interface{}, 3)
	limit, offset := -1, 0
	filterFields["companyId"] = companyID
	filterFields["voucherMonth"] = voucherMonth
	filterFields["status"] = utils.InvalidVoucher
	orderField := ""
	orderDirection := 0
	voucherInfos, err := vs.VInfoDao.SimpleList(ctx, tx, filterFields, iVoucherYear, limit, offset, orderDirection, orderField)
	if err != nil {
		errInfo := fmt.Sprintf("[VoucherService/service/deleteInvalidVoucher] [VInfoDao.List: %s, filterFields: %v]",
			err.Error(), filterFields)
		vs.Logger.ErrorContext(ctx, errInfo)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if len(voucherInfos) > 0 {
		voucherIdSlice := make([]int, len(voucherInfos))
		for i := 0; i < len(voucherInfos); i++ {
			//delete voucher record
			voucherIdSlice = append(voucherIdSlice, voucherInfos[i].VoucherID)
			err = vs.VRecordDao.DeleteByVoucherId(ctx, tx, voucherInfos[i].VoucherID, iVoucherYear)
			if err != nil {
				return NewError(ErrSystem, ErrError, ErrNull, "Delete voucher record failed")
			}
		}
		//batch,delete voucher information
		err = vs.VInfoDao.BatchDelete(ctx, tx, iVoucherYear, voucherIdSlice)
		if err != nil {
			return NewError(ErrSystem, ErrError, ErrNull, "Delete voucher information failed")
		}
	}
	if err = tx.Commit(); err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	vs.Logger.InfoContext(ctx, "deleteInvalidVoucher method end")
	return nil
}

func (vs *VoucherService) arrangeVoucherNum(ctx context.Context, iVoucherYear, companyID int, voucherMonth int) CcError {
	vs.Logger.InfoContext(ctx, "arrangeVoucherNum method begin,companyID:%d ,month:%d", companyID, voucherMonth)
	//Begin transaction
	FuncName := "VoucherService/service/arrangeVoucherNum"
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
	//update the voucher Num
	filterFields := make(map[string]interface{}, 1)
	limit, offset := -1, 0
	orderField := "numOfMonth"
	orderDirection := 0
	filterFields["voucherMonth"] = voucherMonth
	filterFields["companyId"] = companyID
	voucherInfos, err := vs.VInfoDao.SimpleList(ctx, tx, filterFields, iVoucherYear, limit, offset, orderDirection, orderField)
	if err != nil {
		errInfo := fmt.Sprintf("[VoucherService/service/arrangeVoucherNum] [VInfoDao.List: %s, filterFields: %v]",
			err.Error(), filterFields)
		vs.Logger.ErrorContext(ctx, errInfo)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	voucherInfoParams := make(map[string]interface{}, 2)
	for i := 0; i < len(voucherInfos); i++ {
		//update numOfMonth
		voucherInfoParams["numOfMonth"] = i + 1
		voucherInfoParams["updatedAt"] = time.Now()
		//update voucher information
		err = vs.VInfoDao.Update(ctx, tx, voucherInfos[i].VoucherID, iVoucherYear, voucherInfoParams)
		if err != nil {
			return NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
	}
	if err = tx.Commit(); err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	vs.Logger.InfoContext(ctx, "arrangeVoucherNum method end")
	return nil
}

// 该函数时利用过滤voucherRecord的条件来过滤voucherInfo .在该函数实现中，增加了如下的功能
// 1、在where条件里增加between ... and
// 2、增加了like
// 3、增加了多列排序。
func (vs *VoucherService) ListVoucherInfoWithAuxCondition(ctx context.Context,
	params *model.ListVoucherInfoParams) ([]*model.VoucherInfoView, int, CcError) {

	vouInfoViewSlice := make([]*model.VoucherInfoView, 0)
	filterNo := make(map[string]interface{})
	filterFields := make(map[string]interface{})
	intervalFilterFields := make(map[string]interface{})
	limit, offset := -1, 0
	if params.DescLimit != nil {
		limit = *params.DescLimit
		if params.DescOffset != nil {
			offset = *params.DescOffset
		}
	}

	iVoucherYear := 0
	FuncName := "VoucherService/ListVoucherInfoWithAuxCondition"
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
	//这是对voucherInfo 的list
	if params.BasicFilter != nil {
		//var intervalValSlice []int
		for _, f := range params.BasicFilter {
			switch *f.Field {
			case "voucherYear":
				{
					switch f.Value.(type) {
					case float64:
						//从客户端发过来的，解析json时，会解析成float64 (经验证该结论是错误的)
						//正确的结论是，文档显示当把json解析成interface{}时，把number解析成float64
						iVoucherYear = int(f.Value.(float64))
					case int:
						panic("the iVoucherYear is int")
						//从cli发过来的，解析json时，会解析成int (经验证该结论是错误的)
						//iVoucherYear = f.Value.(int)
					}
				}
			case "voucherId", "companyId", "voucherMonth", "numOfMonth", "voucherDate", "voucherFiller":
				fallthrough
			case "voucherAuditor", "status":
				filterFields[*f.Field] = f.Value
			case "status_no":
				filterNo["status"] = f.Value
			case "numOfMonth_interval":
				intervalFilterFields["numOfMonth"] = f.Value
			case "voucherDate_interval":
				intervalFilterFields["voucherDate"] = f.Value
			case "voucherMonth_interval":
				intervalFilterFields["voucherMonth"] = f.Value
			default:
				return vouInfoViewSlice, 0, NewError(ErrVoucherInfo, ErrUnsupported, ErrField, *f.Field)
			}
		}
	}
	//这是对voucherRecord 的list,根据其结果作为voucherInfo的过滤条件。
	if params.AuxiFilter != nil {
		//filterRecNo := make(map[string]interface{})
		filterRecFields := make(map[string]interface{})
		intervalFilterRecFields := make(map[string]interface{})
		fuzzyMatchFields := make(map[string]string)
		for _, f := range params.AuxiFilter {
			switch *f.Field {
			case "recordId", "voucherId", "subjectName", "summary", "subId1":
				filterRecFields[*f.Field] = f.Value
			case "subId1Arr":
				filterRecFields["subId1"] = f.Value
			case "debitMoney_interval":
				intervalFilterRecFields["debitMoney"] = f.Value
			case "creditMoney_interval":
				intervalFilterRecFields["creditMoney"] = f.Value
			case "subjectName_fuzzy":
				fuzzyMatchFields["subjectName"] = f.Value.(string)
			case "summary_fuzzy":
				fuzzyMatchFields["summary"] = f.Value.(string)
			//因为此时只有单个值的时候，就是>,<的操作，所以此时的类型为string
			//这个分析是错误的，因为在数据库这两个字段是decimal类型，不能比较，等以后转换成整形后，再支持这个操作。
			// case "debitMoney", "creditMoney":
			// 	filterRecFields[*f.Field] = f.Value.(string)
			default:
				return vouInfoViewSlice, 0, NewError(ErrVoucherInfo, ErrUnsupported, ErrField, *f.Field)
			}
		}
		voucherRecords, err := vs.VRecordDao.List(ctx, tx, nil, filterRecFields, intervalFilterRecFields,
			fuzzyMatchFields, nil, iVoucherYear, limit, offset)
		if err != nil {
			vs.Logger.ErrorContext(ctx, "[VoucherService/service/ListVoucherInfoWithAuxCondition] [VRecordDao.List: %s, filterRecFields: %v]",
				err.Error(), filterRecFields)
			return vouInfoViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
		voucherIds := make([]int, 0)
		for _, vouRecord := range voucherRecords {
			voucherIds = append(voucherIds, vouRecord.VoucherID)
		}
		if len(voucherIds) > 0 {
			filterFields["voucherId"] = voucherIds
		}
	}

	voucherInfos, err := vs.VInfoDao.List(ctx, tx, filterNo, filterFields, intervalFilterFields, nil,
		params.Order, iVoucherYear, limit, offset)
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
	voucherInfoCount := len(voucherInfos)
	return vouInfoViewSlice, voucherInfoCount, nil
}

// 	voucherInfos, err := vs.VInfoDao.List(ctx, tx, filterNo, filterFields, intervalFilterFields, nil,
// 		orderFilter, iVoucherYear, limit, offset)
// 	if err != nil {
// 		vs.Logger.ErrorContext(ctx, "[VoucherInfoService/service/ListVoucherInfo] [VInfoDao.List: %s, filterFields: %v]", err.Error(), filterFields)
// 		return vouInfoViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
// 	}
// 	if err = tx.Commit(); err != nil {
// 		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
// 		return vouInfoViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
// 	}
// 	bIsRollBack = false

// 	for _, vouInfo := range voucherInfos {
// 		vouInfoView := VoucherInfoModelToView(vouInfo)
// 		vouInfoViewSlice = append(vouInfoViewSlice, vouInfoView)
// 	}
// 	voucherInfoCount := len(voucherInfos)
// 	return vouInfoViewSlice, voucherInfoCount, nil
// }

func (vs *VoucherService) CalcAccuMoney(ctx context.Context,
	params *model.CalAccuMoneyParams, requestId string) (*model.AccuMoneyValueView, CcError) {
	vs.Logger.InfoContext(ctx, "CalcAccuMoney method begin,companyID:%d ,subjectID:%d",
		*params.CompanyID, *params.SubjectID)
	FuncName := "VoucherService/service/CalcAccuMoney"
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
	var calcAccuMoney model.CalAccuMoney
	calcAccuMoney.CompanyID = *params.CompanyID
	calcAccuMoney.SubjectID = *params.SubjectID
	calcAccuMoney.VoucherDate = *params.VoucherDate
	calcAccuMoney.VoucherYear = *params.VoucherYear
	calcAccuMoney.Status = *params.Status
	var accuMoney *model.AccuMoneyValueView
	accuMoney, err = vs.VouDao.CalcAccuMoney(ctx, tx, &calcAccuMoney)
	if err != nil {
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	vs.Logger.InfoContext(ctx, "CalcAccuMoney method end")
	return accuMoney, nil
}

// 批量计算多个accSubId所对应的累计金额
func (vs *VoucherService) BatchCalcAccuMoney(ctx context.Context,
	params *model.BatchCalAccuMoneyParams, requestId string) ([]*model.AccuMoneyValueView, CcError) {
	vs.Logger.InfoContext(ctx, "BatchCalcAccuMoney method begin")
	FuncName := "VoucherService/service/BatchCalcAccuMoney"
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
	resData := make([]*model.AccuMoneyValueView, 0, len(params.SubjectIDArr))
	for _, subId := range params.SubjectIDArr {
		var calcAccuMoney model.CalAccuMoney
		calcAccuMoney.CompanyID = *params.CompanyID
		calcAccuMoney.SubjectID = subId
		calcAccuMoney.VoucherDate = *params.VoucherDate
		calcAccuMoney.VoucherYear = *params.VoucherYear
		calcAccuMoney.Status = *params.Status
		var accuMoney *model.AccuMoneyValueView
		accuMoney, err = vs.VouDao.CalcAccuMoney(ctx, tx, &calcAccuMoney)
		if err != nil {
			return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
		accuMoney.SubjectID = subId
		resData = append(resData, accuMoney)
	}
	if err = tx.Commit(); err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	vs.Logger.InfoContext(ctx, "BatchCalcAccuMoney method end")
	return resData, nil
}

// 批量计算多个accSubId所对应的本期发生额
func (vs *VoucherService) CalcAccountOfPeriod(ctx context.Context,
	params *model.CalAmountOfPeriodParams, requestId string) ([]*model.AccuMoneyValueView, CcError) {
	vs.Logger.InfoContext(ctx, "CalcAccountOfPeriod method begin")
	FuncName := "VoucherService/service/CalcAccountOfPeriod"
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
	recData, err := vs.VouDao.GetPartialVouRecords(ctx, tx, params)
	if err != nil {
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	resDataMap := make(map[int]*model.AccuMoneyValueView, len(recData))
	for _, accPeriod := range recData {
		if itemPtr, ok := resDataMap[accPeriod.SubjectID]; ok {
			itemPtr.AccuDebitMoney += accPeriod.PeriodDebitMoney
			itemPtr.AccuCreditMoney += accPeriod.PeriodCreditMoney
		} else {
			accPeriodView := new(model.AccuMoneyValueView)
			accPeriodView.SubjectID = accPeriod.SubjectID
			accPeriodView.AccuDebitMoney = accPeriod.PeriodDebitMoney
			accPeriodView.AccuCreditMoney = accPeriod.PeriodCreditMoney
			resDataMap[accPeriod.SubjectID] = accPeriodView
		}
	}
	accPeriodViewSlice := make([]*model.AccuMoneyValueView, len(resDataMap))
	for _, accPeriodPtr := range resDataMap {
		accPeriodViewSlice = append(accPeriodViewSlice, accPeriodPtr)
	}
	if err = tx.Commit(); err != nil {
		vs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	vs.Logger.InfoContext(ctx, "CalcAccountOfPeriod method end")
	return accPeriodViewSlice, nil
}
