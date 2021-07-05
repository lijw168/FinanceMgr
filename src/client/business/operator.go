package business

import (
	//"analysis-server/model"
	"analysis-server/sdk/options"
	"client/util"
	"encoding/json"
)

type OperatorGateway struct {
}

func (og *OperatorGateway) ListOperatorInfo(param []byte) (resData []byte, errCode int) {
	return listCmdJson(resource_type_operator, param, cSdk.ListOperatorInfo_json)
}

func (og *OperatorGateway) GetOperatorInfo(param []byte) (resData []byte, errCode int) {
	strName := string(param)
	errCode = util.ErrNull
	var opts options.NameOptions
	opts.Name = strName
	if opts.Name == "" {
		logger.Error("the name param is empty")
		errCode = util.ErrInvalidParam
		return nil, errCode
	}
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

// func (og *OperatorGateway) CreateOperator(param []byte) (errCode int) {
// 	var opts options.CreateOptInfoOptions
// 	errCode = util.ErrNull
// 	if err := json.Unmarshal(param, &opts); err != nil {
// 		logger.Error("the Unmarshal failed,err:%v", err.Error())
// 		errCode = util.ErrUnmarshalFailed
// 		return errCode
// 	}
// 	if views, err := cSdk.CreateOperator(&opts); err != nil {
// 		errCode = util.ErrCreateFailed
// 		logger.Error("the CreateOperator failed,err:%v", err.Error())
// 	} else {
// 		logger.Debug("CreateOperator succeed;views:%v", views)
// 		// resData, err = json.Marshal(views)
// 		// if err != nil {
// 		// 	errCode = util.ErrMarshalFailed
// 		// 	logger.Error("the Marshal failed,err:%v", err.Error())
// 		// }
// 	}
// 	return errCode
// }

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

// func (og *OperatorGateway) UpdateOperator(param []byte) (errCode int) {
// 	var opts options.ModifyOptInfoOptions
// 	errCode = util.ErrNull
// 	if err := json.Unmarshal(param, &opts); err != nil {
// 		logger.Error("the Unmarshal failed,err:%v", err.Error())
// 		errCode = util.ErrUnmarshalFailed
// 		return errCode
// 	}
// 	if err := cSdk.UpdateOperator(&opts); err != nil {
// 		errCode = util.ErrUpdateFailed
// 		logger.Error("the UpdateOperator failed,err:%v", err.Error())
// 	} else {
// 		logger.Debug("UpdateOperator succeed")
// 	}
// 	return errCode
// }

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
	strName := string(param)
	if strName == "" {
		logger.Error("the name param is empty")
		errCode = util.ErrInvalidParam
		return errCode
	}
	errCode = util.ErrNull
	var opts options.NameOptions
	opts.Name = strName
	if err := cSdk.DeleteOperator(&opts); err != nil {
		errCode = util.ErrDeleteFailed
		logger.Error("the DeleteOperator failed,err:%v", err.Error())
	} else {
		logger.Debug("DeleteOperator succeed;")
	}
	return errCode
}
