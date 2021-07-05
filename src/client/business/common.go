package business

import (
	//"analysis-server/model"
	"analysis-server/sdk/options"
	"client/util"
	//"encoding/json"
)

const (
	resource_type_account_sub = iota
	resource_type_company
	resource_type_voucher
	resource_type_voucher_record
	resource_type_voucher_info
	resource_type_login_info
	resource_type_operator
)

func deleteCmd(rsT, id int, delHandler func(*options.BaseOptions) error) (errCode int) {
	errCode = util.ErrNull
	if id <= 0 {
		logger.Error("the id param is: %d", id)
		errCode = util.ErrInvalidParam
		return errCode
	}
	rsName := getResourceName(rsT)
	var opts options.BaseOptions
	opts.ID = id
	if err := delHandler(&opts); err != nil {
		errCode = util.ErrDeleteFailed
		logger.Error("delete a %s failed,err:%v", rsName, err.Error())
	} else {
		logger.Debug("delete a %s succeed.", rsName)
	}
	return errCode
}

type listhandler func([]byte) ([]byte, error)

//params是json格式的参数数据。
func listCmdJson(rsT int, params []byte, handler listhandler) (resData []byte, errCode int) {
	errCode = util.ErrNull
	var err error
	rsName := getResourceName(rsT)
	if resData, err = handler(params); err != nil {
		errCode = util.ErrListFailed
		logger.Error("the List%s failed,err:%v", rsName, err.Error())
	} else {
		logger.Debug("List%s succeed;", rsName)
	}
	return resData, errCode
}

func getResourceName(rsT int) string {
	switch rsT {
	case resource_type_account_sub:
		return "accSub"
	case resource_type_company:
		return "company"
	case resource_type_voucher:
		return "voucher"
	case resource_type_voucher_record:
		return "vouRecord"
	case resource_type_voucher_info:
		return "vouInfo"
	case resource_type_login_info:
		return "loginInfo"
	case resource_type_operator:
		return "operator"
	default:
		panic("Unsupport resource type")
	}
}
