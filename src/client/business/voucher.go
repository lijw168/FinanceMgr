package business

import (
	//"financeMgr/src/analysis-server/model"
	"encoding/binary"
	"financeMgr/src/analysis-server/sdk/options"
	sdkUtil "financeMgr/src/analysis-server/sdk/util"
	"financeMgr/src/client/util"
	//"encoding/json"
)

type VoucherGateway struct {
}

// beigin voucher
func (vg *VoucherGateway) CreateVoucher(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.CreateVoucher_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrCreateFailed
			errMsg = "CreateVoucher failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("CreateVoucher succeed;views.")
	}
	return resData, errCode, errMsg
}

func (vg *VoucherGateway) UpdateVoucher(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.UpdateVoucher_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrUpdateFailed
			errMsg = "UpdateVoucher_json failed,internal error"
		}
		logger.LogError(errMsg)
		resData = nil
	} else {
		logger.Debug("UpdateVoucher_json succeed")
	}
	return resData, errCode, errMsg
}

func (vg *VoucherGateway) GetVoucher(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.GetVoucher_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrShowFailed
			errMsg = "ErrShowFailed failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("GetVoucher_json succeed.")
	}
	return resData, errCode, errMsg
}

func (vg *VoucherGateway) DeleteVoucher(param []byte) (errCode int, errMsg string) {
	errCode = util.ErrNull
	if err := cSdk.DeleteVoucher_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrDeleteFailed
			errMsg = "DeleteVoucher_json failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("DeleteVoucher_json succeed")
	}
	return
}

func (vg *VoucherGateway) ArrangeVoucher(param []byte) (errCode int, errMsg string) {
	errCode = util.ErrNull
	if err := cSdk.ArrangeVoucher_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrVoucherArrangeFailed
			errMsg = "ArrangeVoucher_json failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("ArrangeVoucher_json succeed")
	}
	return
}

//end voucher;

//beigin voucher records
// func (vg *VoucherGateway) CreateVoucherRecords(param []byte) (resData []byte, errCode int) {
// 	errCode = util.ErrNull
// 	if descData, err := cSdk.CreateVoucherRecords_json(param); err != nil {
// 		errCode = util.ErrCreateFailed
// 		logger.Error("the CreateVoucherRecords_json failed,err:%v", err.Error())
// 	} else {
// 		logger.Debug("CreateVoucherRecords_json succeed;descData:%v", descData)
// 		resData, err = json.Marshal(descData)
// 		if err != nil {
// 			errCode = util.ErrMarshalFailed
// 			logger.Error("the Marshal failed,err:%v", err.Error())
// 		}
// 	}
// 	return resData, errCode
// }

// func (vg *VoucherGateway) DeleteVoucherRecord(param []byte) (errCode int) {
// 	id := int(binary.LittleEndian.Uint32(param))
// 	if id <= 0 {
// 		logger.Error("the id param is: %d", id)
// 		errCode = util.ErrInvalidParam
// 		return errCode
// 	}
// 	return deleteCmd(resource_type_voucher_record, id, cSdk.DeleteVoucherRecord)
// }

// func (vg *VoucherGateway) DeleteVoucherRecords(param []byte) (errCode int) {
// 	errCode = util.ErrNull
// 	if err := cSdk.DeleteVoucherRecords_json(param); err != nil {
// 		errCode = util.ErrDeleteFailed
// 		logger.Error("the DeleteVoucherRecords_json failed,err:%v", err.Error())
// 	} else {
// 		logger.Debug("DeleteVoucherRecords_json succeed")
// 	}
// 	return
// }

func (vg *VoucherGateway) ListVoucherRecords(param []byte) (resData []byte, errCode int, errMsg string) {
	return listCmdJson(resource_type_voucher_record, param, cSdk.ListVoucherRecords_json)
}

// func (vg *VoucherGateway) UpdateVoucherRecordByID(param []byte) (errCode int) {
// 	errCode = util.ErrNull
// 	if err := cSdk.UpdateVoucherRecord_json(param); err != nil {
// 		errCode = util.ErrUpdateFailed
// 		logger.Error("the UpdateVoucherRecord_json failed,err:%v", err.Error())
// 	} else {
// 		logger.Debug("UpdateVoucherRecord_json succeed")
// 	}
// 	return errCode
// }

//end voucher records

// begin voucher info
func (vg *VoucherGateway) GetVoucherInfo(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.GetVoucherInfo_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrShowFailed
			errMsg = "GetVoucherInfo_json failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("GetVoucherInfo_json succeed.")
	}
	return resData, errCode, errMsg
}

func (vg *VoucherGateway) GetLatestVoucherInfo(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.GetLatestVoucherInfo_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrGetLatestVoucherInfoFailed
			errMsg = "GetLatestVoucherInfo_json failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("GetLatestVoucherInfo_json succeed.")
	}
	return resData, errCode, errMsg
}

