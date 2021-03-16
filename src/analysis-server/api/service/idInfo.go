package service

import (
	"database/sql"

	cons "common/constant"
	"common/log"
	"analysis-server/api/db"
	"analysis-server/model"
)

type IDInfoService struct {
	Logger    *log.Logger
	IdInfoDao *db.IDInfoDao
	Db        *sql.DB
}

func (is *IDInfoService) CreateIDInfo(params *model.IDInfoParams,
	requestId string) (*model.IDInfoView, CcError) {
	//create
	is.Logger.Info("CreateIDInfo method start")
	idInfo := new(model.IDInfo)
	idInfo.CompanyID = *params.CompanyID
	idInfo.SubjectID = *params.SubjectID
	idInfo.VoucherID = *params.VoucherID
	idInfo.VoucherRecordID = *params.VoucherRecordID
	if err := is.IdInfoDao.Create(is.Db, idInfo); err != nil {
		is.Logger.Error("[CreateIDInfo] [IdInfoDao.Create: %s]", err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	idInfoView := is.IdInfoModelToView(idInfo)
	is.Logger.Info("CreateIDInfo method end ")
	return idInfoView, nil
}

// convert accSubject to accSubjectView ...
func (is *IDInfoService) IdInfoModelToView(idInfo *model.IDInfo) *model.IDInfoView {
	idInfoView := new(model.IDInfoView)
	idInfoView.CompanyID = idInfo.CompanyID
	idInfoView.SubjectID = idInfo.SubjectID
	idInfoView.VoucherID = idInfo.VoucherID
	idInfoView.VoucherRecordID = idInfo.VoucherRecordID
	return idInfoView
}

func (is *IDInfoService) GetIdInfo() (*model.IDInfoView, CcError) {
	idInfo, err := is.IdInfoDao.Get(is.Db)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, NewCcError(cons.CodeAccSubNotExist, ErrAccSub, ErrNotFound, ErrNull, "the idInfo is not exist")
	default:
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	idInfoView := is.IdInfoModelToView(idInfo)
	return idInfoView, nil
}

func (is *IDInfoService) DeleteIdInfo() CcError {
	is.Logger.Info("DeleteIdInfo method begin")
	err := is.IdInfoDao.Delete(is.Db)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, "Delete failed")
	}
	is.Logger.Info("DeleteIdInfo method end")
	return nil
}

func (is *IDInfoService) UpdateIdInfo(params map[string]interface{}) CcError {
	is.Logger.Info("UpdateIdInfo method begin")
	err := is.IdInfoDao.Update(is.Db,params)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	is.Logger.Info("UpdateIdInfo method end")
	return nil
}
