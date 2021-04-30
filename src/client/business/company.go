package business

import (
	"analysis-server/model"
	"analysis-server/sdk/options"
	"client/util"
	"encoding/binary"
	"encoding/json"
)

type CompanyGateway struct {
}

func (cg *CompanyGateway) ListCompany(param []byte) (resData []byte, errCode int) {
	var opts options.ListOptions
	errCode = util.ErrNull
	if err := json.Unmarshal(param, &opts); err != nil {
		logger.Error("the Unmarshal failed,err:%v", err.Error())
		errCode = util.ErrUnmarshalFailed
		return nil, errCode
	}
	if count, views, err := cSdk.ListCompany(&opts); err != nil {
		logger.Error("the ListCompany failed,err:%v", err.Error())
	} else {
		logger.Debug("ListCompany succeed;count:%d,views:%v", count, views)
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
	var opts options.CreateCompanyOptions
	errCode = util.ErrNull
	if err := json.Unmarshal(param, &opts); err != nil {
		logger.Error("the Unmarshal failed,err:%v", err.Error())
		errCode = util.ErrUnmarshalFailed
		return nil, errCode
	}
	if views, err := cSdk.CreateCompany(&opts); err != nil {
		errCode = util.ErrCreateFailed
		logger.Error("the CreateCompany failed,err:%v", err.Error())
	} else {
		logger.Debug("CreateCompany succeed;views:%v", views)
		resData, err = json.Marshal(views)
		if err != nil {
			errCode = util.ErrMarshalFailed
			logger.Error("the Marshal failed,err:%v", err.Error())
		}
	}
	return resData, errCode
}

func (cg *CompanyGateway) UpdateCompany(param []byte) (errCode int) {
	var opts options.ModifyCompanyOptions
	errCode = util.ErrNull
	if err := json.Unmarshal(param, &opts); err != nil {
		logger.Error("the Unmarshal failed,err:%v", err.Error())
		errCode = util.ErrUnmarshalFailed
		return errCode
	}
	if err := cSdk.UpdateCompany(&opts); err != nil {
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
