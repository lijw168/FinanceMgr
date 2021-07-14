package business

import (
	//"analysis-server/model"
	"analysis-server/sdk/options"
	"client/util"
	"encoding/binary"
	"encoding/json"
)

type CompanyGateway struct {
}

func (cg *CompanyGateway) ListCompany(param []byte) (resData []byte, errCode int) {
	return listCmdJson(resource_type_company, param, cSdk.ListCompany_json)
}

func (cg *CompanyGateway) GetCompany(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	id := int(binary.LittleEndian.Uint32(param))
	if id <= 0 {
		logger.Error("the id param is: %d", id)
		errCode = util.ErrInvalidParam
		return nil, errCode
	}
	var opts options.BaseOptions
	opts.ID = id
	view, err := cSdk.GetCompany(&opts)
	if err != nil {
		errCode = util.ErrShowFailed
		logger.Error("the GetCompany failed,err:%v", err.Error())
	} else {
		resData, err = json.Marshal(view)
		if err != nil {
			errCode = util.ErrMarshalFailed
			logger.Error("the Marshal failed,err:%v", err.Error())
		}
		logger.Debug("GetCompany succeed;views:%v", view)
	}
	return resData, errCode
}

func (cg *CompanyGateway) CreateCompany(param []byte) (resData []byte, errCode int) {
	errCode = util.ErrNull
	if views, err := cSdk.CreateCompany_json(param); err != nil {
		errCode = util.ErrCreateFailed
		logger.Error("the CreateCompany failed,err:%v", err.Error())
	} else {
		logger.Debug("CreateCompany succeed;views:%v", views)
		resData = make([]byte, 4)
		binary.LittleEndian.PutUint32(resData, uint32(views.CompanyID))
	}
	return resData, errCode
}

func (cg *CompanyGateway) UpdateCompany(param []byte) (errCode int) {
	errCode = util.ErrNull
	if err := cSdk.UpdateCompany_json(param); err != nil {
		errCode = util.ErrUpdateFailed
		logger.Error("the UpdateCompany failed,err:%v", err.Error())
	} else {
		logger.Debug("UpdateCompany succeed")
	}
	return errCode
}

func (cg *CompanyGateway) DeleteCompany(param []byte) (errCode int) {
	id := int(binary.LittleEndian.Uint32(param))
	return deleteCmd(resource_type_company, id, cSdk.DeleteCompany)
}