func (vg *VoucherGateway) GetMaxNumOfMonth(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	if iCount, err := cSdk.GetMaxNumOfMonth_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrGetMaxNumOfMonthFailed
			errMsg = "GetMaxNumOfMonth_json failed,internal error"
		}
		logger.LogError(errMsg)
		resData = nil
	} else {
		logger.Debug("GetMaxNumOfMonth_json succeed")
		resData = make([]byte, 4)
		binary.LittleEndian.PutUint32(resData, uint32(iCount))
	}
	return resData, errCode, errMsg
}

func (vg *VoucherGateway) ListVoucherInfo(param []byte) (resData []byte, errCode int, errMsg string) {
	return listCmdJson(resource_type_voucher_info, param, cSdk.ListVoucherInfo_json)
}

func (vg *VoucherGateway) ListVoucherInfoWithAuxCondition(param []byte) (resData []byte, errCode int, errMsg string) {
	return listCmdJson(resource_type_voucher_info, param, cSdk.ListVoucherInfoWithAuxCondition_json)
}

func (vg *VoucherGateway) UpdateVoucherInfo(param []byte) (errCode int, errMsg string) {
	errCode = util.ErrNull
	if err := cSdk.UpdateVoucherInfo_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrUpdateFailed
			errMsg = "UpdateVoucherInfo_json failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("UpdateVoucherInfo_json succeed")
	}
	return
}

func (vg *VoucherGateway) BatchAuditVouchers(param []byte) (errCode int, errMsg string) {
	errCode = util.ErrNull
	if err := cSdk.BatchAuditVouchers_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrBatchAuditVouchersFailed
			errMsg = "BatchAuditVouchers_json failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("BatchAuditVouchers_json succeed")
	}
	return
}

// beigin voucher template
func (vg *VoucherGateway) CreateVoucherTemplate(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	var iSerialNum int
	if iSerialNum, err = cSdk.CreateVoucherTemplate_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrCreateFailed
			errMsg = "CreateVoucherTemplate_json failed, internal error"
		}
		logger.LogError(errMsg)
		resData = nil
	} else {
		resData = make([]byte, 4)
		binary.LittleEndian.PutUint32(resData, uint32(iSerialNum))
		logger.Debug("CreateVoucherTemplate_json succeed.")
	}
	return resData, errCode, errMsg
}

func (vg *VoucherGateway) ListVoucherTemplate(param []byte) (resData []byte, errCode int, errMsg string) {
	return listCmdJson(resource_type_voucher_template, param, cSdk.ListVoucherTemplate_json)
}

func (vg *VoucherGateway) DeleteVoucherTemplate(param []byte) (errCode int, errMsg string) {
	id := int(binary.LittleEndian.Uint32(param))
	return deleteCmd(resource_type_voucher_template, id, cSdk.DeleteVoucherTemplate)
}

func (vg *VoucherGateway) GetVoucherTemplate(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	id := int(binary.LittleEndian.Uint32(param))
	if id <= 0 {
		errCode = util.ErrInvalidParam
		errMsg = "Invalid voucher template ID"
		logger.LogError(errMsg)
		return nil, errCode, errMsg
	}
	var opts options.BaseOptions
	opts.ID = id
	var err error
	if resData, err = cSdk.GetVoucherTemplate(&opts); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrShowFailed
			errMsg = "GetVoucherTemplate failed, internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("GetVoucherTemplate succeed.")
	}
	return resData, errCode, errMsg
}

func (vg *VoucherGateway) CalculateAccumulativeMoney(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.CalculateAccumulativeMoney_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrCalAccuMoneyFailed
			errMsg = "CalculateAccumulativeMoney_json failed, internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("CalculateAccumulativeMoney_json succeed.")
	}
	return resData, errCode, errMsg
}

func (vg *VoucherGateway) BatchCalcAccuMoney(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.BatchCalcAccuMoney_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrBatchCalAccuMoneyFailed
			errMsg = "BatchCalcAccuMoney_json failed, internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("BatchCalcAccuMoney_json succeed.")
	}
	return resData, errCode, errMsg
}

func (vg *VoucherGateway) CalcAccountOfPeriod(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.CalcAccountOfPeriod_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrCalAccPeriodFailed
			errMsg = "CalcAccountOfPeriod_json failed, internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("CalcAccountOfPeriod_json succeed.")
	}
	return resData, errCode, errMsg
}

func (vg *VoucherGateway) GetNoAuditedVoucherInfoCount(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	if iCount, err := cSdk.GetNoAuditedVoucherInfoCount_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrGetMaxNumOfMonthFailed
			errMsg = "GetNoAuditedVoucherInfoCount_json failed,internal error"
		}
		logger.LogError(errMsg)
		resData = nil
	} else {
		logger.Debug("GetNoAuditedVoucherInfoCount_json succeed")
		resData = make([]byte, 4)
		binary.LittleEndian.PutUint32(resData, uint32(iCount))
	}
	return resData, errCode, errMsg
}

// voucher report, end
