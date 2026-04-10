package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"financeMgr/src/analysis-server/api/db"
	dbp "financeMgr/src/analysis-server/api/db"
	aUtils "financeMgr/src/analysis-server/api/utils"
	"financeMgr/src/analysis-server/model"
	cons "financeMgr/src/common/constant"
)

type IDInfoService struct {
	idInfoDao         *db.IDInfoDao
	VInfoDao          *db.VoucherInfoDao
	VRecordDao        *db.VoucherRecordDao
	AccSubDao         *db.AccSubDao
	CompanyDao        *db.CompanyDao
	OptInfoDao        *db.OperatorInfoDao
	ComGroupDao       *db.CompanyGroupDao
	VTemplateDao      *db.VoucherTemplateDao
	genSubIdInfo      *aUtils.GenIdInfo
	genComIdInfo      *aUtils.GenIdInfo
	genVouIdInfo      *aUtils.GenIdInfo
	genVouRecIdInfo   *aUtils.GenIdInfo
	genOptIdInfo      *aUtils.GenIdInfo
	genComGroupIdInfo *aUtils.GenIdInfo
	genvVouTempIdInfo *aUtils.GenIdInfo
	dbName            string
}

func NewIDInfoService() *IDInfoService {
	idInfoService := IDInfoService{}
	return &idInfoService
}

func (is *IDInfoService) InitIdInfoService(idInfoDao *db.IDInfoDao) {
	is.idInfoDao = idInfoDao
}

