package business

import (
	//"financeMgr/src/analysis-server/model"
	"encoding/binary"
	"encoding/json"
	"financeMgr/src/analysis-server/sdk/options"
	sdkUtil "financeMgr/src/analysis-server/sdk/util"
	"financeMgr/src/client/util"
	"fmt"
)

type OperatorGateway struct {
}

func (og *OperatorGateway) ListOperatorInfo(param []byte) (resData []byte, errCode int, errMsg string) {
	return listCmdJson(resource_type_operator, param, cSdk.ListOperatorInfo_json)
}

func (og *OperatorGateway) GetOperatorInfo(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	id := int(binary.LittleEndian.Uint32(param))
	if id <= 0 {
		errMsg = fmt.Sprintf("the id param is: %d", id)
		logger.Error(errMsg)
		errCode = util.ErrInvalidParam
		return nil, errCode, errMsg
	}
	var opts options.BaseOptions
	opts.ID = id
	view, err := cSdk.GetOperatorInfo(&opts)
	if err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrShowFailed
			errMsg = "GetOperatorInfo failed,internal error"
		}
		logger.Error("GetOperatorInfo,err:%s", errMsg)
	} else {
		resData, err = json.Marshal(view)
		if err != nil {
			errCode = util.ErrMarshalFailed
			errMsg = fmt.Sprintf("json.Marshal failed, error:%s", err.Error())
			logger.Error("the Marshal failed,err:%s", errMsg)
		}
		logger.Debug("GetOperatorInfo succeed;views:%v", view)
	}
	return resData, errCode, errMsg
}

func (og *OperatorGateway) CreateOperator(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	if views, err := cSdk.CreateOperator_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrCreateFailed
			errMsg = "CreateOperator failed,internal error"
		}
		logger.Error("CreateOperator,err:%s", errMsg)
	} else {
		logger.Debug("CreateOperator succeed;views:%v", views)
		resData = make([]byte, 4)
		binary.LittleEndian.PutUint32(resData, uint32(views.OperatorID))
	}
	return
}

func (og *OperatorGateway) UpdateOperator(param []byte) (errCode int, errMsg string) {
	errCode = util.ErrNull
	if err := cSdk.UpdateOperator_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrUpdateFailed
			errMsg = "UpdateOperator failed,internal error"
		}
		logger.Error("UpdateOperator,err:%s", errMsg)
	} else {
		logger.Debug("UpdateOperator succeed")
	}
	return
}

func (og *OperatorGateway) DeleteOperator(param []byte) (errCode int, errMsg string) {
	id := int(binary.LittleEndian.Uint32(param))
	return deleteCmd(resource_type_operator, id, cSdk.DeleteOperator)
}
