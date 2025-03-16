package business

import (
	//"financeMgr/src/analysis-server/model"
	"encoding/binary"
	"encoding/json"
	"financeMgr/src/analysis-server/sdk/options"
	sdkUtil "financeMgr/src/analysis-server/sdk/util"
	"financeMgr/src/client/util"
	"fmt"
	//"math"
)

type AccSubGateway struct {
}

func (ag *AccSubGateway) ListAccSub(param []byte) (resData []byte, errCode int, errMsg string) {
	return listCmdJson(resource_type_account_sub, param, cSdk.ListAccSub_json)
}

func (ag *AccSubGateway) GetAccSub(param []byte) (resData []byte, errCode int, errMsg string) {
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
	view, err := cSdk.GetAccSub(&opts)
	if err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrShowFailed
			errMsg = "GetAccSub failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		resData, err = json.Marshal(view)
		if err != nil {
			errCode = util.ErrMarshalFailed
			errMsg = fmt.Sprintf("the Marshal failed,err:%v", err.Error())
			logger.LogError(errMsg)
		} else {
			logger.Debug("GetAccSub succeed;views:%v", view)
		}
	}
	return resData, errCode, errMsg
}

func (ag *AccSubGateway) CreateAccSub(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull

	if views, err := cSdk.CreateAccSub_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrCreateFailed
			errMsg = "CreateAccSub failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("CreateAccSub succeed;views:%v", views)
		resData = make([]byte, 4)
		binary.LittleEndian.PutUint32(resData, uint32(views.SubjectID))
	}
	return resData, errCode, errMsg
}

func (ag *AccSubGateway) UpdateAccSub(param []byte) (errCode int, errMsg string) {
	errCode = util.ErrNull
	if err := cSdk.UpdateAccSub_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrUpdateFailed
			errMsg = "UpdateAccSub failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("UpdateAccSub succeed")
	}
	return errCode, errMsg
}

func (ag *AccSubGateway) DeleteAccSub(param []byte) (errCode int, errMsg string) {
	id := int(binary.LittleEndian.Uint32(param))
	if id <= 0 {
		errMsg = fmt.Sprintf("the id param is: %d", id)
		logger.Error(errMsg)
		errCode = util.ErrInvalidParam
		return errCode, errMsg
	}
	return deleteCmd(resource_type_account_sub, id, cSdk.DeleteAccSub)
}

func (ag *AccSubGateway) QueryAccSubReference(param []byte) (resData []byte, errCode int, errMsg string) {
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
	if iRefCount, err := cSdk.QueryAccSubReference(&opts); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrAccSubRefQueryFailed
			errMsg = "QueryAccSubReference failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("QueryAccSubReference succeed;iRefCount:%v", iRefCount)
		resData = make([]byte, 4)
		binary.LittleEndian.PutUint32(resData, uint32(iRefCount))
	}
	return resData, errCode, errMsg
}

func (ag *AccSubGateway) CopyAccSubTemplate(param []byte) (resData []byte, errCode int, errMsg string) {
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
	var err error
	if resData, err = cSdk.CopyAccSubTemplate(&opts); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrCopyAccSubTemplateFailed
			errMsg = "CopyAccSubTemplate failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("CopyAccSubTemplate succeed")
	}
	return resData, errCode, errMsg
}

// func (ag *AccSubGateway) UpdateYearBalance(param []byte) (errCode int) {
// 	errCode = util.ErrNull
// 	if err := cSdk.UpdateYearBalance_json(param); err != nil {
// 		errCode = util.ErrUpdateFailed
// 		logger.Error("the UpdateYearBalance failed,err:%v", err.Error())
// 	} else {
// 		logger.Debug("UpdateYearBalance succeed")
// 	}
// 	return errCode
// }

// func (ag *AccSubGateway) GetYearBalance(param []byte) (resData []byte, errCode int) {
// 	errCode = util.ErrNull
// 	var err error
// 	if resData, err = cSdk.GetYearBalanceById_json(param); err != nil {
// 		errCode = util.ErrShowFailed
// 		logger.Error("the GetYearBalanceById_json failed,err:%v", err.Error())
// 	} else {
// 		logger.Debug("GetYearBalanceById_json succeed.")
// 	}
// 	return resData, errCode
// }
