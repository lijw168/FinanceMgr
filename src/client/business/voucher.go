package business

import (
	//"analysis-server/model"
	"analysis-server/sdk/options"
	"client/util"
	"encoding/binary"
	"encoding/json"
)

type VoucherGateway struct {
}

//beigin voucher
func (vg *VoucherGateway) CreateVoucher(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	if descData, err := cSdk.CreateVoucher_json(param); err != nil {
		errCode = util.ErrCreateFailed
		logger.Error("the CreateVoucher failed,err:%v", err.Error())
	} else {
		logger.Debug("CreateVoucher succeed;views:%v", descData)
		resData, err = json.Marshal(descData)
		if err != nil {
			errCode = util.ErrMarshalFailed
			logger.Error("the Marshal failed,err:%v", err.Error())
		}
	}
	return resData, errCode
}

func (vg *VoucherGateway) UpdateVoucher(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	var err error
	if resData, err = cSdk.UpdateVoucher_json(param); err != nil {
		errCode = util.ErrUpdateFailed
		logger.Error("the UpdateVoucher_json failed,err:%v", err.Error())
		resData = nil
	} else {
		logger.Debug("UpdateVoucher_json succeed")
	}
	return
}

func (vg *VoucherGateway) DeleteVoucher(param []byte) (errCode int) {
	id := int(binary.LittleEndian.Uint32(param))
	if id <= 0 {
		logger.Error("the id param is: %d", id)
		errCode = util.ErrInvalidParam
		return errCode
	}
	return deleteCmd(resource_type_voucher, id, cSdk.DeleteVoucher)
}

func (vg *VoucherGateway) GetVoucher(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	id := int(binary.LittleEndian.Uint32(param))
	if id <= 0 {
		logger.Error("the id param is: %d", id)
		errCode = util.ErrInvalidParam
		return nil, errCode
	}
	var opts options.BaseOptions
	opts.ID = id
	view, err := cSdk.GetVoucher(&opts)
	if err != nil {
		errCode = util.ErrShowFailed
		logger.Error("the GetVoucher failed,err:%v", err.Error())
	} else {
		resData, err = json.Marshal(view)
		if err != nil {
			errCode = util.ErrMarshalFailed
			logger.Error("the Marshal failed,err:%v", err.Error())
		}
		logger.Debug("GetVoucher succeed;views:%v", view)
	}
	return resData, errCode
}

func (vg *VoucherGateway) VoucherAudit(param []byte) (errCode int) {
	errCode = util.ErrNull
	if err := cSdk.VoucherAudit_json(param); err != nil {
		errCode = util.ErrVoucherAuditFailed
		logger.Error("the VoucherAudit_json failed,err:%v", err.Error())
	} else {
		logger.Debug("VoucherAudit_json succeed")
	}
	return errCode
}

//end voucher;

//beigin voucher records
func (vg *VoucherGateway) CreateVoucherRecords(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	if descData, err := cSdk.CreateVoucherRecords_json(param); err != nil {
		errCode = util.ErrCreateFailed
		logger.Error("the CreateVoucherRecords_json failed,err:%v", err.Error())
	} else {
		logger.Debug("CreateVoucherRecords_json succeed;descData:%v", descData)
		resData, err = json.Marshal(descData)
		if err != nil {
			errCode = util.ErrMarshalFailed
			logger.Error("the Marshal failed,err:%v", err.Error())
		}
	}
	return resData, errCode
}

func (vg *VoucherGateway) DeleteVoucherRecord(param []byte) (errCode int) {
	id := int(binary.LittleEndian.Uint32(param))
	if id <= 0 {
		logger.Error("the id param is: %d", id)
		errCode = util.ErrInvalidParam
		return errCode
	}
	return deleteCmd(resource_type_voucher_record, id, cSdk.DeleteVoucherRecord)
}

func (vg *VoucherGateway) DeleteVoucherRecords(param []byte) (errCode int) {
	errCode = util.ErrNull
	if err := cSdk.DeleteVoucherRecords_json(param); err != nil {
		errCode = util.ErrDeleteFailed
		logger.Error("the DeleteVoucherRecords_json failed,err:%v", err.Error())
	} else {
		logger.Debug("DeleteVoucherRecords_json succeed")
	}
	return
}

func (vg *VoucherGateway) ListVoucherRecords(param []byte) (resData []byte, errCode int) {
	return listCmdJson(resource_type_voucher_record, param, cSdk.ListVoucherRecords_json)
}

func (vg *VoucherGateway) UpdateVoucherRecordByID(param []byte) (errCode int) {
	errCode = util.ErrNull
	if err := cSdk.UpdateVoucherRecord_json(param); err != nil {
		errCode = util.ErrUpdateFailed
		logger.Error("the UpdateVoucherRecord_json failed,err:%v", err.Error())
	} else {
		logger.Debug("UpdateVoucherRecord_json succeed")
	}
	return errCode
}

//end voucher records

//begin voucher info
func (vg *VoucherGateway) GetVoucherInfo(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	id := int(binary.LittleEndian.Uint32(param))
	if id <= 0 {
		logger.Error("the id param is: %d", id)
		errCode = util.ErrInvalidParam
		return nil, errCode
	}
	var opts options.BaseOptions
	opts.ID = id
	view, err := cSdk.GetVoucherInfo(&opts)
	if err != nil {
		errCode = util.ErrShowFailed
		logger.Error("the GetVoucherInfo failed,err:%v", err.Error())
	} else {
		resData, err = json.Marshal(view)
		if err != nil {
			errCode = util.ErrMarshalFailed
			logger.Error("the Marshal failed,err:%v", err.Error())
		}
		logger.Debug("GetVoucherInfo succeed;views:%v", view)
	}
	return resData, errCode
}

func (vg *VoucherGateway) GetLatestVoucherInfo(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	id := int(binary.LittleEndian.Uint32(param))
	if id <= 0 {
		logger.Error("the id param is: %d", id)
		errCode = util.ErrInvalidParam
		return nil, errCode
	}
	var opts options.BaseOptions
	opts.ID = id
	inputParam, err := json.Marshal(opts)
	if err != nil {
		errCode = util.ErrMarshalFailed
		logger.Error("the Marshal failed,err:%v", err.Error())
		return nil, errCode
	}
	resData, err = cSdk.GetLatestVoucherInfo_json(inputParam)
	if err != nil {
		errCode = util.ErrGetLatestVoucherInfoFailed
		logger.Error("the GetLatestVoucherInfo failed,err:%v", err.Error())
	} else {
		logger.Debug("GetLatestVoucherInfo succeed.")
	}
	return resData, errCode
}

func (vg *VoucherGateway) GetMaxNumOfMonth(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	if iCount, err := cSdk.GetMaxNumOfMonth_json(param); err != nil {
		errCode = util.ErrGetMaxNumOfMonthFailed
		logger.Error("the GetMaxNumOfMonth_json failed,err:%v", err.Error())
		resData = nil
		return
	} else {
		logger.Debug("GetMaxNumOfMonth_json succeed")
		resData = make([]byte, 4)
		binary.LittleEndian.PutUint32(resData, uint32(iCount))
	}
	return
}

func (vg *VoucherGateway) ListVoucherInfo(param []byte) (resData []byte, errCode int) {
	return listCmdJson(resource_type_voucher_info, param, cSdk.ListVoucherInfo_json)
}

func (vg *VoucherGateway) ListVoucherInfoByMulCondition(param []byte) (resData []byte, errCode int) {
	return listCmdJson(resource_type_voucher_info, param, cSdk.ListVoucherInfoByMulCondition_json)
}
