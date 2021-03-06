package business

import (
	"analysis-server/model"
	"analysis-server/sdk/options"
	"client/util"
	"encoding/binary"
	"encoding/json"
)

type VoucherGateway struct {
}

//beigin voucher
// func (vg *VoucherGateway) CreateVoucher(param []byte) (resData []byte, errCode int) {
// 	//由于从界面进程传该操作的参数是json格式的，并且为了不修改cSdk里的函数代码，
// 	//所以先需要通过VoucherParams获取相应的参数值。
// 	var opts options.VoucherOptions
// 	// var createRecOpt options.CreateVoucherRecordOptions
// 	// opts.RecordsOptions = append(opts.RecordsOptions, createRecOpt)
// 	errCode = util.ErrNull
// 	if err := json.Unmarshal(param, &opts); err != nil {
// 		logger.Error("the Unmarshal failed,err:%v", err.Error())
// 		errCode = util.ErrUnmarshalFailed
// 		return nil, errCode
// 	}
// 	if descData, err := cSdk.CreateVoucher(&opts); err != nil {
// 		errCode = util.ErrCreateFailed
// 		logger.Error("the CreateVoucher failed,err:%v", err.Error())
// 	} else {
// 		logger.Debug("CreateVoucher succeed;views:%v", descData)
// 		resData, err = json.Marshal(descData)
// 		if err != nil {
// 			errCode = util.ErrMarshalFailed
// 			logger.Error("the Marshal failed,err:%v", err.Error())
// 		}
// 	}
// 	return resData, errCode
// }

func (vg *VoucherGateway) CreateVoucher(param []byte) (resData []byte, errCode int) {
	var paramsOpts model.VoucherParams
	errCode = util.ErrNull
	if err := json.Unmarshal(param, &paramsOpts); err != nil {
		logger.Error("the Unmarshal failed,err:%v", err.Error())
		errCode = util.ErrUnmarshalFailed
		return nil, errCode
	}
	if descData, err := cSdk.CreateVoucher_json(&paramsOpts); err != nil {
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

func (vg *VoucherGateway) DeleteVoucher(param []byte) (errCode int) {
	id := int(binary.LittleEndian.Uint32(param))
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

//end voucher;
//beigin voucher records
// func (vg *VoucherGateway) CreateVoucherRecords(param []byte) (resData []byte, errCode int) {
// 	optSlice := []options.CreateVoucherRecordOptions{}
// 	// var opts options.CreateVoucherRecordOptions
// 	// optSlice = append(optSlice, opts)
// 	errCode = util.ErrNull
// 	if err := json.Unmarshal(param, &optSlice); err != nil {
// 		logger.Error("the Unmarshal failed,err:%v", err.Error())
// 		errCode = util.ErrUnmarshalFailed
// 		return nil, errCode
// 	}
// 	if descData, err := cSdk.CreateVoucherRecords(optSlice); err != nil {
// 		errCode = util.ErrCreateFailed
// 		logger.Error("the CreateVoucherRecords failed,err:%v", err.Error())
// 	} else {
// 		logger.Debug("CreateVoucherRecords succeed;descData:%v", descData)
// 		resData, err = json.Marshal(descData)
// 		if err != nil {
// 			errCode = util.ErrMarshalFailed
// 			logger.Error("the Marshal failed,err:%v", err.Error())
// 		}
// 	}
// 	return resData, errCode
// }

func (vg *VoucherGateway) CreateVoucherRecords(param []byte) (resData []byte, errCode int) {
	recordParamSlice := []*model.CreateVoucherRecordParams{}
	errCode = util.ErrNull
	if err := json.Unmarshal(param, &recordParamSlice); err != nil {
		logger.Error("the Unmarshal failed,err:%v", err.Error())
		errCode = util.ErrUnmarshalFailed
		return nil, errCode
	}

	if descData, err := cSdk.CreateVoucherRecords_json(recordParamSlice); err != nil {
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
	return deleteCmd(resource_type_voucher_record, id, cSdk.DeleteVoucherRecord)
}

func (vg *VoucherGateway) ListVoucherRecords(param []byte) (resData []byte, errCode int) {
	return listCmdJson(resource_type_voucher_record, param, cSdk.ListVoucherRecords_json)
}

func (vg *VoucherGateway) UpdateVoucherRecord(param []byte) (errCode int) {
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

func (vg *VoucherGateway) ListVoucherInfo(param []byte) (resData []byte, errCode int) {
	return listCmdJson(resource_type_voucher_info, param, cSdk.ListVoucherInfo_json)
}
