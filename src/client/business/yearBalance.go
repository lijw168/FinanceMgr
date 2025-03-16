package business

import (
	sdkUtil "financeMgr/src/analysis-server/sdk/util"
	"financeMgr/src/client/util"
)

type YearBalGateway struct {
}

func (yg *YearBalGateway) GetYearBalance(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.GetYearBalance_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrShowFailed
			errMsg = "GetYearBalance failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("GetYearBalance succeed.")
	}
	return resData, errCode, errMsg
}

func (yg *YearBalGateway) GetAccSubYearBalValue(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.GetAccSubYearBalValue_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrShowFailed
			errMsg = "GetAccSubYearBalValue failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("GetAccSubYearBalValue succeed.")
	}
	return resData, errCode, errMsg
}

func (yg *YearBalGateway) BatchCreateYearBalance(param []byte) (errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if err = cSdk.BatchCreateYearBalance_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrCreateFailed
			errMsg = "BatchCreateYearBalance failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("BatchCreateYearBalance succeed;")
	}
	return errCode, errMsg
}

func (yg *YearBalGateway) BatchDeleteYearBalance(param []byte) (errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if err = cSdk.BatchDeleteYearBalance_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrDeleteFailed
			errMsg = "BatchDeleteYearBalance failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("BatchDeleteYearBalance succeed;")
	}
	return errCode, errMsg
}

func (yg *YearBalGateway) UpdateYearBalance(param []byte) (errCode int, errMsg string) {
	errCode = util.ErrNull
	if err := cSdk.UpdateYearBalance_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrUpdateFailed
			errMsg = "UpdateYearBalance failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("UpdateYearBalance succeed")
	}
	return errCode, errMsg
}

func (yg *YearBalGateway) CreateYearBalance(param []byte) (errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if err = cSdk.CreateYearBalance_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrCreateFailed
			errMsg = "CreateYearBalance failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("CreateYearBalance succeed;")
	}
	return errCode, errMsg
}

func (yg *YearBalGateway) BatchUpdateBals(param []byte) (errCode int, errMsg string) {
	errCode = util.ErrNull
	if err := cSdk.BatchUpdateBals_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrUpdateFailed
			errMsg = "BatchUpdateBals failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("BatchUpdateBals succeed")
	}
	return errCode, errMsg
}

func (yg *YearBalGateway) DeleteYearBalance(param []byte) (errCode int, errMsg string) {
	errCode = util.ErrNull
	if err := cSdk.DeleteYearBalance_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrDeleteFailed
			errMsg = "DeleteYearBalance failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("DeleteYearBalanceByID succeed")
	}
	return errCode, errMsg
}

func (yg *YearBalGateway) AnnualClosing(param []byte) (errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if err = cSdk.AnnualClosing_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrAnnualClosing
			errMsg = "AnnualClosing failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("AnnualClosing succeed;")
	}
	return errCode, errMsg
}

func (yg *YearBalGateway) CancelAnnualClosing(param []byte) (errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if err = cSdk.CancelAnnualClosing_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrCancelAnnualClosing
			errMsg = "CancelAnnualClosing failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("CancelAnnualClosing succeed;")
	}
	return errCode, errMsg
}

func (yg *YearBalGateway) ListYearBalance(param []byte) (resData []byte, errCode int, errMsg string) {
	return listCmdJson(resource_type_year_balance, param, cSdk.ListYearBalance_json)
}

func (yg *YearBalGateway) GetAnnualClosingStatus(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.GetAnnualClosingStatus_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrGetAnnualClosingStatus
			errMsg = "GetAnnualClosingStatus failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("GetAnnualClosingStatus succeed.")
	}
	return resData, errCode, errMsg
}
