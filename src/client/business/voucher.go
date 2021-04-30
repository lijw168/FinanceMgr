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
func (vg *VoucherGateway) CreateVoucher(param []byte) (resData []byte, errCode int) {
	var opts options.VoucherOptions
	// var createRecOpt options.CreateVoucherRecordOptions
	// opts.RecordsOptions = append(opts.RecordsOptions, createRecOpt)
	errCode = util.ErrNull
	if err := json.Unmarshal(param, &opts); err != nil {
		logger.Error("the Unmarshal failed,err:%v", err.Error())
		errCode = util.ErrUnmarshalFailed
		return nil, errCode
	}
	if descData, err := cSdk.CreateVoucher(&opts); err != nil {
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
func (vg *VoucherGateway) CreateVoucherRecords(param []byte) (resData []byte, errCode int) {
	optSlice := []options.CreateVoucherRecordOptions{}
	// var opts options.CreateVoucherRecordOptions
	// optSlice = append(optSlice, opts)
	errCode = util.ErrNull
	if err := json.Unmarshal(param, &optSlice); err != nil {
		logger.Error("the Unmarshal failed,err:%v", err.Error())
		errCode = util.ErrUnmarshalFailed
		return nil, errCode
	}
	if descData, err := cSdk.CreateVoucherRecords(optSlice); err != nil {
		errCode = util.ErrCreateFailed
		logger.Error("the CreateVoucherRecords failed,err:%v", err.Error())
	} else {
		logger.Debug("CreateVoucherRecords succeed;descData:%v", descData)
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
	var opts options.ListOptions
	errCode = util.ErrNull
	if err := json.Unmarshal(param, &opts); err != nil {
		logger.Error("the Unmarshal failed,err:%v", err.Error())
		errCode = util.ErrUnmarshalFailed
		return nil, errCode
	}
	if count, views, err := cSdk.ListVoucherRecords(&opts); err != nil {
		logger.Error("the ListVoucherRecords failed,err:%v", err.Error())
	} else {
		logger.Debug("ListVoucherRecords succeed;count:%d,views:%v", count, views)
		desc := &(model.DescData{})
		desc.Tc = count
		desc.Elements = views
		resData, err = json.Marshal(desc)
		if err != nil {
			errCode = util.ErrMarshalFailed
			logger.Error("the Marshal failed,err:%v", err.Error())
		}
	}
	return resData, errCode
}

func (vg *VoucherGateway) UpdateVoucherRecord(param []byte) (errCode int) {
	var opts options.ModifyVoucherRecordOptions
	errCode = util.ErrNull
	if err := json.Unmarshal(param, &opts); err != nil {
		logger.Error("the Unmarshal failed,err:%v", err.Error())
		errCode = util.ErrUnmarshalFailed
		return errCode
	}
	if err := cSdk.UpdateVoucherRecord(&opts); err != nil {
		errCode = util.ErrUpdateFailed
		logger.Error("the UpdateVoucherRecord failed,err:%v", err.Error())
	} else {
		logger.Debug("UpdateVoucherRecord succeed")
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
	var opts options.ListOptions
	errCode = util.ErrNull
	if err := json.Unmarshal(param, &opts); err != nil {
		logger.Error("the Unmarshal failed,err:%v", err.Error())
		errCode = util.ErrUnmarshalFailed
		return nil, errCode
	}
	if count, views, err := cSdk.ListVoucherInfo(&opts); err != nil {
		logger.Error("the ListVoucherInfo failed,err:%v", err.Error())
	} else {
		logger.Debug("ListVoucherInfo succeed;count:%d,views:%v", count, views)
		desc := &(model.DescData{})
		desc.Tc = count
		desc.Elements = views
		resData, err = json.Marshal(desc)
		if err != nil {
			errCode = util.ErrMarshalFailed
			logger.Error("the Marshal failed,err:%v", err.Error())
		}
	}
	return resData, errCode
}
