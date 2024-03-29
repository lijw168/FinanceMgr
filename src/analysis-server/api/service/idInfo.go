package service

import (
	"database/sql"

	"financeMgr/src/analysis-server/api/db"
	aUtils "financeMgr/src/analysis-server/api/utils"
	"financeMgr/src/analysis-server/model"
	cons "financeMgr/src/common/constant"
	"financeMgr/src/common/log"
)

type IDInfoService struct {
	logger            *log.Logger
	idInfoDao         *db.IDInfoDao
	_db               *sql.DB
	genSubIdInfo      *aUtils.GenIdInfo
	genComIdInfo      *aUtils.GenIdInfo
	genVouIdInfo      *aUtils.GenIdInfo
	genVouRecIdInfo   *aUtils.GenIdInfo
	genOptIdInfo      *aUtils.GenIdInfo
	genComGroupIdInfo *aUtils.GenIdInfo
	genvVouTempIdInfo *aUtils.GenIdInfo
}

func NewIDInfoService() *IDInfoService {
	idInfoService := IDInfoService{}
	return &idInfoService
}

func (is *IDInfoService) InitIdInfoService(logger *log.Logger, idInfoDao *db.IDInfoDao, _db *sql.DB) {
	is.logger = logger
	is.idInfoDao = idInfoDao
	is._db = _db
}

func (is *IDInfoService) InitIdResource() CcError {
	idInfoView, ccErr := is.GetIdInfo()
	if ccErr != nil {
		return ccErr
	}
	var err error
	is.genSubIdInfo, err = aUtils.NewGenIdInfo(idInfoView.SubjectID)
	if err != nil {
		is.logger.LogError("[InitGenIdInfo] genSubIdInfo,failed: ", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	is.genComIdInfo, err = aUtils.NewGenIdInfo(idInfoView.CompanyID)
	if err != nil {
		is.logger.LogError("[InitGenIdInfo] genComIdInfo,failed: ", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	is.genVouIdInfo, err = aUtils.NewGenIdInfo(idInfoView.VoucherID)
	if err != nil {
		is.logger.LogError("[InitGenIdInfo] genVouIdInfo,failed:", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	is.genVouRecIdInfo, err = aUtils.NewGenIdInfo(idInfoView.VoucherRecordID)
	if err != nil {
		is.logger.LogError("[InitGenIdInfo] genVouRecIdInfo,failed: ", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	is.genOptIdInfo, err = aUtils.NewGenIdInfo(idInfoView.OperatorID)
	if err != nil {
		is.logger.LogError("[InitGenIdInfo] genOptIdInfo,failed: ", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	is.genComGroupIdInfo, err = aUtils.NewGenIdInfo(idInfoView.ComGroupID)
	if err != nil {
		is.logger.LogError("[InitGenIdInfo] genComGroupIdInfo,failed: ", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	is.genvVouTempIdInfo, err = aUtils.NewGenIdInfo(idInfoView.VoucherTemplateID)
	if err != nil {
		is.logger.LogError("[InitGenIdInfo] genvVouTempIdInfo,failed: ", err.Error())
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	return nil
}

func (is *IDInfoService) CreateIDInfo(params *model.IDInfoParams,
	requestId string) (*model.IDInfoView, CcError) {
	//create
	is.logger.Info("CreateIDInfo method start")
	idInfo := new(model.IDInfo)
	idInfo.CompanyID = *params.CompanyID
	idInfo.OperatorID = *params.OperatorID
	idInfo.SubjectID = *params.SubjectID
	idInfo.VoucherID = *params.VoucherID
	idInfo.VoucherRecordID = *params.VoucherRecordID
	idInfo.ComGroupID = *params.ComGroupID
	idInfo.VoucherTemplateID = *params.VoucherTemplateID
	if err := is.idInfoDao.Create(is._db, idInfo); err != nil {
		is.logger.Error("[CreateIDInfo] [IdInfoDao.Create: %s]", err.Error())
		return nil, NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	idInfoView := is.IdInfoModelToView(idInfo)
	is.logger.Info("CreateIDInfo method end ")
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
	idInfo, err := is.idInfoDao.Get(is._db)
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
	is.logger.Info("DeleteIdInfo method begin")
	err := is.idInfoDao.Delete(is._db)
	if err != nil {
		return NewError(ErrIdInfo, ErrError, ErrNull, err.Error())
	}
	is.logger.Info("DeleteIdInfo method end")
	return nil
}

func (is *IDInfoService) UpdateIdInfo(params map[string]interface{}) CcError {
	is.logger.Info("UpdateIdInfo method begin")
	err := is.idInfoDao.Update(is._db, params)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	is.logger.Info("UpdateIdInfo method end")
	return nil
}

func (is *IDInfoService) WriteIdResourceToDb() CcError {
	is.logger.Info("WriteIdResourceToDb method begin")
	subId := is.genSubIdInfo.GetId()
	optId := is.genOptIdInfo.GetId()
	comId := is.genComIdInfo.GetId()
	vouId := is.genVouIdInfo.GetId()
	vouRecId := is.genVouRecIdInfo.GetId()
	comGroupId := is.genComGroupIdInfo.GetId()
	vouTempId := is.genvVouTempIdInfo.GetId()
	updateFields := make(map[string]interface{})
	updateFields["subjectId"] = subId
	updateFields["operatorId"] = optId
	updateFields["companyId"] = comId
	updateFields["voucherId"] = vouId
	updateFields["voucherRecordId"] = vouRecId
	updateFields["companyGroupId"] = comGroupId
	updateFields["voucherTemplateId"] = vouTempId
	ccErr := is.UpdateIdInfo(updateFields)
	if ccErr != nil {
		is.logger.Error("WriteIdResourceToDb failed,errInfo:%s", ccErr.Error())
	}
	is.logger.Info("WriteIdResourceToDb method end")
	return ccErr
}