func (is *IDInfoService) InitIdResource() CcError {
	idInfoView, ccErr := is.GetIdInfo()
	if ccErr != nil {
		return ccErr
	}
	ccErr = is.verifyIdInfoAndUpdate(idInfoView)
	if ccErr != nil {
		return ccErr
	}
	var err error
	is.genSubIdInfo, err = aUtils.NewGenIdInfo(idInfoView.SubjectID)
	if err != nil {
		gLogger.LogError("[InitGenIdInfo] genSubIdInfo,failed: ", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	is.genComIdInfo, err = aUtils.NewGenIdInfo(idInfoView.CompanyID)
	if err != nil {
		gLogger.LogError("[InitGenIdInfo] genComIdInfo,failed: ", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	is.genVouIdInfo, err = aUtils.NewGenIdInfo(idInfoView.VoucherID)
	if err != nil {
		gLogger.LogError("[InitGenIdInfo] genVouIdInfo,failed:", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	is.genVouRecIdInfo, err = aUtils.NewGenIdInfo(idInfoView.VoucherRecordID)
	if err != nil {
		gLogger.LogError("[InitGenIdInfo] genVouRecIdInfo,failed: ", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	is.genOptIdInfo, err = aUtils.NewGenIdInfo(idInfoView.OperatorID)
	if err != nil {
		gLogger.LogError("[InitGenIdInfo] genOptIdInfo,failed: ", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	is.genComGroupIdInfo, err = aUtils.NewGenIdInfo(idInfoView.ComGroupID)
	if err != nil {
		gLogger.LogError("[InitGenIdInfo] genComGroupIdInfo,failed: ", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	is.genvVouTempIdInfo, err = aUtils.NewGenIdInfo(idInfoView.VoucherTemplateID)
	if err != nil {
		gLogger.LogError("[InitGenIdInfo] genvVouTempIdInfo,failed: ", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	return nil
}

func (is *IDInfoService) CreateIDInfo(params *model.IDInfoParams,
	requestId string) (*model.IDInfoView, CcError) {
	//create
	gLogger.Info("CreateIDInfo method start")
	idInfo := new(model.IDInfo)
	idInfo.CompanyID = *params.CompanyID
	idInfo.OperatorID = *params.OperatorID
	idInfo.SubjectID = *params.SubjectID
	idInfo.VoucherID = *params.VoucherID
	idInfo.VoucherRecordID = *params.VoucherRecordID
	idInfo.ComGroupID = *params.ComGroupID
	idInfo.VoucherTemplateID = *params.VoucherTemplateID
	idInfo.UpdatedAt = time.Now()
	if err := is.idInfoDao.Create(gDb, idInfo); err != nil {
		gLogger.Error("[CreateIDInfo] [IdInfoDao.Create: %s]", err.Error())
		return nil, NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	idInfoView := is.IdInfoModelToView(idInfo)
	gLogger.Info("CreateIDInfo method end ")
	return idInfoView, nil
}

// convert accSubject to accSubjectView ...
func (is *IDInfoService) IdInfoModelToView(idInfo *model.IDInfo) *model.IDInfoView {
	idInfoView := new(model.IDInfoView)
	idInfoView.CompanyID = idInfo.CompanyID
	idInfoView.SubjectID = idInfo.SubjectID
	idInfoView.VoucherID = idInfo.VoucherID
	idInfoView.VoucherRecordID = idInfo.VoucherRecordID
	idInfoView.OperatorID = idInfo.OperatorID
	idInfoView.ComGroupID = idInfo.ComGroupID
	idInfoView.VoucherTemplateID = idInfo.VoucherTemplateID
	return idInfoView
}

func (is *IDInfoService) GetIdInfo() (*model.IDInfoView, CcError) {
	idInfo, err := is.idInfoDao.Get(gDb)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, NewCcError(cons.CodeIdInfoNotExist, ErrIdInfo, ErrNotFound, ErrNull, "the idInfo is not exist")
	default:
		return nil, NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	idInfoView := is.IdInfoModelToView(idInfo)
	return idInfoView, nil
}

func (is *IDInfoService) DeleteIdInfo() CcError {
	gLogger.Info("DeleteIdInfo method begin")
	err := is.idInfoDao.Delete(gDb)
	if err != nil {
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	gLogger.Info("DeleteIdInfo method end")
	return nil
}

func (is *IDInfoService) UpdateIdInfo(params map[string]interface{}) CcError {
	gLogger.Info("UpdateIdInfo method begin")
	err := is.idInfoDao.Update(gDb, params)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	gLogger.Info("UpdateIdInfo method end")
	return nil
}

func (is *IDInfoService) WriteIdResourceToDb() CcError {
	gLogger.Info("WriteIdResourceToDb method begin")
	updateFields := make(map[string]interface{})
	//下面的代码，根据是否发生变化，来决定是否更新数据库，如果没有发生变化，就不更新数据库，减少数据库的压力
	if is.genSubIdInfo.IsChanged() {
		subId := is.genSubIdInfo.GetId(true)
		updateFields["subjectId"] = subId
	}
	if is.genComIdInfo.IsChanged() {
		comId := is.genComIdInfo.GetId(true)
		updateFields["companyId"] = comId
	}
	if is.genVouIdInfo.IsChanged() {
		vouId := is.genVouIdInfo.GetId(true)
		updateFields["voucherId"] = vouId
	}
	if is.genVouRecIdInfo.IsChanged() {
		vouRecId := is.genVouRecIdInfo.GetId(true)
		updateFields["voucherRecordId"] = vouRecId
	}
	if is.genOptIdInfo.IsChanged() {
		optId := is.genOptIdInfo.GetId(true)
		updateFields["operatorId"] = optId
	}
	if is.genComGroupIdInfo.IsChanged() {
		comGroupId := is.genComGroupIdInfo.GetId(true)
		updateFields["companyGroupId"] = comGroupId
	}
	if is.genvVouTempIdInfo.IsChanged() {
		vouTempId := is.genvVouTempIdInfo.GetId(true)
		updateFields["voucherTemplateId"] = vouTempId
	}
	updateFields["updatedAt"] = time.Now()
	ccErr := is.UpdateIdInfo(updateFields)
	if ccErr != nil {
		gLogger.Error("WriteIdResourceToDb failed,errInfo:%s", ccErr.Error())
	}
	gLogger.Info("WriteIdResourceToDb method end")
	return ccErr
}

func (is *IDInfoService) verifyIdInfoAndUpdate(idInfoView *model.IDInfoView) CcError {
	//get latest voucherTable,voucherRecordTable
	voucherInfoTab, voucherRecTab, err := is.getLatestYearOfVoucherTable()
	if err != nil {
		gLogger.Error("getLatestYearOfVoucherTable failed,errInfo:%s", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	//get max subjectId
	maxSubId, err := is.AccSubDao.GetMaxSubId(context.TODO(), gDb)
	if err != nil {
		gLogger.Error("GetMaxSubjectId failed,errInfo:%s", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	//get max companyId
	maxCompanyId, err := is.CompanyDao.GetMaxCompanyId(context.TODO(), gDb)
	if err != nil {
		gLogger.Error("GetMaxCompanyId failed,errInfo:%s", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	//get max voucherId
	maxVoucherId, err := is.VInfoDao.GetMaxVoucherIdInLatestVoucherInfoTable(context.TODO(), gDb, voucherInfoTab)
	if err != nil {
		gLogger.Error("GetMaxVoucherIdInLatestVoucherInfoTable failed,errInfo:%s", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	//get max voucherRecordId
	maxVoucherRecordId, err := is.VRecordDao.GetMaxRecordIdInLatestVoucherRecordTable(context.TODO(), gDb, voucherRecTab)
	if err != nil {
		gLogger.Error("GetMaxRecordIdInLatestVoucherRecordTable failed,errInfo:%s", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	//get max operatorId
	maxOperatorId, err := is.OptInfoDao.GetMaxOperatorId(context.TODO(), gDb)
	if err != nil {
		gLogger.Error("GetMaxOperatorId failed,errInfo:%s", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	//get max companyGroupId
	maxCompanyGroupId, err := is.ComGroupDao.GetMaxCompanyGroupId(context.TODO(), gDb)
	if err != nil {
		gLogger.Error("GetMaxCompanyGroupId failed,errInfo:%s", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	//get max voucherTemplateId
	maxVoucherTemplateId, err := is.VTemplateDao.GetMaxVoucherTemplateId(context.TODO(), gDb)
	if err != nil {
		gLogger.Error("GetMaxVoucherTemplateId failed,errInfo:%s", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	updateFields := make(map[string]interface{})
	if idInfoView.SubjectID < maxSubId {
		gLogger.Error("the subjectId in idInfo is invalid,subjectId should be greater than or equal to %d", maxSubId)
		//update idInfo with maxSubId
		idInfoView.SubjectID = maxSubId
		updateFields["subjectId"] = maxSubId
	}
	if idInfoView.CompanyID < maxCompanyId {
		gLogger.Error("the companyId in idInfo is invalid,companyId should be greater than or equal to %d", maxCompanyId)
		//update idInfo with maxCompanyId
		idInfoView.CompanyID = maxCompanyId
		updateFields["companyId"] = maxCompanyId
	}
	if idInfoView.VoucherID < maxVoucherId {
		gLogger.Error("the voucherId in idInfo is invalid,voucherId should be greater than or equal to %d", maxVoucherId)
		//update idInfo with maxVoucherId
		idInfoView.VoucherID = maxVoucherId
		updateFields["voucherId"] = maxVoucherId
	}
	if idInfoView.VoucherRecordID < maxVoucherRecordId {
		gLogger.Error("the voucherRecordId in idInfo is invalid,voucherRecordId should be greater than or equal to %d", maxVoucherRecordId)
		//update idInfo with maxVoucherRecordId
		idInfoView.VoucherRecordID = maxVoucherRecordId
		updateFields["voucherRecordId"] = maxVoucherRecordId

	}
	if idInfoView.OperatorID < maxOperatorId {
		gLogger.Error("the operatorId in idInfo is invalid,operatorId should be greater than or equal to %d", maxOperatorId)
		//update idInfo with maxOperatorId
		idInfoView.OperatorID = maxOperatorId
		updateFields["operatorId"] = maxOperatorId
	}
	if idInfoView.ComGroupID < maxCompanyGroupId {
		gLogger.Error("the companyGroupId in idInfo is invalid,companyGroupId should be greater than or equal to %d", maxCompanyGroupId)
		//update idInfo with maxCompanyGroupId
		idInfoView.ComGroupID = maxCompanyGroupId
		updateFields["companyGroupId"] = maxCompanyGroupId
	}
	if idInfoView.VoucherTemplateID < maxVoucherTemplateId {
		gLogger.Error("the voucherTemplateId in idInfo is invalid,voucherTemplateId should be greater than or equal to %d", maxVoucherTemplateId)
		//update idInfo with maxVoucherTemplateId
		idInfoView.VoucherTemplateID = maxVoucherTemplateId
		updateFields["voucherTemplateId"] = maxVoucherTemplateId
	}
	ccErr := is.UpdateIdInfo(updateFields)
	if ccErr != nil {
		gLogger.Error("UpdateIdInfo failed,errInfo:%s,in verifyIdInfoAndUpdate", ccErr.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, ccErr.Error())
	}
	return nil
}

/*SELECT CAST(SUBSTRING(table_name, LENGTH('voucherInfo_') + 1) AS UNSIGNED) as year
FROM information_schema.TABLES   WHERE table_name LIKE 'voucherInfo_%' ORDER BY year DESC  LIMIT 1;*/
// getLatestYearOfVoucherTable 获取当前最新的凭证表的年份
func (is *IDInfoService) getLatestYearOfVoucherTable() (voucherInfoTab, voucherRecTab string, err error) {
	if is.dbName == "" {
		//get dbName
		err = gDb.QueryRow("SELECT DATABASE()").Scan(&is.dbName)
		if err != nil {
			gLogger.Error("[init/service/getLatestYearOfVoucherTable] [db.QueryRow: %s]", err.Error())
			return
		}
		gLogger.Info("[init/service/getLatestYearOfVoucherTable] current dbName: %s", is.dbName)
	}
	var maxTableName string
	strSql := "select table_name from information_schema.TABLES where table_schema = '" + is.dbName + "' and  table_name like 'voucherInfo_%' order by table_name  desc limit 1"
	err = gDb.QueryRow(strSql).Scan(&maxTableName)
	if err != nil {
		gLogger.Error("[init/service/getLatestYearOfVoucherTable] [db.QueryRowContext: %s]", err.Error())
		return
	}
	// 从表名中提取年份
	if strings.HasPrefix(maxTableName, "voucherInfo_") {
		yearStr := strings.TrimPrefix(maxTableName, "voucherInfo_")
		if year, err := strconv.Atoi(yearStr); err == nil {
			voucherRecTab = dbp.GenTableName(year, "voucherRecordInfo")
			voucherInfoTab = maxTableName
		}
	} else {
		gLogger.Warn("[init/service/getLatestYearOfVoucherTable] No voucherInfo_ table found")
		err = fmt.Errorf("no voucherInfo_ table found")
	}
	return
}
