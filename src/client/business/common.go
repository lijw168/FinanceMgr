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
)

func deleteCmd(rsT, id int, handler func(*options.BaseOptions) error) (errCode int) {
	errCode = util.ErrNull
	if id <= 0 {
		logger.Error("the id param is: %d", id)
		errCode = util.ErrInvalidParam
		return errCode
	}
	rsName := getResourceName(rsT)
	var opts options.BaseOptions
	opts.ID = id
	if err := handler(&opts); err != nil {
		errCode = util.ErrDeleteFailed
		logger.Error("delete a %s failed,err:%v", rsName, err.Error())
	} else {
		logger.Debug("delete a %s succeed.", rsName)
	}
	return errCode
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
	default:
		panic("Unsupport resource type")
	}
}
