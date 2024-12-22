package sdk

import (
	"financeMgr/src/common/log"
	"net/http"
	"time"

	"financeMgr/src/analysis-server/sdk/mgr"
	"financeMgr/src/analysis-server/sdk/util"
)

type CcSdk struct {
	Domain string
	//Tenant  string
	Verbose bool
	Admin   bool
	//由于需要每一个req，一个traceId，所以删除共用的traceId.
	//TraceId string
	Timeout uint64
	//AccessToken string
	mgr.AccSub
	mgr.Company
	mgr.CompanyGroup
	mgr.Operator
	mgr.Voucher
	mgr.Authen
	mgr.MenuInfo
	mgr.YearBalance
	Logger *log.Logger
}

func (c *CcSdk) Setup() {
	util.Domain = c.Domain
	util.AccessToken = ""
	util.Verbose = c.Verbose
	util.Admin = c.Admin
	util.Client = new(http.Client)
	util.Client.Timeout = time.Duration(c.Timeout) * time.Millisecond
	//util.TraceId = c.TraceId
	util.Logger = c.Logger
}

func (c *CcSdk) SetAccessToken(accessToken string) {
	// util.TokenMutex.Lock()
	// defer util.TokenMutex.Unlock()
	util.AccessToken = accessToken
}
