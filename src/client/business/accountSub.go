package business

import (
	//"analysis-server/model"
	"analysis-server/sdk/options"
	"client/util"
	"encoding/binary"
	"encoding/json"
)

type AccSubGateway struct {
}

func (ag *AccSubGateway) ListAccSub(param []byte) (resData []byte, errCode int) {
	return listCmdJson(resource_type_account_sub, param, cSdk.ListAccSub_json)
}

func (ag *AccSubGateway) GetAccSub(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	id := int(binary.LittleEndian.Uint32(param))
	if id <= 0 {
		logger.Error("the id param is: %d", id)
		errCode = util.ErrInvalidParam
		return nil, errCode
	}
	var opts options.BaseOptions
	opts.ID = id
	view, err := cSdk.GetAccSub(&opts)
	if err != nil {
		errCode = util.ErrShowFailed
		logger.Error("the GetAccSub failed,err:%v", err.Error())
	} else {
		resData, err = json.Marshal(view)
		if err != nil {
			errCode = util.ErrMarshalFailed
			logger.Error("the Marshal failed,err:%v", err.Error())
		}
		logger.Debug("GetAccSub succeed;views:%v", view)
	}
	return resData, errCode
}

func (ag *AccSubGateway) CreateAccSub(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull

	if views, err := cSdk.CreateAccSub_json(param); err != nil {
		errCode = util.ErrCreateFailed
		logger.Error("the CreateAccSub failed,err:%v", err.Error())
	} else {
		logger.Debug("CreateAccSub succeed;views:%v", views)
		resData = make([]byte, 4)
		binary.LittleEndian.PutUint32(resData, uint32(views.SubjectID))
	}
	return resData, errCode
}

func (ag *AccSubGateway) UpdateAccSub(param []byte) (errCode int) {
	errCode = util.ErrNull
	if err := cSdk.UpdateAccSub_json(param); err != nil {
		errCode = util.ErrUpdateFailed
		logger.Error("the UpdateAccSub failed,err:%v", err.Error())
	} else {
		logger.Debug("UpdateAccSub succeed")
	}
	return errCode
}

func (ag *AccSubGateway) DeleteAccSub(param []byte) (errCode int) {
	id := int(binary.LittleEndian.Uint32(param))
	if id <= 0 {
		logger.Error("the id param is: %d", id)
		errCode = util.ErrInvalidParam
		return errCode
	}
	return deleteCmd(resource_type_account_sub, id, cSdk.DeleteAccSub)
}

func (ag *AccSubGateway) QueryAccSubReference(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	id := int(binary.LittleEndian.Uint32(param))
	if id <= 0 {
		logger.Error("the id param is: %d", id)
		errCode = util.ErrInvalidParam
		return nil, errCode
	}
	var opts options.BaseOptions
	opts.ID = id
	if iRefCount, err := cSdk.QueryAccSubReference(&opts); err != nil {
		errCode = util.ErrAccSubRefQueryFailed
		logger.Error("the QueryAccSubReference failed,err:%v", err.Error())
	} else {
		logger.Debug("QueryAccSubReference succeed;iRefCount:%v", iRefCount)
		resData = make([]byte, 4)
		binary.LittleEndian.PutUint32(resData, uint32(iRefCount))
	}
	return resData, errCode
}
