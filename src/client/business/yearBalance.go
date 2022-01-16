package business

import (
	"client/util"
)

type YearBalGateway struct {
}

func (yg *YearBalGateway) GetYearBalanceById(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.GetYearBalanceById_json(param); err != nil {
		errCode = util.ErrShowFailed
		logger.Error("the GetYearBalanceById_json failed,err:%v", err.Error())
	} else {
		logger.Debug("GetYearBalanceById_json succeed.")
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

func (yg *YearBalGateway) UpdateYearBalanceById(param []byte) (errCode int) {
	errCode = util.ErrNull
	if err := cSdk.UpdateYearBalanceById_json(param); err != nil {
		errCode = util.ErrUpdateFailed
		logger.Error("the UpdateYearBalanceById failed,err:%v", err.Error())
	} else {
		logger.Debug("UpdateYearBalanceById succeed")
	}
	return errCode
}

func (yg *YearBalGateway) DeleteYearBalanceByID(param []byte) (errCode int) {
	errCode = util.ErrNull
	if err := cSdk.DeleteYearBalanceByID_json(param); err != nil {
		errCode = util.ErrDeleteFailed
		logger.Error("the DeleteYearBalanceByID failed,err:%v", err.Error())
	} else {
		logger.Debug("DeleteYearBalanceByID succeed")
	}
	return errCode
}
