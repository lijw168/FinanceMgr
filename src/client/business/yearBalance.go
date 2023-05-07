package business

import (
	"financeMgr/src/client/util"
)

type YearBalGateway struct {
}

func (yg *YearBalGateway) GetYearBalance(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.GetYearBalance_json(param); err != nil {
		errCode = util.ErrShowFailed
		logger.Error("the GetYearBalance failed,err:%v", err.Error())
	} else {
		logger.Debug("GetYearBalance succeed.")
	}
	return resData, errCode
}

func (yg *YearBalGateway) CreateYearBalance(param []byte) (errCode int) {
	errCode = util.ErrNull
	var err error
	if err = cSdk.CreateYearBalance_json(param); err != nil {
		errCode = util.ErrCreateFailed
		logger.Error("the CreateYearBalance failed,err:%v", err.Error())
	} else {
		logger.Debug("CreateYearBalance succeed;")
	}
	return errCode
}

func (yg *YearBalGateway) BatchCreateYearBalance(param []byte) (errCode int) {
	errCode = util.ErrNull
	var err error
	if err = cSdk.BatchCreateYearBalance_json(param); err != nil {
		errCode = util.ErrCreateFailed
		logger.Error("the BatchCreateYearBalance failed,err:%v", err.Error())
	} else {
		logger.Debug("BatchCreateYearBalance succeed;")
	}
	return errCode
}

func (yg *YearBalGateway) UpdateYearBalance(param []byte) (errCode int) {
	errCode = util.ErrNull
	if err := cSdk.UpdateYearBalance_json(param); err != nil {
		errCode = util.ErrUpdateFailed
		logger.Error("the UpdateYearBalance failed,err:%v", err.Error())
	} else {
		logger.Debug("UpdateYearBalance succeed")
	}
	return errCode
}

func (yg *YearBalGateway) BatchUpdateYearBalance(param []byte) (errCode int) {
	errCode = util.ErrNull
	if err := cSdk.BatchUpdateYearBalance_json(param); err != nil {
		errCode = util.ErrUpdateFailed
		logger.Error("the BatchUpdateYearBalance failed,err:%v", err.Error())
	} else {
		logger.Debug("BatchUpdateYearBalance succeed")
	}
	return errCode
}

func (yg *YearBalGateway) DeleteYearBalance(param []byte) (errCode int) {
	errCode = util.ErrNull
	if err := cSdk.DeleteYearBalance_json(param); err != nil {
		errCode = util.ErrDeleteFailed
		logger.Error("the DeleteYearBalanceByID failed,err:%v", err.Error())
	} else {
		logger.Debug("DeleteYearBalanceByID succeed")
	}
	return errCode
}

func (yg *YearBalGateway) ListYearBalance(param []byte) (resData []byte, errCode int) {
	return listCmdJson(resource_type_year_balance, param, cSdk.ListYearBalance_json)
}
