package business

import (
//"analysis-server/model"
// "analysis-server/sdk/options"
// "client/util"
// "encoding/binary"
// "encoding/json"
)

type MenuInfoGateway struct {
}

func (ag *MenuInfoGateway) ListMenuInfo(param []byte) (resData []byte, errCode int) {
	return listCmdJson(resource_type_menu_info, param, cSdk.ListMenuInfo_json)
}
