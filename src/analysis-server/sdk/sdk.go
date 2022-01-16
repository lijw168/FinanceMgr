package sdk

import (
	"net/http"
	"time"

	"analysis-server/sdk/mgr"
	"analysis-server/sdk/util"
	"common/log"
)

type CcSdk struct {
	Domain string
	//Tenant  string
	Verbose bool
	Admin   bool
	TraceId string
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
	util.TraceId = c.TraceId
	util.Logger = c.Logger
}

func (c *CcSdk) SetAccessToken(accessToken string) {
	// util.TokenMutex.Lock()
	// defer util.TokenMutex.Unlock()
	util.AccessToken = accessToken
}
