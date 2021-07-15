package handler

import (
	//"errors"
	//"net/http"
	"analysis-server/model"
)

// func Auth(secretIds []string) bool {
// 	auth := func(secretId string) bool {
// 		return true
// 	}
// 	for _, secretId := range secretIds {
// 		if auth(secretId) {
// 			return true
// 		}
// 	}
// 	return false
// }

// func isAdmin(req *http.Request) bool {
// 	return Auth(req.Header["Secret-Id"])
// }
//在baseParams中，只要满足其中一个就可以。
func isLackBaseParams(baseParams []string, queryParams []*model.FilterItem) bool {
	bRet := true
end:
	for _, bP := range baseParams {
		for _, f := range queryParams {
			if *f.Field == bP {
				bRet = false
				break end
			}
		}
	}
	return bRet
}
