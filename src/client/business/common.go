package business

import (
	//"financeMgr/src/analysis-server/model"
	"financeMgr/src/analysis-server/sdk/options"
	sdkUtil "financeMgr/src/analysis-server/sdk/util"
	"financeMgr/src/client/util"
	"fmt"
	//"encoding/json"
)

const (
	resource_type_account_sub = iota
	resource_type_year_balance
	resource_type_company
	resource_type_voucher
	resource_type_voucher_record
	resource_type_voucher_info
	resource_type_voucher_template
	resource_type_login_info
	resource_type_operator
	resource_type_menu_info
)

func deleteCmd(rsT, id int, delHandler func(*options.BaseOptions) error) (errCode int, errMsg string) {
	errCode = util.ErrNull
	if id <= 0 {
		logger.Error("the id param is: %d", id)
		errCode = util.ErrInvalidParam
		return errCode, errMsg
	}
	rsName := getResourceName(rsT)
	var opts options.BaseOptions
	opts.ID = id
	if err := delHandler(&opts); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrDeleteFailed
			errMsg = fmt.Sprintf("delete %s failed,internal error", rsName)
		}
		logger.Error("the delete %s failed,err:%s", rsName, errMsg)
	} else {
		logger.Debug("delete a %s succeed.", rsName)
	}
	return errCode, errMsg
}

type listhandler func([]byte) ([]byte, error)

// params是json格式的参数数据。
func listCmdJson(rsT int, params []byte, handler listhandler) (resData []byte, errCode int, errMsg string) {
	errCode = util.ErrNull
	var err error
	rsName := getResourceName(rsT)
	if resData, err = handler(params); err != nil {
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrListFailed
			errMsg = fmt.Sprintf("List %s failed,internal error", rsName)
		}
		logger.Error("the List %s failed,err:%s", rsName, errMsg)
	} else {
		logger.Debug("List %s succeed;", rsName)
	}
	return resData, errCode, errMsg
}

func getResourceName(rsT int) string {
	switch rsT {
	case resource_type_account_sub:
		return "accSub"
	case resource_type_year_balance:
		return "yearBalance"
	case resource_type_company:
		return "company"
	case resource_type_voucher:
		return "voucher"
	case resource_type_voucher_record:
		return "vouRecord"
	case resource_type_voucher_info:
		return "vouInfo"
	case resource_type_voucher_template:
		return "vouTemplate"
	case resource_type_login_info:
		return "loginInfo"
	case resource_type_operator:
		return "operator"
	case resource_type_menu_info:
		return "menuInfo"
	default:
		panic("Unsupport resource type")
	}
}
