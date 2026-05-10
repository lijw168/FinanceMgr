package business

import (
	"encoding/binary"
	"encoding/json"
	"financeMgr/src/analysis-server/sdk/options"
	sdkUtil "financeMgr/src/analysis-server/sdk/util"
	"financeMgr/src/client/util"
	"fmt"
)

type CompanyGateway struct {
}

func (cg *CompanyGateway) ListCompany(param []byte) (resData []byte, errCode int, errMsg string) {
	return listCmdJson(resource_type_company, param, cSdk.ListCompany_json)
}

func (cg *CompanyGateway) GetCompany(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	//因为是直接调用的API接口，所以是通过小端的方式把companyId传过来的。
	id := int(binary.LittleEndian.Uint32(param))
	if id <= 0 {
		errMsg = fmt.Sprintf("the id param is: %d", id)
		logger.Error(errMsg)
		errCode = util.ErrInvalidParam
		return nil, errCode, errMsg
	}
	var opts options.BaseOptions
	opts.ID = id
	view, err := cSdk.GetCompany(&opts)
	if err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrShowFailed
			errMsg = "GetCompany failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		resData, err = json.Marshal(view)
		if err != nil {
			errCode = util.ErrMarshalFailed
			errMsg = fmt.Sprintf("the Marshal failed,err:%v", err.Error())
			logger.LogError(errMsg)
		}
		logger.Debug("GetCompany succeed;views:%v", view)
	}
	return resData, errCode, errMsg
}

func (cg *CompanyGateway) CreateCompany(param []byte) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	if views, err := cSdk.CreateCompany_json(param); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrCreateFailed
			errMsg = "CreateCompany failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("CreateCompany succeed;views:%v", views)
		resData = make([]byte, 4)
		binary.LittleEndian.PutUint32(resData, uint32(views.CompanyID))
	}
	return resData, errCode, errMsg
}

func (cg *CompanyGateway) UpdateCompany(param []byte) (errCode int, errMsg string) {
	errCode = util.ErrNull
	if err := cSdk.UpdateCompany_json(param); err != nil {
		errCode = util.ErrUpdateFailed
		logger.Error("the UpdateCompany failed,err:%v", err.Error())
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrUpdateFailed
			errMsg = "UpdateCompany failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("UpdateCompany succeed")
	}
	return errCode, errMsg
}

func (cg *CompanyGateway) DeleteCompany(param []byte) (errCode int, errMsg string) {
	id := int(binary.LittleEndian.Uint32(param))
	return deleteCmd(resource_type_company, id, cSdk.DeleteCompany)
}

func (cg *CompanyGateway) ListCompanyAccountYearInfo(param []byte) (resData []byte, errCode int, errMsg string) {
	//因为是直接调用的API接口，所以是通过小端的方式把companyId传过来的。
	id := int(binary.LittleEndian.Uint32(param))
	if id <= 0 {
		errMsg = fmt.Sprintf("the id param is: %d", id)
		logger.Error(errMsg)
		errCode = util.ErrInvalidParam
		return nil, errCode, errMsg
	}
	var opts options.BaseOptions
	opts.ID = id
	errCode = util.ErrNull
	if views, err := cSdk.ListCompanyAccountYearInfo(&opts); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrListCompanyAccountYearInfo
			errMsg = "ListCompanyAccountYearInfo failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		resData, err = json.Marshal(views)
		if err != nil {
			errCode = util.ErrMarshalFailed
			errMsg = fmt.Sprintf("the Marshal failed,err:%v", err.Error())
			logger.LogError(errMsg)
		} else {
			logger.Debug("ListCompanyAccountYearInfo succeed;views:%v", views)
		}
	}
	return resData, errCode, errMsg
}
