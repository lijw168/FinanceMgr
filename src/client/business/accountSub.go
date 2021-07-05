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

// func (ag *AccSubGateway) ListAccSub(param []byte) (resData []byte, errCode int) {
// 	var opts options.ListOptions
// 	errCode = util.ErrNull
// 	if err := json.Unmarshal(param, &opts); err != nil {
// 		logger.Error("the Unmarshal failed,err:%v", err.Error())
// 		errCode = util.ErrUnmarshalFailed
// 		return nil, errCode
// 	}
// 	if count, views, err := cSdk.ListAccSub(&opts); err != nil {
// 		logger.Error("the ListAccSub failed,err:%v", err.Error())
// 	} else {
// 		logger.Debug("ListAccSub succeed;count:%d,views:%v", count, views)
// 		desc := &(model.DescData{})
// 		desc.Tc = count
// 		desc.Elements = views
// 		resData, err = json.Marshal(desc)
// 		if err != nil {
// 			errCode = util.ErrMarshalFailed
// 			logger.Error("the Marshal failed,err:%v", err.Error())
// 		}
// 	}
// 	return resData, errCode
// }

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

// func (ag *AccSubGateway) CreateAccSub(param []byte) (resData []byte, errCode int) {
// 	var opts options.CreateSubjectOptions
// 	errCode = util.ErrNull
// 	if err := json.Unmarshal(param, &opts); err != nil {
// 		logger.Error("the Unmarshal failed,err:%v", err.Error())
// 		errCode = util.ErrUnmarshalFailed
// 		return nil, errCode
// 	}
// 	if views, err := cSdk.CreateAccSub(&opts); err != nil {
// 		errCode = util.ErrCreateFailed
// 		logger.Error("the CreateAccSub failed,err:%v", err.Error())
// 	} else {
// 		logger.Debug("CreateAccSub succeed;views:%v", views)
// 		resData = make([]byte, 4)
// 		binary.LittleEndian.PutUint32(resData, uint32(views.SubjectID))
// 		// resData, err = json.Marshal(views)
// 		// if err != nil {
// 		// 	errCode = util.ErrMarshalFailed
// 		// 	logger.Error("the Marshal failed,err:%v", err.Error())
// 		// }
// 	}
// 	return resData, errCode
// }

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

// func (ag *AccSubGateway) UpdateAccSub(param []byte) (errCode int) {
// 	var opts options.ModifySubjectOptions
// 	errCode = util.ErrNull
// 	if err := json.Unmarshal(param, &opts); err != nil {
// 		logger.Error("the Unmarshal failed,err:%v", err.Error())
// 		errCode = util.ErrUnmarshalFailed
// 		return errCode
// 	}
// 	if err := cSdk.UpdateAccSub(&opts); err != nil {
// 		errCode = util.ErrUpdateFailed
// 		logger.Error("the UpdateAccSub failed,err:%v", err.Error())
// 	} else {
// 		logger.Debug("UpdateAccSub succeed")
// 	}
// 	return errCode
// }

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
	return deleteCmd(resource_type_account_sub, id, cSdk.DeleteAccSub)
}
