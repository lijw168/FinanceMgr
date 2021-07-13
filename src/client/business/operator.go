package business

import (
	//"analysis-server/model"
	"analysis-server/sdk/options"
	"client/util"
	"encoding/binary"
	"encoding/json"
)

type OperatorGateway struct {
}

func (og *OperatorGateway) ListOperatorInfo(param []byte) (resData []byte, errCode int) {
	return listCmdJson(resource_type_operator, param, cSdk.ListOperatorInfo_json)
}

func (og *OperatorGateway) GetOperatorInfo(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	id := int(binary.LittleEndian.Uint32(param))
	if id <= 0 {
		logger.Error("the id param is: %d", id)
		errCode = util.ErrInvalidParam
		return nil, errCode
	}
	var opts options.BaseOptions
	opts.ID = id
	view, err := cSdk.GetOperatorInfo(&opts)
	if err != nil {
		errCode = util.ErrShowFailed
		logger.Error("the GetOperatorInfo failed,err:%v", err.Error())
	} else {
		resData, err = json.Marshal(view)
		if err != nil {
			errCode = util.ErrMarshalFailed
			logger.Error("the Marshal failed,err:%v", err.Error())
		}
		logger.Debug("GetOperatorInfo succeed;views:%v", view)
	}
	return resData, errCode
}

func (og *OperatorGateway) CreateOperator(param []byte) (errCode int) {
	errCode = util.ErrNull
	if views, err := cSdk.CreateOperator_json(param); err != nil {
		errCode = util.ErrCreateFailed
		logger.Error("the CreateOperator failed,err:%v", err.Error())
	} else {
		logger.Debug("CreateOperator succeed;views:%v", views)
	}
	return errCode
}

func (og *OperatorGateway) UpdateOperator(param []byte) (errCode int) {
	errCode = util.ErrNull
	if err := cSdk.UpdateOperator_json(param); err != nil {
		errCode = util.ErrUpdateFailed
		logger.Error("the UpdateOperator failed,err:%v", err.Error())
	} else {
		logger.Debug("UpdateOperator succeed")
	}
	return errCode
}

func (og *OperatorGateway) DeleteOperator(param []byte) (errCode int) {
	id := int(binary.LittleEndian.Uint32(param))
	return deleteCmd(resource_type_operator, id, cSdk.DeleteOperator)
}
