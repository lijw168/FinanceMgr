package business

import (
	sdkUtil "financeMgr/src/analysis-server/sdk/util"
	"financeMgr/src/client/util"
)

type YearBalGateway struct {
}

func (yg *YearBalGateway) GetYearBalance(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.GetYearBalance_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
		} else {
			errCode = util.ErrShowFailed
		}
		logger.Error("the GetYearBalance failed,err:%v", err.Error())
	} else {
		logger.Debug("GetYearBalance succeed.")
	}
	return resData, errCode
}

func (yg *YearBalGateway) GetAccSubYearBalValue(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.GetAccSubYearBalValue_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
		} else {
			errCode = util.ErrShowFailed
		}
		logger.Error("the GetAccSubYearBalValue failed,err:%v", err.Error())
	} else {
		logger.Debug("GetAccSubYearBalValue succeed.")
	}
	return resData, errCode
}

func (yg *YearBalGateway) CreateYearBalance(param []byte) (errCode int) {
	errCode = util.ErrNull
	var err error
	if err = cSdk.CreateYearBalance_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
		} else {
			errCode = util.ErrCreateFailed
		}
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
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
		} else {
			errCode = util.ErrCreateFailed
		}
		logger.Error("the BatchCreateYearBalance failed,err:%v", err.Error())
	} else {
		logger.Debug("BatchCreateYearBalance succeed;")
	}
	return errCode
}

func (yg *YearBalGateway) BatchDeleteYearBalance(param []byte) (errCode int) {
	errCode = util.ErrNull
	var err error
	if err = cSdk.BatchDeleteYearBalance_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
		} else {
			errCode = util.ErrDeleteFailed
		}
		logger.Error("the BatchDeleteYearBalance failed,err:%v", err.Error())
	} else {
		logger.Debug("BatchDeleteYearBalance succeed;")
	}
	return errCode
}

func (yg *YearBalGateway) UpdateYearBalance(param []byte) (errCode int) {
	errCode = util.ErrNull
	if err := cSdk.UpdateYearBalance_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
		} else {
			errCode = util.ErrUpdateFailed
		}
		logger.Error("the UpdateYearBalance failed,err:%v", err.Error())
	} else {
		logger.Debug("UpdateYearBalance succeed")
	}
	return errCode
}

func (yg *YearBalGateway) BatchUpdateBals(param []byte) (errCode int) {
	errCode = util.ErrNull
	if err := cSdk.BatchUpdateBals_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
		} else {
			errCode = util.ErrUpdateFailed
		}
		logger.Error("the BatchUpdateBals failed,err:%v", err.Error())
	} else {
		logger.Debug("BatchUpdateBals succeed")
	}
	return errCode
}

func (yg *YearBalGateway) DeleteYearBalance(param []byte) (errCode int) {
	errCode = util.ErrNull
	if err := cSdk.DeleteYearBalance_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
		} else {
			errCode = util.ErrDeleteFailed
		}
		logger.Error("the DeleteYearBalanceByID failed,err:%v", err.Error())
	} else {
		logger.Debug("DeleteYearBalanceByID succeed")
	}
	return errCode
}

func (yg *YearBalGateway) ListYearBalance(param []byte) (resData []byte, errCode int) {
	return listCmdJson(resource_type_year_balance, param, cSdk.ListYearBalance_json)
}

func (yg *YearBalGateway) AnnualClosing(param []byte) (errCode int) {
	errCode = util.ErrNull
	var err error
	if err = cSdk.AnnualClosing_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
		} else {
			errCode = util.ErrAnnualClosing
		}
		logger.Error("the AnnualClosing failed,err:%v", err.Error())
	} else {
		logger.Debug("AnnualClosing succeed;")
	}
	return errCode
}

func (yg *YearBalGateway) CancelAnnualClosing(param []byte) (errCode int) {
	errCode = util.ErrNull
	var err error
	if err = cSdk.CancelAnnualClosing_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
		} else {
			errCode = util.ErrCancelAnnualClosing
		}
		logger.Error("the CancelAnnualClosing failed,err:%v", err.Error())
	} else {
		logger.Debug("CancelAnnualClosing succeed;")
	}
	return errCode
}

func (yg *YearBalGateway) GetAnnualClosingStatus(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.GetAnnualClosingStatus_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
		} else {
			errCode = util.ErrGetAnnualClosingStatus
		}
		logger.Error("the GetAnnualClosingStatus failed,err:%v", err.Error())
	} else {
		logger.Debug("GetAnnualClosingStatus succeed.")
	}
	return resData, errCode
}
